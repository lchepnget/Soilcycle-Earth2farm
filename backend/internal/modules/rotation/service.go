package rotation

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strings"

	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

type cropCandidate struct {
	ID               uuid.UUID
	Name             string
	Family           string
	Category         string
	NitrogenFixing   bool
	WaterRequirement string
	NutrientDemand   string
	GrowthDuration   *int
}

type cropHistory struct {
	CropID uuid.UUID
	Family string
}

// GenerateRotation creates and saves alternative crop rotation plans.
func (s *Service) GenerateRotation(
	ctx context.Context,
	req CreateRotationRequest,
) (*RotationPlan, []RotationResult, error) {

	if err := req.Validate(); err != nil {
		return nil, nil, err
	}

	farmID, err := uuid.Parse(req.FarmID)
	if err != nil {
		return nil, nil, errors.New("invalid farm ID")
	}

	exists, err := s.repo.FarmExists(ctx, farmID)
	if err != nil {
		return nil, nil, err
	}

	if !exists {
		return nil, nil, errors.New("farm not found")
	}

	crops, err := s.getAvailableCrops(ctx)
	if err != nil {
		return nil, nil, err
	}

	if len(crops) == 0 {
		return nil, nil, errors.New("no crops available for rotation planning")
	}

	history, err := s.getCropHistory(ctx, farmID)
	if err != nil {
		return nil, nil, err
	}

	prioritiesJSON, err := json.Marshal(req.FarmerPriorities)
	if err != nil {
		return nil, nil, fmt.Errorf("encode farmer priorities: %w", err)
	}

	snapshot := map[string]interface{}{
		"crop_history": history,
		"engine":       "rule_based_v1",
		"data_limitations": []string{
			"Recommendations depend on available crop records and characteristics.",
			"Climate and soil suitability are not yet incorporated into scoring.",
			"Recommendations require local agronomic validation.",
		},
	}

	snapshotJSON, err := json.Marshal(snapshot)
	if err != nil {
		return nil, nil, fmt.Errorf("encode input snapshot: %w", err)
	}

	plan := &RotationPlan{
		FarmID:           farmID,
		FarmerPriorities: prioritiesJSON,
		InputSnapshot:    snapshotJSON,
		NumberOfSeasons:  req.NumberOfSeasons,
	}

	if err := s.repo.CreatePlan(ctx, plan); err != nil {
		return nil, nil, err
	}

	results := s.buildRotationOptions(
		crops,
		history,
		req.NumberOfSeasons,
		req.FarmerPriorities,
		plan.ID,
	)

	if len(results) == 0 {
		return plan, nil, errors.New("could not generate rotation options")
	}

	if err := s.repo.SaveResults(ctx, results); err != nil {
		return plan, nil, err
	}

	savedResults, err := s.repo.GetResultsByPlan(ctx, plan.ID)
	if err != nil {
		return plan, nil, err
	}

	return plan, savedResults, nil
}

// GetPlan retrieves a saved rotation plan.
func (s *Service) GetPlan(
	ctx context.Context,
	planID string,
) (*RotationPlan, error) {

	id, err := uuid.Parse(planID)
	if err != nil {
		return nil, errors.New("invalid rotation plan ID")
	}

	plan, err := s.repo.GetPlanByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if plan == nil {
		return nil, errors.New("rotation plan not found")
	}

	return plan, nil
}

// GetPlanResults retrieves the generated alternatives for a plan.
func (s *Service) GetPlanResults(
	ctx context.Context,
	planID string,
) ([]RotationResult, error) {

	id, err := uuid.Parse(planID)
	if err != nil {
		return nil, errors.New("invalid rotation plan ID")
	}

	plan, err := s.repo.GetPlanByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if plan == nil {
		return nil, errors.New("rotation plan not found")
	}

	return s.repo.GetResultsByPlan(ctx, id)
}

// ListFarmPlans retrieves a farm's saved rotation plans.
func (s *Service) ListFarmPlans(
	ctx context.Context,
	farmID string,
	limit int,
	offset int,
) ([]RotationPlan, error) {

	id, err := uuid.Parse(farmID)
	if err != nil {
		return nil, errors.New("invalid farm ID")
	}

	exists, err := s.repo.FarmExists(ctx, id)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, errors.New("farm not found")
	}

	return s.repo.ListPlansByFarm(ctx, id, limit, offset)
}

// getAvailableCrops retrieves crops and their characteristics.
func (s *Service) getAvailableCrops(
	ctx context.Context,
) ([]cropCandidate, error) {

	rows, err := s.repo.db.Query(
		ctx,
		`SELECT
			c.id,
			c.name,
			COALESCE(c.crop_family, ''),
			COALESCE(c.category, ''),
			COALESCE(cc.nitrogen_fixing, FALSE),
			COALESCE(cc.water_requirement, ''),
			COALESCE(cc.nutrient_demand, ''),
			cc.growth_duration_days
		FROM crops c
		LEFT JOIN crop_characteristics cc ON cc.crop_id = c.id
		ORDER BY c.name`,
	)

	if err != nil {
		return nil, fmt.Errorf("retrieve available crops: %w", err)
	}

	defer rows.Close()

	crops := make([]cropCandidate, 0)

	for rows.Next() {
		var crop cropCandidate

		err := rows.Scan(
			&crop.ID,
			&crop.Name,
			&crop.Family,
			&crop.Category,
			&crop.NitrogenFixing,
			&crop.WaterRequirement,
			&crop.NutrientDemand,
			&crop.GrowthDuration,
		)

		if err != nil {
			return nil, fmt.Errorf("scan crop: %w", err)
		}

		crops = append(crops, crop)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate crops: %w", err)
	}

	return crops, nil
}

// getCropHistory retrieves crops previously grown on a farm.
func (s *Service) getCropHistory(
	ctx context.Context,
	farmID uuid.UUID,
) ([]cropHistory, error) {

	rows, err := s.repo.db.Query(
		ctx,
		`SELECT
			fc.crop_id,
			COALESCE(c.crop_family, '')
		FROM farm_crops fc
		JOIN crops c ON c.id = fc.crop_id
		WHERE fc.farm_id = $1
		ORDER BY fc.planting_date DESC NULLS LAST, fc.created_at DESC
		LIMIT 10`,
		farmID,
	)

	if err != nil {
		return nil, fmt.Errorf("retrieve farm crop history: %w", err)
	}

	defer rows.Close()

	history := make([]cropHistory, 0)

	for rows.Next() {
		var item cropHistory

		if err := rows.Scan(&item.CropID, &item.Family); err != nil {
			return nil, fmt.Errorf("scan crop history: %w", err)
		}

		history = append(history, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate crop history: %w", err)
	}

	return history, nil
}

// buildRotationOptions generates three alternative rotation strategies.
func (s *Service) buildRotationOptions(
	crops []cropCandidate,
	history []cropHistory,
	seasons int,
	priorities []string,
	planID uuid.UUID,
) []RotationResult {

	strategies := []string{
		"Balanced Rotation",
		"Soil Restoration Focus",
		"Water Conscious Rotation",
	}

	results := make([]RotationResult, 0, len(strategies))

	for strategyIndex, strategy := range strategies {

		sequence := make([]CropRecommendation, 0, seasons)
		used := make(map[uuid.UUID]bool)

		lastCropID := uuid.Nil
		lastFamily := ""

		for season := 1; season <= seasons; season++ {

			crop := chooseCrop(
				crops,
				used,
				lastCropID,
				lastFamily,
				history,
				priorities,
				strategyIndex,
			)

			if crop == nil {
				break
			}

			reason, benefits := explainCrop(
				crop,
				lastFamily,
				priorities,
				strategyIndex,
			)

			sequence = append(sequence, CropRecommendation{
				CropID:   crop.ID,
				CropName: crop.Name,
				Season:   season,
				Reason:   reason,
				Benefits: benefits,
			})

			used[crop.ID] = true
			lastCropID = crop.ID
			lastFamily = crop.Family
		}

		if len(sequence) == 0 {
			continue
		}

		score := calculateRotationScore(
			sequence,
			crops,
			history,
		)

		cropJSON, _ := json.Marshal(sequence)

		scoresJSON, _ := json.Marshal(score)

		explanationJSON, _ := json.Marshal([]string{
			"Generated using transparent rule-based crop selection.",
			"Crop diversity and nitrogen-fixing characteristics contribute to the score.",
			"Climate suitability and measured soil conditions are not yet included.",
			"Review recommendations with a local agricultural extension officer.",
		})

		results = append(results, RotationResult{
			PlanID:       planID,
			PlanName:     strategy,
			CropSequence: cropJSON,
			Scores:       scoresJSON,
			Explanation:  explanationJSON,
			OverallScore: &score.OverallScore,
		})
	}

	return results
}

func chooseCrop(
	crops []cropCandidate,
	used map[uuid.UUID]bool,
	lastCropID uuid.UUID,
	lastFamily string,
	history []cropHistory,
	priorities []string,
	strategy int,
) *cropCandidate {

	var best *cropCandidate
	bestScore := math.Inf(-1)

	for i := range crops {
		crop := &crops[i]

		if used[crop.ID] {
			continue
		}

		score := 0.0

		if crop.NitrogenFixing {
			score += 3
		}

		if crop.ID != lastCropID {
			score += 2
		}

		if crop.Family != "" && crop.Family != lastFamily {
			score += 2
		}

		for _, previous := range history {
			if previous.CropID == crop.ID {
				score -= 1
			}

			if crop.Family != "" && previous.Family == crop.Family {
				score -= 0.5
			}
		}

		if matchesPriority(crop, priorities) {
			score += 3
		}

		switch strategy {
		case 1:
			if crop.NitrogenFixing {
				score += 4
			}
		case 2:
			if strings.EqualFold(crop.WaterRequirement, "low") {
				score += 4
			}
		}

		if score > bestScore {
			bestScore = score
			best = crop
		}
	}

	return best
}

func matchesPriority(
	crop *cropCandidate,
	priorities []string,
) bool {

	for _, priority := range priorities {
		priority = strings.ToLower(strings.TrimSpace(priority))

		if priority == "" {
			continue
		}

		if strings.Contains(strings.ToLower(crop.Name), priority) ||
			strings.Contains(strings.ToLower(crop.Category), priority) ||
			strings.Contains(strings.ToLower(crop.Family), priority) {
			return true
		}

		if priority == "soil health" && crop.NitrogenFixing {
			return true
		}

		if priority == "water conservation" &&
			strings.EqualFold(crop.WaterRequirement, "low") {
			return true
		}
	}

	return false
}

func explainCrop(
	crop *cropCandidate,
	lastFamily string,
	priorities []string,
	strategy int,
) (string, []string) {

	reasons := make([]string, 0)
	benefits := make([]string, 0)

	if crop.NitrogenFixing {
		reasons = append(reasons, "This crop is recorded as nitrogen-fixing.")
		benefits = append(benefits, "Potential contribution to biological nitrogen fixation.")
	}

	if crop.Family != "" && crop.Family != lastFamily {
		reasons = append(reasons, "It introduces a different crop family from the previous season.")
		benefits = append(benefits, "Supports crop-family diversification.")
	}

	if matchesPriority(crop, priorities) {
		reasons = append(reasons, "It matches one or more farmer-stated priorities.")
		benefits = append(benefits, "Aligned with farmer preferences.")
	}

	if strategy == 2 &&
		strings.EqualFold(crop.WaterRequirement, "low") {
		reasons = append(reasons, "It is recorded as having low water requirements.")
		benefits = append(benefits, "May be a useful option where water is limited.")
	}

	if len(reasons) == 0 {
		reasons = append(reasons, "Selected based on crop diversity and available crop records.")
	}

	return strings.Join(reasons, " "), benefits
}

func calculateRotationScore(
	sequence []CropRecommendation,
	crops []cropCandidate,
	history []cropHistory,
) RotationScore {

	if len(sequence) == 0 {
		return RotationScore{}
	}

	cropMap := make(map[uuid.UUID]cropCandidate)

	for _, crop := range crops {
		cropMap[crop.ID] = crop
	}

	families := make(map[string]bool)

	nitrogenFixingCount := 0

	for _, recommendation := range sequence {
		crop, ok := cropMap[recommendation.CropID]
		if !ok {
			continue
		}

		if crop.Family != "" {
			families[crop.Family] = true
		}

		if crop.NitrogenFixing {
			nitrogenFixingCount++
		}
	}

	diversityScore := math.Min(
		100,
		float64(len(families))*25,
	)

	soilHealthScore := 40.0 +
		float64(nitrogenFixingCount)*20

	soilHealthScore = math.Min(100, soilHealthScore)

	climateScore := 50.0

	// These are provisional heuristic scores, not measured outcomes.
	// Climate data will be incorporated in a later engine version.

	overall :=
		soilHealthScore*0.4 +
			climateScore*0.2 +
			diversityScore*0.4

	return RotationScore{
		SoilHealthScore: soilHealthScore,
		ClimateScore:    climateScore,
		DiversityScore:  diversityScore,
		OverallScore:    math.Round(overall*100) / 100,
	}
}

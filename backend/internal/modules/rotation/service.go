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
	repo rotationRepository
}

type rotationRepository interface {
	FarmExists(context.Context, uuid.UUID) (bool, error)
	ListAvailableCrops(context.Context) ([]cropCandidate, error)
	ListCropHistory(context.Context, uuid.UUID) ([]cropHistory, error)
	CreatePlan(context.Context, *RotationPlan) error
	SaveResults(context.Context, []RotationResult) error
	GetResultsByPlan(context.Context, uuid.UUID) ([]RotationResult, error)
	GetPlanByID(context.Context, uuid.UUID) (*RotationPlan, error)
	ListPlansByFarm(context.Context, uuid.UUID, int, int) ([]RotationPlan, error)
}

var errInsufficientCrops = errors.New("insufficient crop catalogue")

func NewService(repo rotationRepository) *Service {
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

	crops, err := s.repo.ListAvailableCrops(ctx)
	if err != nil {
		return nil, nil, err
	}

	if len(crops) == 0 {
		return nil, nil, errors.New("no crops available for rotation planning")
	}

	if len(crops) < req.NumberOfSeasons {
		return nil, nil, fmt.Errorf(
			"%w: requested %d seasons but only %d distinct crops are available",
			errInsufficientCrops,
			req.NumberOfSeasons,
			len(crops),
		)
	}

	history, err := s.repo.ListCropHistory(ctx, farmID)
	if err != nil {
		return nil, nil, err
	}

	results, err := buildRotationOptions(
		crops,
		history,
		req.NumberOfSeasons,
		req.FarmerPriorities,
	)
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
			"Climate observations and measured soil data are not used in selection or scoring.",
			"The climate score is a fixed provisional placeholder, not a climate-suitability measurement.",
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

	for i := range results {
		results[i].PlanID = plan.ID
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

type rotationStrategy int

const (
	balancedRotation rotationStrategy = iota
	soilRestorationRotation
	waterConsciousRotation
)

type rotationStrategyOption struct {
	name     string
	strategy rotationStrategy
}

// buildRotationOptions generates complete alternatives using the available catalogue.
func buildRotationOptions(
	crops []cropCandidate,
	history []cropHistory,
	seasons int,
	priorities []string,
) ([]RotationResult, error) {
	if seasons < 1 {
		return nil, errors.New("number of seasons must be greater than zero")
	}

	if len(crops) < seasons {
		return nil, fmt.Errorf(
			"%w: requested %d seasons but only %d distinct crops are available",
			errInsufficientCrops,
			seasons,
			len(crops),
		)
	}

	strategies := []rotationStrategyOption{
		{name: "Balanced Rotation", strategy: balancedRotation},
		{name: "Soil Restoration Focus", strategy: soilRestorationRotation},
		{name: "Water Conscious Rotation", strategy: waterConsciousRotation},
	}

	results := make([]RotationResult, 0, len(strategies))

	for _, strategy := range strategies {
		sequence := make([]CropRecommendation, 0, seasons)
		used := make(map[uuid.UUID]bool)
		usedFamilies := make(map[string]int)

		for season := 1; season <= seasons; season++ {
			crop := chooseCrop(
				crops,
				used,
				usedFamilies,
				history,
				priorities,
				strategy.strategy,
			)
			if crop == nil {
				return nil, fmt.Errorf(
					"%w: could not complete the requested %d-season sequence",
					errInsufficientCrops,
					seasons,
				)
			}

			reason, benefits := explainCrop(
				crop,
				usedFamilies,
				history,
				priorities,
				strategy.strategy,
			)

			sequence = append(sequence, CropRecommendation{
				CropID:   crop.ID,
				CropName: crop.Name,
				Season:   season,
				Reason:   reason,
				Benefits: benefits,
			})

			used[crop.ID] = true
			family := normalizeLabel(crop.Family)
			if family != "" {
				usedFamilies[family]++
			}
		}

		score, err := calculateRotationScore(sequence, crops)
		if err != nil {
			return nil, fmt.Errorf("calculate %s score: %w", strategy.name, err)
		}

		cropJSON, err := json.Marshal(sequence)
		if err != nil {
			return nil, fmt.Errorf("encode %s crop sequence: %w", strategy.name, err)
		}

		scoresJSON, err := json.Marshal(score)
		if err != nil {
			return nil, fmt.Errorf("encode %s scores: %w", strategy.name, err)
		}

		explanationJSON, err := json.Marshal([]string{
			"Generated using deterministic, rule-based crop selection and catalogue records.",
			"All scores are provisional heuristics, not measurements or predictions.",
			"Diversity score is the share of the sequence represented by distinct recorded crop families.",
			"Soil health proxy starts at 40 and adds up to 35 points for recorded nitrogen-fixing share and 25 for low recorded nutrient-demand share.",
			"Climate score is fixed at 50 because climate observations are not used; it is not a climate-suitability score.",
			"Soil health score is a catalogue-based proxy, not a measured soil condition or guaranteed restoration outcome.",
			"Overall score weights the soil proxy at 40%, the climate placeholder at 20%, and diversity at 40%.",
			"Unknown crop characteristics are not treated as low water or nutrient demand.",
			"Review recommendations with a local agricultural extension officer.",
		})
		if err != nil {
			return nil, fmt.Errorf("encode %s explanation: %w", strategy.name, err)
		}

		results = append(results, RotationResult{
			PlanName:     strategy.name,
			CropSequence: cropJSON,
			Scores:       scoresJSON,
			Explanation:  explanationJSON,
			OverallScore: &score.OverallScore,
		})
	}

	return results, nil
}

func chooseCrop(
	crops []cropCandidate,
	used map[uuid.UUID]bool,
	usedFamilies map[string]int,
	history []cropHistory,
	priorities []string,
	strategy rotationStrategy,
) *cropCandidate {
	var best *cropCandidate
	bestScore := math.Inf(-1)

	for i := range crops {
		crop := &crops[i]

		if used[crop.ID] {
			continue
		}

		score := 0.0

		family := normalizeLabel(crop.Family)
		if family != "" {
			if count := usedFamilies[family]; count == 0 {
				diversityWeight := 4.0
				if hasPriority(priorities, "crop diversity") {
					diversityWeight += 3
				}
				score += diversityWeight
			} else {
				score -= float64(count) * 2
			}
		}

		for i, previous := range history {
			if previous.CropID == crop.ID {
				score -= 4 / float64(i+1)
			}

			if family != "" && family == normalizeLabel(previous.Family) {
				score -= 1 / float64(i+1)
			}
		}

		if matchesPriority(crop, priorities) {
			score += 3
		}

		switch strategy {
		case soilRestorationRotation:
			if crop.NitrogenFixing {
				score += 6
			}

			switch strings.ToLower(strings.TrimSpace(crop.NutrientDemand)) {
			case "low":
				score += 2
			case "medium":
				score++
			case "high":
				score--
			}
		case waterConsciousRotation:
			switch strings.ToLower(strings.TrimSpace(crop.WaterRequirement)) {
			case "low":
				score += 6
			case "high":
				score -= 2
			}
		}

		if score > bestScore || (score == bestScore && candidateLess(crop, best)) {
			bestScore = score
			best = crop
		}
	}

	return best
}

func candidateLess(left, right *cropCandidate) bool {
	if right == nil {
		return true
	}

	leftName := normalizeLabel(left.Name)
	rightName := normalizeLabel(right.Name)
	if leftName != rightName {
		return leftName < rightName
	}

	return left.ID.String() < right.ID.String()
}

func matchesPriority(
	crop *cropCandidate,
	priorities []string,
) bool {
	for _, priority := range priorities {
		priority = normalizeLabel(priority)

		if priority == "" {
			continue
		}

		switch priority {
		case "vegetables":
			priority = "vegetable"
		case "cereals":
			priority = "cereal"
		case "legumes":
			priority = "legume"
		case "common bean", "common beans":
			priority = "beans"
		}

		if strings.Contains(normalizeLabel(crop.Name), priority) ||
			strings.Contains(normalizeLabel(crop.Category), priority) ||
			strings.Contains(normalizeLabel(crop.Family), priority) {
			return true
		}

		if (priority == "soil health" || priority == "soil restoration") &&
			crop.NitrogenFixing {
			return true
		}

		if (priority == "water conservation" ||
			priority == "water efficiency" ||
			priority == "water saving") &&
			strings.EqualFold(strings.TrimSpace(crop.WaterRequirement), "low") {
			return true
		}
	}

	return false
}

func hasPriority(priorities []string, expected string) bool {
	for _, priority := range priorities {
		if normalizeLabel(priority) == expected {
			return true
		}
	}

	return false
}

func explainCrop(
	crop *cropCandidate,
	usedFamilies map[string]int,
	history []cropHistory,
	priorities []string,
	strategy rotationStrategy,
) (string, []string) {
	reasons := make([]string, 0)
	benefits := make([]string, 0)

	if crop.NitrogenFixing {
		reasons = append(reasons, "The catalogue records nitrogen-fixing potential; actual fixation depends on field conditions.")
		if strategy == soilRestorationRotation {
			benefits = append(benefits, "Nitrogen-fixing potential is weighted by this strategy, not a claim of measured soil improvement.")
		}
	}

	family := normalizeLabel(crop.Family)
	if family != "" && usedFamilies[family] == 0 {
		reasons = append(reasons, "It adds a crop family not yet used in this sequence.")
		benefits = append(benefits, "Supports crop-family diversification.")
		if hasPriority(priorities, "crop diversity") {
			reasons = append(reasons, "This also advances the farmer-stated crop-diversity priority.")
		}
	}

	if matchesPriority(crop, priorities) {
		reasons = append(reasons, "It matches one or more farmer-stated priorities.")
		benefits = append(benefits, "Aligned with farmer preferences.")
	}

	if strategy == soilRestorationRotation {
		switch strings.ToLower(strings.TrimSpace(crop.NutrientDemand)) {
		case "low":
			reasons = append(reasons, "Its recorded nutrient demand is low, which this strategy modestly favors.")
		case "medium":
			reasons = append(reasons, "Its recorded nutrient demand is medium.")
		case "high":
			reasons = append(reasons, "Its recorded nutrient demand is high, so this strategy does not favor it on that attribute.")
		}
	}

	if strategy == waterConsciousRotation {
		switch strings.ToLower(strings.TrimSpace(crop.WaterRequirement)) {
		case "low":
			reasons = append(reasons, "It is recorded as having low water requirements; this does not mean it is drought-proof or suitable for every dry region.")
			benefits = append(benefits, "Low recorded water requirement is favored, not a guarantee of water suitability.")
		case "medium":
			reasons = append(reasons, "Its recorded water requirement is medium.")
		case "high":
			reasons = append(reasons, "Its recorded water requirement is high, so this strategy does not favor it on that attribute.")
		}
	}

	if !appearsInHistory(crop.ID, history) {
		reasons = append(reasons, "This crop is not present in the recent farm history considered.")
	}

	if crop.GrowthDuration != nil {
		reasons = append(reasons, fmt.Sprintf(
			"Catalogue duration is approximately %d days to maturity or first harvest; actual timing varies by variety and site.",
			*crop.GrowthDuration,
		))
	}

	if len(reasons) == 0 {
		reasons = append(reasons, "Selected using family diversity, recent crop history, farmer priorities, and available catalogue data.")
	}

	if crop.Family == "" || crop.Category == "" || crop.WaterRequirement == "" ||
		crop.NutrientDemand == "" || crop.GrowthDuration == nil {
		reasons = append(reasons, "Some crop characteristics are not recorded; unknown values were left neutral.")
	}

	return strings.Join(reasons, " "), benefits
}

func appearsInHistory(cropID uuid.UUID, history []cropHistory) bool {
	for _, previous := range history {
		if previous.CropID == cropID {
			return true
		}
	}

	return false
}

func normalizeLabel(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	value = strings.NewReplacer("_", " ", "-", " ").Replace(value)

	return strings.Join(strings.Fields(value), " ")
}

func calculateRotationScore(
	sequence []CropRecommendation,
	crops []cropCandidate,
) (RotationScore, error) {
	if len(sequence) == 0 {
		return RotationScore{}, errors.New("rotation sequence cannot be empty")
	}

	cropMap := make(map[uuid.UUID]cropCandidate)

	for _, crop := range crops {
		cropMap[crop.ID] = crop
	}

	families := make(map[string]bool)
	nitrogenFixingCount := 0
	lowNutrientDemandCount := 0
	seenCrops := make(map[uuid.UUID]bool)

	for i, recommendation := range sequence {
		crop, ok := cropMap[recommendation.CropID]
		if recommendation.CropID == uuid.Nil || !ok {
			return RotationScore{}, fmt.Errorf("rotation sequence contains unknown crop ID %q", recommendation.CropID)
		}

		if seenCrops[recommendation.CropID] {
			return RotationScore{}, fmt.Errorf("rotation sequence repeats crop %q", crop.Name)
		}

		if recommendation.Season != i+1 {
			return RotationScore{}, fmt.Errorf("rotation sequence has invalid season number %d at position %d", recommendation.Season, i+1)
		}
		seenCrops[recommendation.CropID] = true

		if family := normalizeLabel(crop.Family); family != "" {
			families[family] = true
		}

		if crop.NitrogenFixing {
			nitrogenFixingCount++
		}

		if strings.EqualFold(strings.TrimSpace(crop.NutrientDemand), "low") {
			lowNutrientDemandCount++
		}
	}

	diversityScore := float64(len(families)) / float64(len(sequence)) * 100

	soilHealthScore := 40.0 +
		float64(nitrogenFixingCount)/float64(len(sequence))*35 +
		float64(lowNutrientDemandCount)/float64(len(sequence))*25

	climateScore := 50.0

	overall :=
		soilHealthScore*0.4 +
			climateScore*0.2 +
			diversityScore*0.4

	return RotationScore{
		SoilHealthScore: soilHealthScore,
		ClimateScore:    climateScore,
		DiversityScore:  diversityScore,
		OverallScore:    math.Round(overall*100) / 100,
	}, nil
}

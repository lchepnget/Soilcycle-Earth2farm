package rotation

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func rotationTestCrops() []cropCandidate {
	return []cropCandidate{
		{
			ID:               uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			Name:             "A maize",
			Family:           "Poaceae",
			Category:         "cereal",
			WaterRequirement: "medium",
			NutrientDemand:   "high",
		},
		{
			ID:               uuid.MustParse("00000000-0000-0000-0000-000000000002"),
			Name:             "B beans",
			Family:           "Fabaceae",
			Category:         "legume",
			NitrogenFixing:   true,
			WaterRequirement: "medium",
			NutrientDemand:   "medium",
		},
		{
			ID:               uuid.MustParse("00000000-0000-0000-0000-000000000003"),
			Name:             "C sorghum",
			Family:           "Poaceae",
			Category:         "cereal",
			WaterRequirement: "low",
			NutrientDemand:   "low",
		},
		{
			ID:               uuid.MustParse("00000000-0000-0000-0000-000000000004"),
			Name:             "D tomatoes",
			Family:           "Solanaceae",
			Category:         "vegetable",
			WaterRequirement: "high",
			NutrientDemand:   "high",
		},
	}
}

func TestChooseCropDiversifiesFamilies(t *testing.T) {
	crops := rotationTestCrops()
	usedFamilies := map[string]int{"poaceae": 1}
	used := map[uuid.UUID]bool{crops[0].ID: true, crops[2].ID: true}

	got := chooseCrop(crops, used, usedFamilies, nil, nil, balancedRotation)
	if got == nil || got.Family == "Poaceae" {
		t.Fatalf("chooseCrop() = %#v, want an unused crop from a different family", got)
	}
}

func TestChooseCropPrioritizesNitrogenFixingForSoilStrategy(t *testing.T) {
	crops := rotationTestCrops()

	balanced := chooseCrop(crops, map[uuid.UUID]bool{}, map[string]int{}, nil, nil, balancedRotation)
	soil := chooseCrop(crops, map[uuid.UUID]bool{}, map[string]int{}, nil, nil, soilRestorationRotation)

	if balanced == nil || balanced.Name != "A maize" {
		t.Fatalf("balanced choice = %#v, want A maize as the deterministic tie-break", balanced)
	}
	if soil == nil || !soil.NitrogenFixing {
		t.Fatalf("soil strategy choice = %#v, want a nitrogen-fixing crop", soil)
	}
}

func TestChooseCropPrioritizesLowWaterForWaterStrategy(t *testing.T) {
	crops := rotationTestCrops()

	balanced := chooseCrop(crops, map[uuid.UUID]bool{}, map[string]int{}, nil, nil, balancedRotation)
	water := chooseCrop(crops, map[uuid.UUID]bool{}, map[string]int{}, nil, nil, waterConsciousRotation)

	if balanced == nil || balanced.Name != "A maize" {
		t.Fatalf("balanced choice = %#v, want A maize as the deterministic tie-break", balanced)
	}
	if water == nil || !strings.EqualFold(water.WaterRequirement, "low") {
		t.Fatalf("water strategy choice = %#v, want a crop recorded as low-water", water)
	}
}

func TestChooseCropHonorsFarmerWaterPriority(t *testing.T) {
	crops := rotationTestCrops()

	got := chooseCrop(
		crops,
		map[uuid.UUID]bool{},
		map[string]int{},
		nil,
		[]string{"water_efficiency"},
		balancedRotation,
	)
	if got == nil || !strings.EqualFold(got.WaterRequirement, "low") {
		t.Fatalf("chooseCrop() = %#v, want the crop matching water_efficiency", got)
	}
}

func TestChooseCropUsesNutrientDemandForSoilStrategy(t *testing.T) {
	crops := []cropCandidate{
		{
			ID:             uuid.MustParse("00000000-0000-0000-0000-000000000021"),
			Name:           "A legume",
			Family:         "Fabaceae",
			NitrogenFixing: true,
			NutrientDemand: "high",
		},
		{
			ID:             uuid.MustParse("00000000-0000-0000-0000-000000000022"),
			Name:           "B legume",
			Family:         "Fabaceae",
			NitrogenFixing: true,
			NutrientDemand: "low",
		},
	}

	got := chooseCrop(crops, map[uuid.UUID]bool{}, map[string]int{}, nil, nil, soilRestorationRotation)
	if got == nil || got.NutrientDemand != "low" {
		t.Fatalf("chooseCrop() = %#v, want the lower recorded nutrient-demand crop", got)
	}
}

func TestChooseCropIsDeterministicRegardlessOfInputOrder(t *testing.T) {
	crops := rotationTestCrops()
	first := chooseCrop(crops, map[uuid.UUID]bool{}, map[string]int{}, nil, nil, balancedRotation)
	reversed := []cropCandidate{crops[3], crops[2], crops[1], crops[0]}
	second := chooseCrop(reversed, map[uuid.UUID]bool{}, map[string]int{}, nil, nil, balancedRotation)

	if first == nil || second == nil || first.ID != second.ID {
		t.Fatalf("choices differ by input order: first=%#v, reversed=%#v", first, second)
	}
}

func TestMatchesPrioritySupportsFarmerTerms(t *testing.T) {
	crops := rotationTestCrops()
	tests := []struct {
		name     string
		crop     cropCandidate
		priority string
	}{
		{name: "underscore water alias", crop: crops[2], priority: "water_efficiency"},
		{name: "soil health", crop: crops[1], priority: "soil_health"},
		{name: "crop name", crop: crops[0], priority: "maize"},
		{name: "category", crop: crops[1], priority: "legumes"},
		{name: "family", crop: crops[1], priority: "fabaceae"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if !matchesPriority(&test.crop, []string{test.priority}) {
				t.Fatalf("matchesPriority(%q) = false", test.priority)
			}
		})
	}

	if matchesPriority(&crops[0], []string{"water_efficiency"}) {
		t.Fatal("matchesPriority() treated a medium-water crop as water efficient")
	}
	uncharacterized := cropCandidate{Name: "Uncharacterized crop"}
	if matchesPriority(&uncharacterized, []string{"water_efficiency"}) {
		t.Fatal("matchesPriority() treated an unknown water requirement as low")
	}
}

func TestChooseCropConsidersRecentCropHistory(t *testing.T) {
	crops := rotationTestCrops()
	history := []cropHistory{{CropID: crops[0].ID, Family: crops[0].Family}}

	got := chooseCrop(crops[:2], map[uuid.UUID]bool{}, map[string]int{}, history, nil, balancedRotation)
	if got == nil || got.ID == crops[0].ID {
		t.Fatalf("chooseCrop() = %#v, want the crop not recently recorded", got)
	}
}

func TestBuildRotationOptionsReturnsThreeCompleteStrategies(t *testing.T) {
	results, err := buildRotationOptions(rotationTestCrops(), nil, 3, nil)
	if err != nil {
		t.Fatalf("buildRotationOptions() error = %v", err)
	}

	wantNames := []string{
		"Balanced Rotation",
		"Soil Restoration Focus",
		"Water Conscious Rotation",
	}
	if len(results) != len(wantNames) {
		t.Fatalf("got %d strategies, want %d", len(results), len(wantNames))
	}

	firstCrops := make(map[uuid.UUID]bool)

	for i, result := range results {
		if result.PlanName != wantNames[i] {
			t.Errorf("strategy %d = %q, want %q", i, result.PlanName, wantNames[i])
		}

		var sequence []CropRecommendation
		if err := json.Unmarshal(result.CropSequence, &sequence); err != nil {
			t.Fatalf("decode %s sequence: %v", result.PlanName, err)
		}
		if len(sequence) != 3 {
			t.Errorf("%s sequence has %d seasons, want 3", result.PlanName, len(sequence))
			continue
		}
		if firstCrops[sequence[0].CropID] {
			t.Errorf("%s used a first crop already chosen by another strategy", result.PlanName)
		}
		firstCrops[sequence[0].CropID] = true
		for season, recommendation := range sequence {
			if recommendation.Season != season+1 {
				t.Errorf("%s has season %d at position %d", result.PlanName, recommendation.Season, season+1)
			}
		}

		var explanation []string
		if err := json.Unmarshal(result.Explanation, &explanation); err != nil {
			t.Fatalf("decode %s explanation: %v", result.PlanName, err)
		}
		joinedExplanation := strings.Join(explanation, " ")
		if !strings.Contains(joinedExplanation, "fixed at 50") ||
			!strings.Contains(joinedExplanation, "provisional heuristics") {
			t.Errorf("%s explanation does not identify provisional scores: %v", result.PlanName, explanation)
		}
	}

	if len(firstCrops) != len(wantNames) {
		t.Errorf("strategies did not produce different first choices: %v", firstCrops)
	}
}

func TestGeneratedExplanationTextPreservesWordSpacing(t *testing.T) {
	crops := rotationTestCrops()
	duration := 120
	crop := crops[1]
	crop.WaterRequirement = "low"
	crop.NutrientDemand = "low"
	crop.GrowthDuration = &duration

	reason, benefits := explainCrop(
		&crop,
		map[string]int{},
		nil,
		[]string{"crop_diversity", "water_efficiency"},
		waterConsciousRotation,
	)
	for _, expected := range []string{
		"It adds a crop family not yet used in this sequence.",
		"This also advances the farmer-stated crop-diversity priority.",
		"This crop is not present in the recent farm history considered.",
		"The catalogue records nitrogen-fixing potential; actual fixation depends on field conditions.",
		"It is recorded as having low water requirements; this does not mean it is drought-proof or suitable for every dry region.",
		"approximately 120 days to maturity or first harvest;",
	} {
		if !strings.Contains(reason, expected) {
			t.Errorf("crop reason %q does not contain %q", reason, expected)
		}
	}

	_, soilBenefits := explainCrop(
		&crop,
		map[string]int{},
		nil,
		nil,
		soilRestorationRotation,
	)
	soilBenefitText := strings.Join(soilBenefits, " ")
	if !strings.Contains(soilBenefitText, "Nitrogen-fixing potential is weighted by this strategy, not a claim of measured soil improvement.") {
		t.Errorf("soil benefits %q do not preserve the nitrogen-fixing explanation spacing", soilBenefitText)
	}

	attributeCases := []struct {
		name       string
		strategy   rotationStrategy
		water      string
		nutrients  string
		wantReason string
	}{
		{
			name:       "low nutrient demand",
			strategy:   soilRestorationRotation,
			water:      "medium",
			nutrients:  "low",
			wantReason: "Its recorded nutrient demand is low, which this strategy modestly favors.",
		},
		{
			name:       "medium nutrient demand",
			strategy:   soilRestorationRotation,
			water:      "medium",
			nutrients:  "medium",
			wantReason: "Its recorded nutrient demand is medium.",
		},
		{
			name:       "high nutrient demand",
			strategy:   soilRestorationRotation,
			water:      "medium",
			nutrients:  "high",
			wantReason: "Its recorded nutrient demand is high, so this strategy does not favor it on that attribute.",
		},
		{
			name:       "medium water requirement",
			strategy:   waterConsciousRotation,
			water:      "medium",
			nutrients:  "medium",
			wantReason: "Its recorded water requirement is medium.",
		},
		{
			name:       "high water requirement",
			strategy:   waterConsciousRotation,
			water:      "high",
			nutrients:  "medium",
			wantReason: "Its recorded water requirement is high, so this strategy does not favor it on that attribute.",
		},
	}
	for _, test := range attributeCases {
		t.Run(test.name, func(t *testing.T) {
			candidate := crop
			candidate.NitrogenFixing = false
			candidate.WaterRequirement = test.water
			candidate.NutrientDemand = test.nutrients
			candidate.GrowthDuration = nil
			reason, _ := explainCrop(
				&candidate,
				map[string]int{"fabaceae": 1},
				[]cropHistory{{CropID: candidate.ID, Family: candidate.Family}},
				nil,
				test.strategy,
			)
			if !strings.Contains(reason, test.wantReason) {
				t.Errorf("reason %q does not contain %q", reason, test.wantReason)
			}
		})
	}

	fallbackCrop := cropCandidate{
		ID:     uuid.MustParse("00000000-0000-0000-0000-000000000099"),
		Family: "Poaceae",
	}
	fallback, _ := explainCrop(
		&fallbackCrop,
		map[string]int{"poaceae": 1},
		[]cropHistory{{CropID: fallbackCrop.ID, Family: fallbackCrop.Family}},
		nil,
		balancedRotation,
	)
	if !strings.Contains(fallback, "Selected using family diversity, recent crop history, farmer priorities, and available catalogue data.") {
		t.Errorf("fallback reason %q does not preserve its complete sentence", fallback)
	}

	benefitText := strings.Join(benefits, " ")
	for _, expected := range []string{
		"Supports crop-family diversification.",
		"Aligned with farmer preferences.",
		"Low recorded water requirement is favored, not a guarantee of water suitability.",
	} {
		if !strings.Contains(benefitText, expected) {
			t.Errorf("crop benefits %q do not contain %q", benefitText, expected)
		}
	}

	results, err := buildRotationOptions(crops, nil, 3, nil)
	if err != nil {
		t.Fatalf("buildRotationOptions() error = %v", err)
	}
	for _, result := range results {
		var explanation []string
		if err := json.Unmarshal(result.Explanation, &explanation); err != nil {
			t.Fatalf("decode %s explanation: %v", result.PlanName, err)
		}

		text := strings.Join(explanation, " ")
		for _, expected := range []string{
			"Generated using deterministic, rule-based crop selection and catalogue records.",
			"All scores are provisional heuristics, not measurements or predictions.",
			"Diversity score is the share of the sequence represented by distinct recorded crop families.",
			"Soil health proxy starts at 40 and adds up to 35 points for recorded nitrogen-fixing share and 25 for low recorded nutrient-demand share.",
			"Climate score is fixed at 50 because climate observations are not used; it is not a climate-suitability score.",
			"Soil health score is a catalogue-based proxy, not a measured soil condition or guaranteed restoration outcome.",
			"Overall score weights the soil proxy at 40%, the climate placeholder at 20%, and diversity at 40%.",
			"Unknown crop characteristics are not treated as low water or nutrient demand.",
			"Review recommendations with a local agricultural extension officer.",
		} {
			if !strings.Contains(text, expected) {
				t.Errorf("%s explanation %q does not contain %q", result.PlanName, text, expected)
			}
		}

		for _, malformed := range []string{
			"measurementsor",
			"pointsfor",
			"thissequence",
			"cropis",
			"potentialis",
			"onfield",
			"at40%",
			"meanit",
		} {
			allText := text + " " + reason + " " + soilBenefitText + " " + benefitText
			if strings.Contains(allText, malformed) {
				t.Errorf("generated text contains malformed word boundary %q", malformed)
			}
		}
	}
}

func TestBuildRotationOptionsRejectsInsufficientCatalogue(t *testing.T) {
	_, err := buildRotationOptions(rotationTestCrops()[:2], nil, 3, nil)
	if !errors.Is(err, errInsufficientCrops) {
		t.Fatalf("buildRotationOptions() error = %v, want insufficient-crops error", err)
	}
}

func TestGenerateRotationRejectsEmptyCatalogue(t *testing.T) {
	repo := &fakeRotationRepository{farmExists: true}
	service := NewService(repo)

	_, _, err := service.GenerateRotation(context.Background(), CreateRotationRequest{
		FarmID:          uuid.NewString(),
		NumberOfSeasons: 1,
	})
	if err == nil || err.Error() != "no crops available for rotation planning" {
		t.Fatalf("GenerateRotation() error = %v, want empty-catalogue error", err)
	}
	if repo.planCreated {
		t.Fatal("GenerateRotation() created a plan for an empty catalogue")
	}
}

func TestGenerateRotationRejectsMoreSeasonsThanCrops(t *testing.T) {
	repo := &fakeRotationRepository{
		farmExists: true,
		crops:      rotationTestCrops()[:2],
	}
	service := NewService(repo)

	_, _, err := service.GenerateRotation(context.Background(), CreateRotationRequest{
		FarmID:          uuid.NewString(),
		NumberOfSeasons: 3,
	})
	if !errors.Is(err, errInsufficientCrops) {
		t.Fatalf("GenerateRotation() error = %v, want insufficient-crops error", err)
	}
	if repo.planCreated {
		t.Fatal("GenerateRotation() created a plan for an incomplete sequence")
	}
}

func TestCalculateRotationScoreUsesCatalogueHeuristics(t *testing.T) {
	crops := rotationTestCrops()
	sequence := []CropRecommendation{
		{CropID: crops[1].ID, Season: 1},
		{CropID: crops[2].ID, Season: 2},
	}

	score, err := calculateRotationScore(sequence, crops)
	if err != nil {
		t.Fatalf("calculateRotationScore() error = %v", err)
	}
	if score.SoilHealthScore != 70 {
		t.Errorf("soil score = %v, want 70", score.SoilHealthScore)
	}
	if score.ClimateScore != 50 {
		t.Errorf("climate score = %v, want fixed provisional value 50", score.ClimateScore)
	}
	if score.DiversityScore != 100 {
		t.Errorf("diversity score = %v, want 100", score.DiversityScore)
	}
	if score.OverallScore != 78 {
		t.Errorf("overall score = %v, want 78", score.OverallScore)
	}
}

func TestCalculateRotationScoreRejectsEmptyAndInvalidSequences(t *testing.T) {
	crops := rotationTestCrops()
	tests := []struct {
		name     string
		sequence []CropRecommendation
	}{
		{name: "empty"},
		{
			name: "unknown crop",
			sequence: []CropRecommendation{{
				CropID: uuid.New(),
				Season: 1,
			}},
		},
		{
			name: "duplicate crop",
			sequence: []CropRecommendation{
				{CropID: crops[0].ID, Season: 1},
				{CropID: crops[0].ID, Season: 2},
			},
		},
		{
			name: "invalid season order",
			sequence: []CropRecommendation{{
				CropID: crops[0].ID,
				Season: 2,
			}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if _, err := calculateRotationScore(test.sequence, crops); err == nil {
				t.Fatal("calculateRotationScore() error = nil, want validation error")
			}
		})
	}
}

type fakeRotationRepository struct {
	farmExists  bool
	crops       []cropCandidate
	history     []cropHistory
	planCreated bool
}

func (r *fakeRotationRepository) FarmExists(context.Context, uuid.UUID) (bool, error) {
	return r.farmExists, nil
}

func (r *fakeRotationRepository) ListAvailableCrops(context.Context) ([]cropCandidate, error) {
	return r.crops, nil
}

func (r *fakeRotationRepository) ListCropHistory(context.Context, uuid.UUID) ([]cropHistory, error) {
	return r.history, nil
}

func (r *fakeRotationRepository) CreatePlan(_ context.Context, plan *RotationPlan) error {
	r.planCreated = true
	plan.ID = uuid.New()
	return nil
}

func (*fakeRotationRepository) SaveResults(context.Context, []RotationResult) error {
	return nil
}

func (*fakeRotationRepository) GetResultsByPlan(context.Context, uuid.UUID) ([]RotationResult, error) {
	return nil, nil
}

func (*fakeRotationRepository) GetPlanByID(context.Context, uuid.UUID) (*RotationPlan, error) {
	return nil, nil
}

func (*fakeRotationRepository) ListPlansByFarm(context.Context, uuid.UUID, int, int) ([]RotationPlan, error) {
	return nil, nil
}

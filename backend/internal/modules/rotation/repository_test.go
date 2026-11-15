package rotation

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestReadAvailableCropsScansRowsAndMissingCharacteristics(t *testing.T) {
	cropID := uuid.MustParse("00000000-0000-0000-0000-000000000010")
	growthDuration := 120
	queryer := &fakeAvailableCropQueryer{
		rows: &fakeAvailableCropRows{
			records: [][]any{
				{
					cropID,
					"Maize",
					"Poaceae",
					"cereal",
					false,
					"medium",
					"high",
					growthDuration,
				},
				{
					uuid.MustParse("00000000-0000-0000-0000-000000000011"),
					"Uncharacterized crop",
					"",
					"",
					false,
					"",
					"",
					nil,
				},
			},
		},
	}

	crops, err := readAvailableCrops(context.Background(), queryer)
	if err != nil {
		t.Fatalf("readAvailableCrops() error = %v", err)
	}
	if len(crops) != 2 {
		t.Fatalf("readAvailableCrops() returned %d crops, want 2", len(crops))
	}
	if crops[0].ID != cropID || crops[0].Name != "Maize" ||
		crops[0].GrowthDuration == nil || *crops[0].GrowthDuration != growthDuration {
		t.Errorf("first crop was not scanned correctly: %#v", crops[0])
	}
	if crops[1].Name != "Uncharacterized crop" ||
		crops[1].NitrogenFixing ||
		crops[1].WaterRequirement != "" ||
		crops[1].NutrientDemand != "" ||
		crops[1].GrowthDuration != nil {
		t.Errorf("missing characteristics were not kept neutral: %#v", crops[1])
	}
	reason, _ := explainCrop(&crops[1], map[string]int{}, nil, nil, balancedRotation)
	if !strings.Contains(reason, "unknown values were left neutral") ||
		strings.Contains(reason, "low water requirements") {
		t.Errorf("missing characteristics were not explained neutrally: %q", reason)
	}
	if !strings.Contains(queryer.query, "LEFT JOIN crop_characteristics") {
		t.Errorf("available-crops query does not retain crops without characteristics: %q", queryer.query)
	}
	if !queryer.rows.closed {
		t.Error("available-crops rows were not closed")
	}
}

func TestReadAvailableCropsReturnsEmptyCatalogue(t *testing.T) {
	queryer := &fakeAvailableCropQueryer{rows: &fakeAvailableCropRows{}}

	crops, err := readAvailableCrops(context.Background(), queryer)
	if err != nil {
		t.Fatalf("readAvailableCrops() error = %v", err)
	}
	if crops == nil || len(crops) != 0 {
		t.Fatalf("readAvailableCrops() = %#v, want an empty non-nil slice", crops)
	}
}

func TestReadAvailableCropsReportsQueryError(t *testing.T) {
	queryer := &fakeAvailableCropQueryer{err: errors.New("query failed")}

	_, err := readAvailableCrops(context.Background(), queryer)
	if err == nil || !strings.Contains(err.Error(), "retrieve available crops") {
		t.Fatalf("readAvailableCrops() error = %v, want wrapped query error", err)
	}
}

type fakeAvailableCropQueryer struct {
	query string
	rows  *fakeAvailableCropRows
	err   error
}

func (q *fakeAvailableCropQueryer) Query(
	_ context.Context,
	query string,
	_ ...any,
) (availableCropRows, error) {
	q.query = query
	if q.err != nil {
		return nil, q.err
	}

	return q.rows, nil
}

type fakeAvailableCropRows struct {
	records [][]any
	index   int
	closed  bool
	err     error
}

func (r *fakeAvailableCropRows) Close() {
	r.closed = true
}

func (r *fakeAvailableCropRows) Next() bool {
	if r.index >= len(r.records) {
		return false
	}

	r.index++
	return true
}

func (r *fakeAvailableCropRows) Scan(destinations ...any) error {
	if r.index == 0 || r.index > len(r.records) {
		return errors.New("Scan called without a current row")
	}

	record := r.records[r.index-1]
	if len(record) != len(destinations) {
		return errors.New("unexpected number of scan destinations")
	}

	for i, destination := range destinations {
		value := record[i]
		switch target := destination.(type) {
		case *uuid.UUID:
			*target = value.(uuid.UUID)
		case *string:
			*target = value.(string)
		case *bool:
			*target = value.(bool)
		case **int:
			if value == nil {
				*target = nil
			} else {
				duration := value.(int)
				*target = &duration
			}
		default:
			return errors.New("unexpected scan destination type")
		}
	}

	return nil
}

func (r *fakeAvailableCropRows) Err() error {
	return r.err
}

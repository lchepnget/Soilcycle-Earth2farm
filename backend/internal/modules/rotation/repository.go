package rotation

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

// FarmExists checks whether a farm exists.
func (r *Repository) FarmExists(
	ctx context.Context,
	farmID uuid.UUID,
) (bool, error) {
	var exists bool

	err := r.db.QueryRow(
		ctx,
		`SELECT EXISTS (
			SELECT 1 FROM farms WHERE id = $1
		)`,
		farmID,
	).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("check farm existence: %w", err)
	}

	return exists, nil
}

// CreatePlan saves a new rotation plan.
func (r *Repository) CreatePlan(
	ctx context.Context,
	plan *RotationPlan,
) error {
	if plan == nil {
		return errors.New("rotation plan cannot be nil")
	}

	if plan.NumberOfSeasons < 1 || plan.NumberOfSeasons > 10 {
		return errors.New("number of seasons must be between 1 and 10")
	}

	priorities := normalizeJSON(plan.FarmerPriorities, "[]")
	snapshot := normalizeJSON(plan.InputSnapshot, "{}")

	err := r.db.QueryRow(
		ctx,
		`INSERT INTO rotation_plans (
			farm_id,
			farmer_priorities,
			input_snapshot,
			number_of_seasons
		)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`,
		plan.FarmID,
		priorities,
		snapshot,
		plan.NumberOfSeasons,
	).Scan(
		&plan.ID,
		&plan.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("create rotation plan: %w", err)
	}

	plan.FarmerPriorities = priorities
	plan.InputSnapshot = snapshot

	return nil
}

// GetPlanByID retrieves a rotation plan by its ID.
func (r *Repository) GetPlanByID(
	ctx context.Context,
	planID uuid.UUID,
) (*RotationPlan, error) {
	plan := &RotationPlan{}

	err := r.db.QueryRow(
		ctx,
		`SELECT
			id,
			farm_id,
			farmer_priorities,
			input_snapshot,
			number_of_seasons,
			created_at
		FROM rotation_plans
		WHERE id = $1`,
		planID,
	).Scan(
		&plan.ID,
		&plan.FarmID,
		&plan.FarmerPriorities,
		&plan.InputSnapshot,
		&plan.NumberOfSeasons,
		&plan.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("get rotation plan: %w", err)
	}

	return plan, nil
}

// ListPlansByFarm retrieves rotation plans belonging to a farm.
func (r *Repository) ListPlansByFarm(
	ctx context.Context,
	farmID uuid.UUID,
	limit int,
	offset int,
) ([]RotationPlan, error) {
	if limit < 1 || limit > 100 {
		limit = 20
	}

	if offset < 0 {
		offset = 0
	}

	rows, err := r.db.Query(
		ctx,
		`SELECT
			id,
			farm_id,
			farmer_priorities,
			input_snapshot,
			number_of_seasons,
			created_at
		FROM rotation_plans
		WHERE farm_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`,
		farmID,
		limit,
		offset,
	)

	if err != nil {
		return nil, fmt.Errorf("list rotation plans: %w", err)
	}

	defer rows.Close()

	plans := make([]RotationPlan, 0)

	for rows.Next() {
		var plan RotationPlan

		err := rows.Scan(
			&plan.ID,
			&plan.FarmID,
			&plan.FarmerPriorities,
			&plan.InputSnapshot,
			&plan.NumberOfSeasons,
			&plan.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("scan rotation plan: %w", err)
		}

		plans = append(plans, plan)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate rotation plans: %w", err)
	}

	return plans, nil
}

// SaveResults saves generated rotation results.
// All results are inserted in one transaction.
func (r *Repository) SaveResults(
	ctx context.Context,
	results []RotationResult,
) error {
	if len(results) == 0 {
		return nil
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin rotation results transaction: %w", err)
	}

	defer tx.Rollback(ctx)

	for _, result := range results {
		cropSequence := normalizeJSON(result.CropSequence, "[]")
		scores := normalizeJSON(result.Scores, "{}")
		explanation := normalizeJSON(result.Explanation, "[]")

		_, err := tx.Exec(
			ctx,
			`INSERT INTO rotation_results (
				plan_id,
				plan_name,
				crop_sequence,
				scores,
				explanation,
				overall_score
			)
			VALUES ($1, $2, $3, $4, $5, $6)`,
			result.PlanID,
			result.PlanName,
			cropSequence,
			scores,
			explanation,
			result.OverallScore,
		)

		if err != nil {
			return fmt.Errorf("save rotation result: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit rotation results: %w", err)
	}

	return nil
}

// GetResultsByPlan retrieves all results belonging to a rotation plan.
func (r *Repository) GetResultsByPlan(
	ctx context.Context,
	planID uuid.UUID,
) ([]RotationResult, error) {
	rows, err := r.db.Query(
		ctx,
		`SELECT
			id,
			plan_id,
			plan_name,
			crop_sequence,
			scores,
			explanation,
			overall_score,
			created_at
		FROM rotation_results
		WHERE plan_id = $1
		ORDER BY overall_score DESC NULLS LAST, created_at DESC`,
		planID,
	)

	if err != nil {
		return nil, fmt.Errorf("get rotation results: %w", err)
	}

	defer rows.Close()

	results := make([]RotationResult, 0)

	for rows.Next() {
		var result RotationResult

		err := rows.Scan(
			&result.ID,
			&result.PlanID,
			&result.PlanName,
			&result.CropSequence,
			&result.Scores,
			&result.Explanation,
			&result.OverallScore,
			&result.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("scan rotation result: %w", err)
		}

		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate rotation results: %w", err)
	}

	return results, nil
}

// normalizeJSON provides a valid JSON default for empty JSON fields.
func normalizeJSON(value json.RawMessage, fallback string) json.RawMessage {
	if len(value) == 0 || !json.Valid(value) {
		return json.RawMessage(fallback)
	}

	return value
}

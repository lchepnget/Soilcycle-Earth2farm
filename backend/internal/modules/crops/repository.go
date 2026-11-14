package crops

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CropRepository struct {
	db *pgxpool.Pool
}

func NewCropRepository(db *pgxpool.Pool) *CropRepository {
	return &CropRepository{db: db}
}

func (r *CropRepository) Create(ctx context.Context, req CreateCropRequest) (*Crop, error) {
	query := `
		INSERT INTO crops (
			name, scientific_name, crop_family, category, description
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, name, scientific_name, crop_family,
		          category, description, created_at
	`

	var crop Crop

	err := r.db.QueryRow(
		ctx,
		query,
		req.Name,
		req.ScientificName,
		req.CropFamily,
		req.Category,
		req.Description,
	).Scan(
		&crop.ID,
		&crop.Name,
		&crop.ScientificName,
		&crop.CropFamily,
		&crop.Category,
		&crop.Description,
		&crop.CreatedAt,
	)

	if err != nil {
		return nil, mapDatabaseError(err)
	}

	return &crop, nil
}

func (r *CropRepository) GetByID(ctx context.Context, id uuid.UUID) (*Crop, error) {
	query := `
		SELECT id, name, scientific_name, crop_family,
		       category, description, created_at
		FROM crops
		WHERE id = $1
	`

	var crop Crop

	err := r.db.QueryRow(ctx, query, id).Scan(
		&crop.ID,
		&crop.Name,
		&crop.ScientificName,
		&crop.CropFamily,
		&crop.Category,
		&crop.Description,
		&crop.CreatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrCropNotFound
	}

	if err != nil {
		return nil, err
	}

	return &crop, nil
}

func (r *CropRepository) List(
	ctx context.Context,
	search, category string,
	limit, offset int,
) ([]Crop, int64, error) {

	query := `
		SELECT id, name, scientific_name, crop_family,
		       category, description, created_at
		FROM crops
		WHERE ($1 = '' OR name ILIKE '%' || $1 || '%')
		  AND ($2 = '' OR category = $2)
		ORDER BY name ASC
		LIMIT $3 OFFSET $4
	`

	countQuery := `
		SELECT COUNT(*)
		FROM crops
		WHERE ($1 = '' OR name ILIKE '%' || $1 || '%')
		  AND ($2 = '' OR category = $2)
	`

	var total int64

	err := r.db.QueryRow(ctx, countQuery, search, category).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(ctx, query, search, category, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	cropsList := make([]Crop, 0)

	for rows.Next() {
		var crop Crop

		err := rows.Scan(
			&crop.ID,
			&crop.Name,
			&crop.ScientificName,
			&crop.CropFamily,
			&crop.Category,
			&crop.Description,
			&crop.CreatedAt,
		)

		if err != nil {
			return nil, 0, err
		}

		cropsList = append(cropsList, crop)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return cropsList, total, nil
}

func (r *CropRepository) Update(
	ctx context.Context,
	id uuid.UUID,
	req UpdateCropRequest,
) (*Crop, error) {

	query := `
		UPDATE crops
		SET name = $1,
		    scientific_name = $2,
		    crop_family = $3,
		    category = $4,
		    description = $5
		WHERE id = $6
		RETURNING id, name, scientific_name, crop_family,
		          category, description, created_at
	`

	var crop Crop

	err := r.db.QueryRow(
		ctx,
		query,
		req.Name,
		req.ScientificName,
		req.CropFamily,
		req.Category,
		req.Description,
		id,
	).Scan(
		&crop.ID,
		&crop.Name,
		&crop.ScientificName,
		&crop.CropFamily,
		&crop.Category,
		&crop.Description,
		&crop.CreatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrCropNotFound
	}

	if err != nil {
		return nil, mapDatabaseError(err)
	}

	return &crop, nil
}

func (r *CropRepository) Delete(ctx context.Context, id uuid.UUID) error {
	tag, err := r.db.Exec(
		ctx,
		`DELETE FROM crops WHERE id = $1`,
		id,
	)

	if err != nil {
		return mapDatabaseError(err)
	}

	if tag.RowsAffected() == 0 {
		return ErrCropNotFound
	}

	return nil
}

func (r *CropRepository) GetCharacteristics(
	ctx context.Context,
	cropID uuid.UUID,
) (*CropCharacteristic, error) {

	query := `
		SELECT id, crop_id, water_requirement,
		       growth_duration_days, nitrogen_fixing,
		       nutrient_demand, suitable_seasons, created_at
		FROM crop_characteristics
		WHERE crop_id = $1
	`

	var characteristic CropCharacteristic

	err := r.db.QueryRow(ctx, query, cropID).Scan(
		&characteristic.ID,
		&characteristic.CropID,
		&characteristic.WaterRequirement,
		&characteristic.GrowthDurationDays,
		&characteristic.NitrogenFixing,
		&characteristic.NutrientDemand,
		&characteristic.SuitableSeasons,
		&characteristic.CreatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrCropNotFound
	}

	if err != nil {
		return nil, err
	}

	return &characteristic, nil
}

func (r *CropRepository) UpsertCharacteristics(
	ctx context.Context,
	cropID uuid.UUID,
	req CropCharacteristicRequest,
) (*CropCharacteristic, error) {

	query := `
		INSERT INTO crop_characteristics (
			crop_id, water_requirement, growth_duration_days,
			nitrogen_fixing, nutrient_demand, suitable_seasons
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (crop_id)
		DO UPDATE SET
			water_requirement = EXCLUDED.water_requirement,
			growth_duration_days = EXCLUDED.growth_duration_days,
			nitrogen_fixing = EXCLUDED.nitrogen_fixing,
			nutrient_demand = EXCLUDED.nutrient_demand,
			suitable_seasons = EXCLUDED.suitable_seasons
		RETURNING id, crop_id, water_requirement,
		          growth_duration_days, nitrogen_fixing,
		          nutrient_demand, suitable_seasons, created_at
	`

	var characteristic CropCharacteristic

	err := r.db.QueryRow(
		ctx,
		query,
		cropID,
		req.WaterRequirement,
		req.GrowthDurationDays,
		req.NitrogenFixing,
		req.NutrientDemand,
		req.SuitableSeasons,
	).Scan(
		&characteristic.ID,
		&characteristic.CropID,
		&characteristic.WaterRequirement,
		&characteristic.GrowthDurationDays,
		&characteristic.NitrogenFixing,
		&characteristic.NutrientDemand,
		&characteristic.SuitableSeasons,
		&characteristic.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &characteristic, nil
}

func mapDatabaseError(err error) error {
	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return ErrDuplicateCrop
		case "23503":
			return fmt.Errorf("referenced record does not exist: %w", err)
		}
	}

	return err
}

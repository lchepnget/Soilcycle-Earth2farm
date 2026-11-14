package soil

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrSoilNotFound = errors.New("soil record not found")
	ErrFarmNotFound = errors.New("farm not found")
)

type SoilRepository struct {
	db *pgxpool.Pool
}

func NewSoilRepository(db *pgxpool.Pool) *SoilRepository {
	return &SoilRepository{db: db}
}

// Converts an optional YYYY-MM-DD string into a nullable time value.
func parseMeasurementDate(value *string) (*time.Time, error) {
	if value == nil {
		return nil, nil
	}

	parsed, err := time.Parse("2006-01-02", *value)
	if err != nil {
		return nil, fmt.Errorf("invalid measurement date: %w", err)
	}

	return &parsed, nil
}

// Checks whether the farm exists.
func (r *SoilRepository) farmExists(
	ctx context.Context,
	farmID uuid.UUID,
) (bool, error) {
	var exists bool

	err := r.db.QueryRow(
		ctx,
		`SELECT EXISTS(SELECT 1 FROM farms WHERE id = $1)`,
		farmID,
	).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("check farm existence: %w", err)
	}

	return exists, nil
}

// Shared scanner for soil records.
func scanSoil(row pgx.Row) (*SoilData, error) {
	var soil SoilData

	err := row.Scan(
		&soil.ID,
		&soil.FarmID,
		&soil.PH,
		&soil.OrganicCarbon,
		&soil.Nitrogen,
		&soil.Phosphorus,
		&soil.Potassium,
		&soil.SoilTexture,
		&soil.DataSource,
		&soil.MeasurementDate,
		&soil.CreatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrSoilNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("scan soil record: %w", err)
	}

	return &soil, nil
}

const soilColumns = `
	id,
	farm_id,
	ph,
	organic_carbon,
	nitrogen,
	phosphorus,
	potassium,
	soil_texture,
	data_source,
	measurement_date,
	created_at
`

// Create a soil record.
func (r *SoilRepository) Create(
	ctx context.Context,
	req CreateSoilRequest,
) (*SoilData, error) {
	farmID, err := uuid.Parse(req.FarmID)
	if err != nil {
		return nil, errors.New("invalid farm ID")
	}

	exists, err := r.farmExists(ctx, farmID)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, ErrFarmNotFound
	}

	measurementDate, err := parseMeasurementDate(req.MeasurementDate)
	if err != nil {
		return nil, err
	}

	query := `
		INSERT INTO soil_data (
			farm_id,
			ph,
			organic_carbon,
			nitrogen,
			phosphorus,
			potassium,
			soil_texture,
			data_source,
			measurement_date
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING ` + soilColumns

	row := r.db.QueryRow(
		ctx,
		query,
		farmID,
		req.PH,
		req.OrganicCarbon,
		req.Nitrogen,
		req.Phosphorus,
		req.Potassium,
		req.SoilTexture,
		req.DataSource,
		measurementDate,
	)

	return scanSoil(row)
}

// Get a soil record by ID.
func (r *SoilRepository) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*SoilData, error) {
	query := `SELECT ` + soilColumns + `
		FROM soil_data
		WHERE id = $1`

	return scanSoil(r.db.QueryRow(ctx, query, id))
}

// List all soil records with pagination.
func (r *SoilRepository) List(
	ctx context.Context,
	limit, offset int,
) ([]SoilData, int64, error) {
	var total int64

	err := r.db.QueryRow(
		ctx,
		`SELECT COUNT(*) FROM soil_data`,
	).Scan(&total)

	if err != nil {
		return nil, 0, fmt.Errorf("count soil records: %w", err)
	}

	query := `
		SELECT ` + soilColumns + `
		FROM soil_data
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list soil records: %w", err)
	}
	defer rows.Close()

	records := make([]SoilData, 0)

	for rows.Next() {
		soil, err := scanSoil(rows)
		if err != nil {
			return nil, 0, err
		}

		records = append(records, *soil)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate soil records: %w", err)
	}

	return records, total, nil
}

// List soil records belonging to a specific farm.
func (r *SoilRepository) ListByFarm(
	ctx context.Context,
	farmID uuid.UUID,
	limit, offset int,
) ([]SoilData, int64, error) {
	exists, err := r.farmExists(ctx, farmID)
	if err != nil {
		return nil, 0, err
	}

	if !exists {
		return nil, 0, ErrFarmNotFound
	}

	var total int64

	err = r.db.QueryRow(
		ctx,
		`SELECT COUNT(*) FROM soil_data WHERE farm_id = $1`,
		farmID,
	).Scan(&total)

	if err != nil {
		return nil, 0, fmt.Errorf("count farm soil records: %w", err)
	}

	query := `
		SELECT ` + soilColumns + `
		FROM soil_data
		WHERE farm_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(
		ctx,
		query,
		farmID,
		limit,
		offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("list farm soil records: %w", err)
	}
	defer rows.Close()

	records := make([]SoilData, 0)

	for rows.Next() {
		soil, err := scanSoil(rows)
		if err != nil {
			return nil, 0, err
		}

		records = append(records, *soil)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate farm soil records: %w", err)
	}

	return records, total, nil
}

// Update a soil record.
func (r *SoilRepository) Update(
	ctx context.Context,
	id uuid.UUID,
	req UpdateSoilRequest,
) (*SoilData, error) {
	measurementDate, err := parseMeasurementDate(req.MeasurementDate)
	if err != nil {
		return nil, err
	}

	query := `
		UPDATE soil_data
		SET
			ph = $1,
			organic_carbon = $2,
			nitrogen = $3,
			phosphorus = $4,
			potassium = $5,
			soil_texture = $6,
			data_source = $7,
			measurement_date = $8
		WHERE id = $9
		RETURNING ` + soilColumns

	row := r.db.QueryRow(
		ctx,
		query,
		req.PH,
		req.OrganicCarbon,
		req.Nitrogen,
		req.Phosphorus,
		req.Potassium,
		req.SoilTexture,
		req.DataSource,
		measurementDate,
		id,
	)

	return scanSoil(row)
}

// Delete a soil record.
func (r *SoilRepository) Delete(
	ctx context.Context,
	id uuid.UUID,
) error {
	tag, err := r.db.Exec(
		ctx,
		`DELETE FROM soil_data WHERE id = $1`,
		id,
	)

	if err != nil {
		return fmt.Errorf("delete soil record: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return ErrSoilNotFound
	}

	return nil
}

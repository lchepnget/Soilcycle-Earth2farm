package farms

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrFarmNotFound  = errors.New("farm not found")
	ErrFarmerMissing = errors.New("farmer does not exist")
)

type FarmRepository struct {
	db *pgxpool.Pool
}

func NewFarmRepository(db *pgxpool.Pool) *FarmRepository {
	return &FarmRepository{db: db}
}

const farmColumns = `
	id, user_id, name, county, sub_county,
	latitude, longitude, area_hectares,
	water_availability, created_at
`

func scanFarm(row pgx.Row) (*Farm, error) {
	var farm Farm

	err := row.Scan(
		&farm.ID,
		&farm.UserID,
		&farm.Name,
		&farm.County,
		&farm.SubCounty,
		&farm.Latitude,
		&farm.Longitude,
		&farm.AreaHectares,
		&farm.WaterAvailability,
		&farm.CreatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrFarmNotFound
	}

	if err != nil {
		return nil, err
	}

	return &farm, nil
}

func (r *FarmRepository) Create(
	ctx context.Context,
	req CreateFarmRequest,
) (*Farm, error) {

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, err
	}

	// Ensure the owner exists and is a farmer.
	var exists bool

	err = r.db.QueryRow(
		ctx,
		`SELECT EXISTS(
			SELECT 1 FROM users WHERE id = $1 AND role = 'farmer'
		)`,
		userID,
	).Scan(&exists)

	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, ErrFarmerMissing
	}

	query := `
		INSERT INTO farms (
			user_id, name, county, sub_county,
			latitude, longitude, area_hectares,
			water_availability
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		RETURNING ` + farmColumns

	return scanFarm(r.db.QueryRow(
		ctx,
		query,
		userID,
		req.Name,
		req.County,
		req.SubCounty,
		req.Latitude,
		req.Longitude,
		req.AreaHectares,
		req.WaterAvailability,
	))
}

func (r *FarmRepository) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*Farm, error) {

	query := `
		SELECT ` + farmColumns + `
		FROM farms
		WHERE id = $1
	`

	return scanFarm(r.db.QueryRow(ctx, query, id))
}

func (r *FarmRepository) List(
	ctx context.Context,
	userID *uuid.UUID,
	limit, offset int,
) ([]Farm, int64, error) {

	var total int64

	countQuery := `
		SELECT COUNT(*)
		FROM farms
		WHERE ($1::uuid IS NULL OR user_id = $1)
	`

	err := r.db.QueryRow(ctx, countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `
		SELECT ` + farmColumns + `
		FROM farms
		WHERE ($1::uuid IS NULL OR user_id = $1)
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	result := make([]Farm, 0)

	for rows.Next() {
		var farm Farm

		err := rows.Scan(
			&farm.ID,
			&farm.UserID,
			&farm.Name,
			&farm.County,
			&farm.SubCounty,
			&farm.Latitude,
			&farm.Longitude,
			&farm.AreaHectares,
			&farm.WaterAvailability,
			&farm.CreatedAt,
		)

		if err != nil {
			return nil, 0, err
		}

		result = append(result, farm)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return result, total, nil
}

func (r *FarmRepository) Update(
	ctx context.Context,
	id, userID uuid.UUID,
	req UpdateFarmRequest,
) (*Farm, error) {

	query := `
		UPDATE farms
		SET name = $1,
		    county = $2,
		    sub_county = $3,
		    latitude = $4,
		    longitude = $5,
		    area_hectares = $6,
		    water_availability = $7
		WHERE id = $8 AND user_id = $9
		RETURNING ` + farmColumns

	farm, err := scanFarm(r.db.QueryRow(
		ctx,
		query,
		req.Name,
		req.County,
		req.SubCounty,
		req.Latitude,
		req.Longitude,
		req.AreaHectares,
		req.WaterAvailability,
		id,
		userID,
	))

	return farm, err
}

func (r *FarmRepository) Delete(
	ctx context.Context,
	id, userID uuid.UUID,
) error {

	tag, err := r.db.Exec(
		ctx,
		`DELETE FROM farms WHERE id = $1 AND user_id = $2`,
		id,
		userID,
	)

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return fmt.Errorf("farm has dependent records and cannot be deleted: %w", err)
		}

		return err
	}

	if tag.RowsAffected() == 0 {
		return ErrFarmNotFound
	}

	return nil
}

package farmers

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrFarmerNotFound = errors.New("farmer not found")
	ErrFarmerExists   = errors.New("email or phone already registered")
)

type FarmerRepository struct {
	db *pgxpool.Pool
}

func NewFarmerRepository(db *pgxpool.Pool) *FarmerRepository {
	return &FarmerRepository{db: db}
}

const farmerColumns = `
	id, full_name, email, phone,
	role, preferred_language, created_at, updated_at
`

func scanFarmer(row pgx.Row) (*Farmer, error) {
	var farmer Farmer

	err := row.Scan(
		&farmer.ID,
		&farmer.FullName,
		&farmer.Email,
		&farmer.Phone,
		&farmer.Role,
		&farmer.PreferredLanguage,
		&farmer.CreatedAt,
		&farmer.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrFarmerNotFound
	}

	if err != nil {
		return nil, err
	}

	return &farmer, nil
}

func (r *FarmerRepository) Create(
	ctx context.Context,
	req CreateFarmerRequest,
	passwordHash string,
) (*Farmer, error) {

	query := `
		INSERT INTO users (
			full_name, email, phone,
			password_hash, role, preferred_language
		)
		VALUES ($1, $2, $3, $4, 'farmer', $5)
		RETURNING ` + farmerColumns

	farmer, err := scanFarmer(r.db.QueryRow(
		ctx,
		query,
		req.FullName,
		req.Email,
		req.Phone,
		passwordHash,
		req.PreferredLanguage,
	))

	if err != nil {
		return nil, mapFarmerError(err)
	}

	return farmer, nil
}

func (r *FarmerRepository) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*Farmer, error) {

	query := `
		SELECT ` + farmerColumns + `
		FROM users
		WHERE id = $1 AND role = 'farmer'
	`

	return scanFarmer(r.db.QueryRow(ctx, query, id))
}

func (r *FarmerRepository) List(
	ctx context.Context,
	limit, offset int,
) ([]Farmer, int64, error) {

	var total int64

	err := r.db.QueryRow(
		ctx,
		`SELECT COUNT(*) FROM users WHERE role = 'farmer'`,
	).Scan(&total)

	if err != nil {
		return nil, 0, err
	}

	query := `
		SELECT ` + farmerColumns + `
		FROM users
		WHERE role = 'farmer'
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	result := make([]Farmer, 0)

	for rows.Next() {
		var farmer Farmer

		err := rows.Scan(
			&farmer.ID,
			&farmer.FullName,
			&farmer.Email,
			&farmer.Phone,
			&farmer.Role,
			&farmer.PreferredLanguage,
			&farmer.CreatedAt,
			&farmer.UpdatedAt,
		)

		if err != nil {
			return nil, 0, err
		}

		result = append(result, farmer)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return result, total, nil
}

func (r *FarmerRepository) Update(
	ctx context.Context,
	id uuid.UUID,
	req UpdateFarmerRequest,
) (*Farmer, error) {

	query := `
		UPDATE users
		SET full_name = $1,
		    email = $2,
		    phone = $3,
		    preferred_language = $4,
		    updated_at = NOW()
		WHERE id = $5 AND role = 'farmer'
		RETURNING ` + farmerColumns

	farmer, err := scanFarmer(r.db.QueryRow(
		ctx,
		query,
		req.FullName,
		req.Email,
		req.Phone,
		req.PreferredLanguage,
		id,
	))

	if err != nil {
		return nil, mapFarmerError(err)
	}

	return farmer, nil
}

func mapFarmerError(err error) error {
	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return ErrFarmerExists
	}

	return err
}

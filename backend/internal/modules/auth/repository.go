package auth

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrUserNotFound = errors.New("user not found")

type UserCredentials struct {
	ID           uuid.UUID
	PasswordHash string
	Role         string
}

type AuthRepository struct {
	db *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) FindByIdentifier(ctx context.Context, identifier string) (*UserCredentials, error) {
	var user UserCredentials
	var query string

	if strings.Contains(identifier, "@") {
		query = `SELECT id, COALESCE(password_hash, ''), role FROM users WHERE email = $1 LIMIT 1`
	} else {
		query = `SELECT id, COALESCE(password_hash, ''), role FROM users WHERE phone = $1 LIMIT 1`
	}

	err := r.db.QueryRow(ctx, query, identifier).Scan(&user.ID, &user.PasswordHash, &user.Role)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

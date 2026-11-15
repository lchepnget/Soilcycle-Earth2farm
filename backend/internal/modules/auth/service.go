package auth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrAuthNotConfigured  = errors.New("authentication is not configured")
)

const (
	tokenLifetime       = time.Hour
	minimumSecretLength = 32
	tokenIssuer         = "soilcycle-api"
)

type AuthService struct {
	repo      *AuthRepository
	jwtSecret []byte
}

func NewAuthService(repo *AuthRepository, secret string) *AuthService {
	return &AuthService{repo: repo, jwtSecret: []byte(secret)}
}

type LoginResult struct {
	AccessToken string    `json:"access_token"`
	TokenType   string    `json:"token_type"`
	ExpiresIn   int64     `json:"expires_in"`
	User        LoginUser `json:"user"`
}

type LoginUser struct {
	ID   string `json:"id"`
	Role string `json:"role"`
}

type tokenClaims struct {
	Subject   string `json:"sub"`
	Role      string `json:"role"`
	ExpiresAt int64  `json:"exp"`
	IssuedAt  int64  `json:"iat"`
	Issuer    string `json:"iss"`
}

func (s *AuthService) Login(ctx context.Context, req LoginRequest) (*LoginResult, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if len(s.jwtSecret) < minimumSecretLength {
		return nil, ErrAuthNotConfigured
	}

	user, err := s.repo.FindByIdentifier(ctx, req.Identifier)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}
	if user.PasswordHash == "" {
		return nil, ErrInvalidCredentials
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	now := time.Now().UTC()
	claims := tokenClaims{
		Subject:   user.ID.String(),
		Role:      user.Role,
		ExpiresAt: now.Add(tokenLifetime).Unix(),
		IssuedAt:  now.Unix(),
		Issuer:    tokenIssuer,
	}
	token, err := s.generateToken(claims)
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   int64(tokenLifetime.Seconds()),
		User:        LoginUser{ID: user.ID.String(), Role: user.Role},
	}, nil
}

func (s *AuthService) generateToken(claims tokenClaims) (string, error) {
	headerJSON, err := json.Marshal(map[string]string{"alg": "HS256", "typ": "JWT"})
	if err != nil {
		return "", err
	}
	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	unsigned := base64.RawURLEncoding.EncodeToString(headerJSON) + "." +
		base64.RawURLEncoding.EncodeToString(claimsJSON)

	mac := hmac.New(sha256.New, s.jwtSecret)
	_, _ = mac.Write([]byte(unsigned))
	signature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	return unsigned + "." + signature, nil
}

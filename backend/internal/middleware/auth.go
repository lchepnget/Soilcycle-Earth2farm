// Package middleware contains HTTP middleware shared by the API.
package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

const (
	UserIDLocal         = "authUserID"
	UserRoleLocal       = "authUserRole"
	minimumSecretLength = 32
)

// Claims is the minimal set of JWT claims accepted by this middleware.
type Claims struct {
	Subject   string `json:"sub"`
	Role      string `json:"role,omitempty"`
	ExpiresAt int64  `json:"exp"`
	NotBefore int64  `json:"nbf,omitempty"`
	IssuedAt  int64  `json:"iat,omitempty"`
	Issuer    string `json:"iss,omitempty"`
}

// RequireAuth validates an Authorization: Bearer <JWT> token.
// The signing secret is supplied by the application configuration.
func RequireAuth(jwtSecret string) fiber.Handler {
	secret := []byte(jwtSecret)
	secretReady := len(secret) >= minimumSecretLength

	return func(c fiber.Ctx) error {
		if isPublicEndpoint(c.Method(), c.Path()) {
			return c.Next()
		}

		if !secretReady {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"success": false,
				"message": "Authentication is not configured",
			})
		}

		header := strings.TrimSpace(c.Get(fiber.HeaderAuthorization))
		parts := strings.Fields(header)

		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			return unauthorized(c)
		}

		claims, err := validateToken(parts[1], secret, time.Now())
		if err != nil {
			return unauthorized(c)
		}

		c.Locals(UserIDLocal, claims.Subject)
		c.Locals(UserRoleLocal, claims.Role)

		return c.Next()
	}
}

func unauthorized(c fiber.Ctx) error {
	c.Set(fiber.HeaderWWWAuthenticate, `Bearer realm="soilcycle-api"`)

	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"success": false,
		"message": "A valid bearer token is required",
	})
}

// validateToken deliberately accepts only HS256 JWTs.
func validateToken(token string, secret []byte, now time.Time) (*Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		return nil, errors.New("malformed JWT")
	}

	headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, errors.New("invalid JWT header encoding")
	}

	var header struct {
		Algorithm string `json:"alg"`
		Type      string `json:"typ"`
	}

	if err := json.Unmarshal(headerBytes, &header); err != nil || header.Algorithm != "HS256" {
		return nil, errors.New("unsupported JWT algorithm")
	}

	signature, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return nil, errors.New("invalid JWT signature encoding")
	}

	mac := hmac.New(sha256.New, secret)
	_, _ = mac.Write([]byte(parts[0] + "." + parts[1]))
	expected := mac.Sum(nil)

	if len(signature) != len(expected) ||
		subtle.ConstantTimeCompare(signature, expected) != 1 {
		return nil, errors.New("invalid JWT signature")
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, errors.New("invalid JWT payload encoding")
	}

	var claims Claims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, errors.New("invalid JWT claims")
	}

	if _, err := uuid.Parse(claims.Subject); err != nil {
		return nil, errors.New("JWT subject must be a UUID")
	}

	nowUnix := now.Unix()

	if claims.ExpiresAt == 0 || nowUnix >= claims.ExpiresAt {
		return nil, errors.New("JWT is expired or has no expiry")
	}

	if claims.NotBefore != 0 && nowUnix < claims.NotBefore {
		return nil, errors.New("JWT is not active yet")
	}

	if claims.IssuedAt != 0 && claims.IssuedAt > nowUnix+60 {
		return nil, errors.New("JWT issued-at time is in the future")
	}

	if claims.Role != "" &&
		claims.Role != "farmer" &&
		claims.Role != "admin" &&
		claims.Role != "agronomist" {
		return nil, errors.New("JWT role is invalid")
	}

	return &claims, nil
}

// isPublicEndpoint lists routes that do not require a user token.
func isPublicEndpoint(method, path string) bool {
	method = strings.ToUpper(method)
	path = strings.TrimSuffix(path, "/")

	if method == fiber.MethodGet &&
		(path == "/api/v1/crops" || strings.HasPrefix(path, "/api/v1/crops/")) {
		return true
	}

	if method == fiber.MethodPost && path == "/api/v1/farmers" {
		return true
	}

	if (method == fiber.MethodPost || method == fiber.MethodGet) &&
		(path == "/api/v1/auth/login" || path == "/api/v1/auth/register") {
		return true
	}

	return false
}

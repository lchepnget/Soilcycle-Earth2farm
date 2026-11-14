package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"testing"
	"time"
)

func signTestToken(t *testing.T, secret []byte, header map[string]string, claims Claims) string {
	t.Helper()
	headerJSON, err := json.Marshal(header)
	if err != nil {
		t.Fatal(err)
	}
	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		t.Fatal(err)
	}
	unsigned := base64.RawURLEncoding.EncodeToString(headerJSON) + "." + base64.RawURLEncoding.EncodeToString(claimsJSON)
	mac := hmac.New(sha256.New, secret)
	_, _ = mac.Write([]byte(unsigned))
	return unsigned + "." + base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func TestValidateToken(t *testing.T) {
	secret := []byte("a-test-secret-that-is-long-enough-123")
	now := time.Unix(1_800_000_000, 0)
	validClaims := Claims{
		Subject:   "c3d6e8c8-73a6-4df3-ae16-fb8c6e52c6e7",
		Role:      "farmer",
		ExpiresAt: now.Add(time.Hour).Unix(),
		IssuedAt:  now.Unix(),
	}

	t.Run("valid token", func(t *testing.T) {
		token := signTestToken(t, secret, map[string]string{"alg": "HS256", "typ": "JWT"}, validClaims)
		got, err := validateToken(token, secret, now)
		if err != nil {
			t.Fatalf("expected valid token, got error: %v", err)
		}
		if got.Subject != validClaims.Subject || got.Role != validClaims.Role {
			t.Fatalf("unexpected claims: %#v", got)
		}
	})

	t.Run("expired token", func(t *testing.T) {
		claims := validClaims
		claims.ExpiresAt = now.Unix()
		token := signTestToken(t, secret, map[string]string{"alg": "HS256"}, claims)
		if _, err := validateToken(token, secret, now); err == nil {
			t.Fatal("expected expired token to be rejected")
		}
	})

	t.Run("unsupported algorithm", func(t *testing.T) {
		token := signTestToken(t, secret, map[string]string{"alg": "none"}, validClaims)
		if _, err := validateToken(token, secret, now); err == nil {
			t.Fatal("expected unsupported algorithm to be rejected")
		}
	})

	t.Run("wrong secret", func(t *testing.T) {
		token := signTestToken(t, secret, map[string]string{"alg": "HS256"}, validClaims)
		if _, err := validateToken(token, []byte("another-secret-that-is-long-enough"), now); err == nil {
			t.Fatal("expected invalid signature to be rejected")
		}
	})
}

func TestIsPublicEndpoint(t *testing.T) {
	tests := []struct {
		method string
		path   string
		want   bool
	}{
		{"GET", "/api/v1/crops", true},
		{"GET", "/api/v1/crops/abc/characteristics", true},
		{"POST", "/api/v1/farmers", true},
		{"GET", "/api/v1/farms", false},
		{"DELETE", "/api/v1/crops/abc", false},
		{"POST", "/api/v1/farms", false},
	}
	for _, tt := range tests {
		t.Run(tt.method+" "+tt.path, func(t *testing.T) {
			if got := isPublicEndpoint(tt.method, tt.path); got != tt.want {
				t.Fatalf("isPublicEndpoint() = %v, want %v", got, tt.want)
			}
		})
	}
}

package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// ErrInvalidToken is returned when a token is missing, malformed, or expired.
var ErrInvalidToken = errors.New("invalid or expired token")

// Claims are the JWT claims carried by the API. Subject = user id; Role = role name.
type Claims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

// Issuer mints and validates JWTs with a shared secret and TTL (config-driven).
type Issuer struct {
	secret []byte
	ttl    time.Duration
}

// NewIssuer builds an Issuer from the configured secret and token TTL.
func NewIssuer(secret string, ttl time.Duration) *Issuer {
	return &Issuer{secret: []byte(secret), ttl: ttl}
}

// TTL exposes the configured token lifetime (used for the login expiresIn field).
func (i *Issuer) TTL() time.Duration { return i.ttl }

// Issue returns a signed JWT for the given user id and role.
func (i *Issuer) Issue(userID uuid.UUID, role string) (string, error) {
	now := time.Now()
	claims := Claims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(i.ttl)),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(i.secret)
}

// Parse validates a token string and returns its claims, or ErrInvalidToken.
func (i *Issuer) Parse(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	tok, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return i.secret, nil
	})
	if err != nil || !tok.Valid {
		return nil, ErrInvalidToken
	}
	return claims, nil
}

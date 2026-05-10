package security

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var ErrInvalidToken = errors.New("invalid token")

const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

type Claims struct {
	UserID string
	Role   string
}

type claimsContextKey struct{}

func GenerateJWT(secretKey string, userID string, role string) (string, error) {
	now := time.Now().UTC()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"iat":     now.Unix(),
		"exp":     now.Add(24 * time.Hour).Unix(),
	})

	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign jwt: %w", err)
	}

	return signedToken, nil
}

func ParseJWT(secretKey string, tokenString string) (Claims, error) {
	parser := jwt.NewParser(
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithExpirationRequired(),
	)

	claims := jwt.MapClaims{}
	token, err := parser.ParseWithClaims(tokenString, claims, func(_ *jwt.Token) (any, error) {
		return []byte(secretKey), nil
	})
	if err != nil || !token.Valid {
		return Claims{}, ErrInvalidToken
	}

	userID, ok := claims["user_id"].(string)
	if !ok || userID == "" {
		return Claims{}, ErrInvalidToken
	}

	role, ok := claims["role"].(string)
	if !ok || role == "" {
		return Claims{}, ErrInvalidToken
	}

	return Claims{
		UserID: userID,
		Role:   role,
	}, nil
}

func WithClaims(ctx context.Context, claims Claims) context.Context {
	return context.WithValue(ctx, claimsContextKey{}, claims)
}

func ClaimsFromContext(ctx context.Context) (Claims, bool) {
	claims, ok := ctx.Value(claimsContextKey{}).(Claims)
	return claims, ok
}

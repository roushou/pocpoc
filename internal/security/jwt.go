package security

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func NewJWT(claims jwt.Claims, secretKey string, duration time.Duration) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secretKey))
}

func ParseJWT(token, secretKey string) (jwt.Claims, error) {
	parser := jwt.NewParser(jwt.WithValidMethods([]string{"HS256"}))

	parsedToken, err := parser.Parse(token, func(t *jwt.Token) (any, error) {
		return []byte(secretKey), nil
	})
	if err != nil || !parsedToken.Valid {
		return nil, err
	}

	return parsedToken.Claims, nil
}

func ParseJWTWithClaims(tokenString string, claims jwt.Claims, secretKey string) (jwt.Claims, error) {
	parser := jwt.NewParser(jwt.WithValidMethods([]string{"HS256"}))

	token, err := parser.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		return []byte(secretKey), nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}

	return claims, nil
}

package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func IssueToken(secret, userID string, ttl time.Duration) (string, time.Time, error) {
	now := time.Now()
	exp := now.Add(ttl)
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(exp),
		IssuedAt:  jwt.NewNumericDate(now),
	})
	s, err := t.SignedString([]byte(secret))
	if err != nil {
		return "", time.Time{}, err
	}
	return s, exp, nil
}

func ParseToken(secret, tokenString string) (userID string, err error) {
	claims := &jwt.RegisteredClaims{}
	t, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(_ *jwt.Token) (any, error) { return []byte(secret), nil },
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithLeeway(time.Minute),
	)
	if err != nil {
		return "", err
	}
	parsedClaims, ok := t.Claims.(*jwt.RegisteredClaims)
	if !ok || !t.Valid {
		return "", fmt.Errorf("invalid token")
	}
	if parsedClaims.Subject == "" {
		return "", fmt.Errorf("invalid token")
	}
	return parsedClaims.Subject, nil
}

package service

import (
	"strings"
	"time"

	"focus/account-cabinet/internal/auth"
	"focus/account-cabinet/internal/config"
	"focus/account-cabinet/internal/repository"
)

type AuthResult struct {
	UserID      string
	Email       string
	AccessToken string
	ExpiresAt   time.Time
}

func normalizeEmail(email string) string {
	return strings.TrimSpace(strings.ToLower(email))
}

func issueAuthResult(cfg config.Config, authUser *repository.User) (*AuthResult, error) {
	token, exp, err := auth.IssueToken(cfg.JWTSecret, authUser.ID, config.JWTExpiry)
	if err != nil {
		return nil, err
	}
	return &AuthResult{
		UserID:      authUser.ID,
		Email:       authUser.Email,
		AccessToken: token,
		ExpiresAt:   exp.UTC(),
	}, nil
}

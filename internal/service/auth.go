package service

import (
	"context"
	"errors"
	"focus/account-cabinet/internal/appclient"
	"focus/account-cabinet/internal/auth"
	"focus/account-cabinet/internal/config"
	"focus/account-cabinet/internal/repository"
	"strings"
)

type AuthService struct {
	cfg      config.Config
	users    *repository.UserRepository
	appUsers *appclient.Client
}

func NewAuthService(cfg config.Config, users *repository.UserRepository, appUsers *appclient.Client) *AuthService {
	return &AuthService{
		cfg:      cfg,
		users:    users,
		appUsers: appUsers,
	}
}

func (s *AuthService) Register(ctx context.Context, userName, email, password string) (*AuthResult, error) {
	hash, err := auth.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user, err := s.users.Create(ctx, repository.CreateUserParams{
		Email:        normalizeEmail(email),
		PasswordHash: hash,
	})
	if err != nil {
		return nil, err
	}

	err = s.appUsers.ProvisionUser(ctx, appclient.UserProvisionRequest{
		UserID:   user.ID,
		UserName: strings.TrimSpace(userName),
		Email:    user.Email,
	})
	if err != nil {
		_ = s.users.DeleteByID(ctx, user.ID)
		return nil, err
	}

	return issueAuthResult(s.cfg, user)
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*AuthResult, error) {
	user, err := s.users.GetByEmail(ctx, normalizeEmail(email))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}
	if !auth.CheckPassword(user.PasswordHash, password) {
		return nil, ErrInvalidCredentials
	}
	return issueAuthResult(s.cfg, user)
}

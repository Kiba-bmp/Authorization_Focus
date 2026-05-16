package service

import (
	"context"
	"errors"
	"fmt"
	"focus/account-cabinet/internal/appclient"
	"focus/account-cabinet/internal/auth"
	"focus/account-cabinet/internal/config"
	"focus/account-cabinet/internal/repository"
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

func (s *AuthService) Register(ctx context.Context, email, password string) (*AuthResult, error) {
	email = normalizeEmail(email)

	if err := validateRegisterInput(email, password); err != nil {
		return nil, err
	}

	hash, err := auth.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("не удалось подготовить пароль: %w", err)
	}

	user, err := s.users.Create(ctx, repository.CreateUserParams{
		Email:        email,
		PasswordHash: hash,
	})
	if err != nil {
		if errors.Is(err, repository.ErrEmailTaken) {
			return nil, conflict("Пользователь с такой почтой уже зарегистрирован.")
		}
		return nil, fmt.Errorf("не удалось создать пользователя auth: %w", err)
	}

	err = s.appUsers.ProvisionUser(ctx, appclient.UserProvisionRequest{
		UserID: user.ID,
		Email:  user.Email,
	})
	if err != nil {
		_ = s.users.DeleteByID(ctx, user.ID)
		var responseErr *appclient.ResponseError
		if errors.As(err, &responseErr) {
			return nil, &AppError{
				StatusCode: responseErr.StatusCode,
				Message:    responseErr.Message,
			}
		}
		return nil, fmt.Errorf("не удалось создать пользователя в backend: %w", err)
	}

	return issueAuthResult(s.cfg, user)
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*AuthResult, error) {
	email = normalizeEmail(email)
	if err := validateLoginInput(email, password); err != nil {
		return nil, err
	}

	user, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, unauthorized("Неверная почта или пароль.")
		}
		return nil, fmt.Errorf("не удалось загрузить пользователя: %w", err)
	}
	if !auth.CheckPassword(user.PasswordHash, password) {
		return nil, unauthorized("Неверная почта или пароль.")
	}
	return issueAuthResult(s.cfg, user)
}

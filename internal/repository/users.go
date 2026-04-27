package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, params CreateUserParams) (*User, error) {
	rec := userRecord{
		ID:           uuid.New(),
		Email:        normalizeEmail(params.Email),
		PasswordHash: params.PasswordHash,
	}
	if err := r.db.WithContext(ctx).Create(&rec).Error; err != nil {
		if isUniqueViolation(err) {
			return nil, ErrEmailTaken
		}
		return nil, fmt.Errorf("insert user: %w", err)
	}
	return toUser(rec), nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	var rec userRecord
	err := r.db.WithContext(ctx).
		Where("email = ?", normalizeEmail(email)).
		First(&rec).Error
	if isRecordNotFound(err) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return toUser(rec), nil
}

func (r *UserRepository) DeleteByID(ctx context.Context, id string) error {
	userID, err := parseUUID(id)
	if err != nil {
		return err
	}

	res := r.db.WithContext(ctx).Delete(&userRecord{}, "id = ?", userID)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

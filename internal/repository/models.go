package repository

import (
	"time"

	"github.com/google/uuid"
)

type userRecord struct {
	ID           uuid.UUID `gorm:"column:id;type:uuid;primaryKey"`
	Email        string    `gorm:"column:email;not null;uniqueIndex"`
	PasswordHash string    `gorm:"column:password_hash;not null"`
	CreatedAt    time.Time `gorm:"column:created_at;not null"`
	UpdatedAt    time.Time `gorm:"column:updated_at;not null"`
}

func (userRecord) TableName() string { return "users" }

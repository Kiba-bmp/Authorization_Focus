package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"focus/account-cabinet/internal/config"

	"github.com/pressly/goose/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Store struct {
	db    *gorm.DB
	sqlDB *sql.DB
}

func Open(cfg config.Config) (*Store, error) {
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть подключение к postgres: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("не удалось получить sql-подключение из gorm: %w", err)
	}

	if err := sqlDB.PingContext(context.Background()); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("не удалось проверить доступность postgres: %w", err)
	}

	if err := goose.SetDialect("postgres"); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("не удалось настроить dialect для goose: %w", err)
	}
	if err := goose.Up(sqlDB, config.MigrationsDir); err != nil && !errors.Is(err, goose.ErrNoNextVersion) {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("не удалось применить миграции: %w", err)
	}

	return &Store{db: db, sqlDB: sqlDB}, nil
}

func (s *Store) Close() error {
	if s == nil || s.sqlDB == nil {
		return nil
	}
	return s.sqlDB.Close()
}

func (s *Store) DB() *gorm.DB {
	if s == nil {
		return nil
	}
	return s.db
}

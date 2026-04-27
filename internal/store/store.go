package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"focus/account-cabinet/internal/config"

	"github.com/jackc/pgx/v5/pgconn"
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
		return nil, fmt.Errorf("open postgres: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("gorm db handle: %w", err)
	}

	if err := sqlDB.PingContext(context.Background()); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	if err := goose.SetDialect("postgres"); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("set goose dialect: %w", err)
	}
	if err := goose.Up(sqlDB, config.MigrationsDir); err != nil && !errors.Is(err, goose.ErrNoNextVersion) {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("apply migrations: %w", err)
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

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}

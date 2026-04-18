package store

import (
	"context"
	"database/sql"
	"time"

	_ "modernc.org/sqlite"
)

type Store struct {
	db *sql.DB
}

func Open(path string) (*Store, error) {
	dsn := "file:" + path + "?_pragma=foreign_keys(1)"
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	db.SetConnMaxLifetime(time.Hour)
	if err := migrate(db); err != nil {
		_ = db.Close()
		return nil, err
	}
	return &Store{db: db}, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

func migrate(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			user_name TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE,
			password_hash TEXT NOT NULL,
			email_verified INTEGER NOT NULL DEFAULT 0,
			invite_code TEXT NOT NULL,
			avatar_emoji TEXT,
			today_minutes INTEGER NOT NULL DEFAULT 0,
			week_minutes INTEGER NOT NULL DEFAULT 0,
			created_at INTEGER NOT NULL
		);
		CREATE TABLE IF NOT EXISTS email_codes (
			email TEXT NOT NULL PRIMARY KEY,
			code TEXT NOT NULL,
			expires_at INTEGER NOT NULL
		);
	`)
	return err
}

func (s *Store) SaveEmailCode(ctx context.Context, email, code string, expiresAt time.Time) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO email_codes (email, code, expires_at) VALUES (lower(?), ?, ?)
		ON CONFLICT(email) DO UPDATE SET code = excluded.code, expires_at = excluded.expires_at
	`, email, code, expiresAt.Unix())
	return err
}

func (s *Store) ConsumeEmailCode(ctx context.Context, email, code string) (bool, error) {
	row := s.db.QueryRowContext(ctx, `SELECT code, expires_at FROM email_codes WHERE lower(email) = lower(?)`, email)
	var stored string
	var exp int64
	if err := row.Scan(&stored, &exp); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	if time.Now().Unix() > exp {
		_, _ = s.db.ExecContext(ctx, `DELETE FROM email_codes WHERE lower(email) = lower(?)`, email)
		return false, nil
	}
	if stored != code {
		return false, nil
	}
	_, err := s.db.ExecContext(ctx, `DELETE FROM email_codes WHERE lower(email) = lower(?)`, email)
	if err != nil {
		return false, err
	}
	return true, nil
}

package store

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

var ErrNotFound = errors.New("not found")
var ErrEmailTaken = errors.New("email already registered")

type User struct {
	ID            string
	UserName      string
	Email         string
	PasswordHash  string
	EmailVerified bool
	InviteCode    string
	AvatarEmoji   *string
	TodayMinutes  int
	WeekMinutes   int
	CreatedAt     time.Time
}

func randomInviteCode() string {
	b := make([]byte, 4)
	_, _ = rand.Read(b)
	return strings.ToUpper(hex.EncodeToString(b))
}

func (s *Store) CreateUser(ctx context.Context, userName, email, passwordHash string) (*User, error) {
	id := uuid.NewString()
	invite := randomInviteCode()
	now := time.Now().Unix()
	email = strings.TrimSpace(strings.ToLower(email))
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO users (id, user_name, email, password_hash, email_verified, invite_code, avatar_emoji, today_minutes, week_minutes, created_at)
		VALUES (?, ?, ?, ?, 0, ?, NULL, 0, 0, ?)
	`, id, userName, email, passwordHash, invite, now)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return nil, ErrEmailTaken
		}
		return nil, fmt.Errorf("insert user: %w", err)
	}
	return &User{
		ID: id, UserName: userName, Email: email, PasswordHash: passwordHash,
		EmailVerified: false, InviteCode: invite, TodayMinutes: 0, WeekMinutes: 0,
		CreatedAt: time.Unix(now, 0),
	}, nil
}

func (s *Store) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	row := s.db.QueryRowContext(ctx, `
		SELECT id, user_name, email, password_hash, email_verified, invite_code, avatar_emoji, today_minutes, week_minutes, created_at
		FROM users WHERE lower(email) = ?
	`, email)
	return scanUser(row)
}

func (s *Store) GetUserByID(ctx context.Context, id string) (*User, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT id, user_name, email, password_hash, email_verified, invite_code, avatar_emoji, today_minutes, week_minutes, created_at
		FROM users WHERE id = ?
	`, id)
	return scanUser(row)
}

func scanUser(row *sql.Row) (*User, error) {
	var u User
	var verified int
	var avatar sql.NullString
	var created int64
	err := row.Scan(&u.ID, &u.UserName, &u.Email, &u.PasswordHash, &verified, &u.InviteCode, &avatar, &u.TodayMinutes, &u.WeekMinutes, &created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	u.EmailVerified = verified != 0
	if avatar.Valid {
		s := avatar.String
		u.AvatarEmoji = &s
	}
	u.CreatedAt = time.Unix(created, 0)
	return &u, nil
}

func (s *Store) SetEmailVerified(ctx context.Context, email string, verified bool) error {
	v := 0
	if verified {
		v = 1
	}
	email = strings.TrimSpace(strings.ToLower(email))
	res, err := s.db.ExecContext(ctx, `UPDATE users SET email_verified = ? WHERE lower(email) = ?`, v, email)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *Store) UpdateProfile(ctx context.Context, id string, userName *string, email *string, passwordHash *string, avatar *string) error {
	u, err := s.GetUserByID(ctx, id)
	if err != nil {
		return err
	}
	newName := u.UserName
	if userName != nil {
		newName = *userName
	}
	newEmail := u.Email
	if email != nil {
		newEmail = strings.TrimSpace(strings.ToLower(*email))
	}
	newHash := u.PasswordHash
	if passwordHash != nil {
		newHash = *passwordHash
	}
	var avatarVal interface{}
	if avatar != nil {
		if *avatar == "" {
			avatarVal = nil
		} else {
			avatarVal = *avatar
		}
	} else if u.AvatarEmoji != nil {
		avatarVal = *u.AvatarEmoji
	} else {
		avatarVal = nil
	}
	verified := u.EmailVerified
	if email != nil && newEmail != u.Email {
		verified = false
	}
	v := 0
	if verified {
		v = 1
	}
	_, err = s.db.ExecContext(ctx, `
		UPDATE users SET user_name = ?, email = ?, password_hash = ?, avatar_emoji = ?, email_verified = ?
		WHERE id = ?
	`, newName, newEmail, newHash, avatarVal, v, id)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return ErrEmailTaken
		}
		return err
	}
	return nil
}

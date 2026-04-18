package handlers

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"time"

	"focus/account-cabinet/internal/api"
	"focus/account-cabinet/internal/auth"
	"focus/account-cabinet/internal/store"

	"github.com/gin-gonic/gin"
)

func sixDigitCode() string {
	n, err := rand.Int(rand.Reader, big.NewInt(1_000_000))
	if err != nil {
		return "000000"
	}
	return fmt.Sprintf("%06d", n.Int64())
}

func (s *Server) authResponse(ctx context.Context, u *store.User) (api.AuthResponse, error) {
	token, exp, err := auth.IssueToken(s.cfg.JWTSecret, u.ID, s.cfg.JWTExpiry)
	if err != nil {
		return api.AuthResponse{}, err
	}
	return api.AuthResponse{
		UserID:      u.ID,
		UserName:    u.UserName,
		Email:       u.Email,
		InviteCode:  u.InviteCode,
		AccessToken: token,
		ExpiresAt:   exp.UTC().Format(time.RFC3339Nano),
	}, nil
}

func (s *Server) Register(c *gin.Context) {
	var req api.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, api.ErrorBody{Error: "invalid request body"})
		return
	}
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, api.ErrorBody{Error: "could not process password"})
		return
	}
	u, err := s.store.CreateUser(c.Request.Context(), strings.TrimSpace(req.UserName), req.Email, hash)
	if err != nil {
		if err == store.ErrEmailTaken {
			c.JSON(http.StatusConflict, api.ErrorBody{Error: "email already registered"})
			return
		}
		c.JSON(http.StatusInternalServerError, api.ErrorBody{Error: "registration failed"})
		return
	}
	if s.cfg.SkipEmailVerification {
		if err := s.store.SetEmailVerified(c.Request.Context(), u.Email, true); err != nil {
			c.JSON(http.StatusInternalServerError, api.ErrorBody{Error: "could not verify email flag"})
			return
		}
		u.EmailVerified = true
		resp, err := s.authResponse(c.Request.Context(), u)
		if err != nil {
			c.JSON(http.StatusInternalServerError, api.ErrorBody{Error: "could not issue token"})
			return
		}
		c.JSON(http.StatusOK, resp)
		return
	}
	code := sixDigitCode()
	expires := time.Now().Add(s.cfg.VerificationCodeTTL)
	if err := s.store.SaveEmailCode(c.Request.Context(), u.Email, code, expires); err != nil {
		c.JSON(http.StatusInternalServerError, api.ErrorBody{Error: "could not save verification code"})
		return
	}
	s.mail.SendVerificationCode(u.Email, code)
	c.JSON(http.StatusCreated, api.RegisterPendingResponse{
		VerificationRequired: true,
		Email:                u.Email,
		Message:              "Проверьте почту и введите код подтверждения",
	})
}

func (s *Server) Login(c *gin.Context) {
	var req api.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, api.ErrorBody{Error: "invalid request body"})
		return
	}
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	u, err := s.store.GetUserByEmail(c.Request.Context(), req.Email)
	if err != nil {
		if err == store.ErrNotFound {
			c.JSON(http.StatusUnauthorized, api.ErrorBody{Error: "invalid email or password"})
			return
		}
		c.JSON(http.StatusInternalServerError, api.ErrorBody{Error: "login failed"})
		return
	}
	if !auth.CheckPassword(u.PasswordHash, req.Password) {
		c.JSON(http.StatusUnauthorized, api.ErrorBody{Error: "invalid email or password"})
		return
	}
	if !u.EmailVerified {
		c.JSON(http.StatusForbidden, api.ErrorBody{Error: "email not verified"})
		return
	}
	resp, err := s.authResponse(c.Request.Context(), u)
	if err != nil {
		c.JSON(http.StatusInternalServerError, api.ErrorBody{Error: "could not issue token"})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (s *Server) VerifyEmail(c *gin.Context) {
	var req api.VerifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, api.ErrorBody{Error: "invalid request body"})
		return
	}
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	ok, err := s.store.ConsumeEmailCode(c.Request.Context(), req.Email, strings.TrimSpace(req.Code))
	if err != nil {
		c.JSON(http.StatusInternalServerError, api.ErrorBody{Error: "verification failed"})
		return
	}
	if !ok {
		c.JSON(http.StatusBadRequest, api.ErrorBody{Error: "invalid or expired code"})
		return
	}
	if err := s.store.SetEmailVerified(c.Request.Context(), req.Email, true); err != nil {
		c.JSON(http.StatusInternalServerError, api.ErrorBody{Error: "could not update user"})
		return
	}
	u, err := s.store.GetUserByEmail(c.Request.Context(), req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, api.ErrorBody{Error: "user lookup failed"})
		return
	}
	resp, err := s.authResponse(c.Request.Context(), u)
	if err != nil {
		c.JSON(http.StatusInternalServerError, api.ErrorBody{Error: "could not issue token"})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (s *Server) ResendVerification(c *gin.Context) {
	var req api.ResendVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, api.ErrorBody{Error: "invalid request body"})
		return
	}
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	u, err := s.store.GetUserByEmail(c.Request.Context(), req.Email)
	if err != nil {
		if err == store.ErrNotFound {
			c.JSON(http.StatusNotFound, api.ErrorBody{Error: "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, api.ErrorBody{Error: "lookup failed"})
		return
	}
	if u.EmailVerified {
		c.JSON(http.StatusBadRequest, api.ErrorBody{Error: "email already verified"})
		return
	}
	code := sixDigitCode()
	expires := time.Now().Add(s.cfg.VerificationCodeTTL)
	if err := s.store.SaveEmailCode(c.Request.Context(), u.Email, code, expires); err != nil {
		c.JSON(http.StatusInternalServerError, api.ErrorBody{Error: "could not save code"})
		return
	}
	s.mail.SendVerificationCode(u.Email, code)
	c.Status(http.StatusNoContent)
}

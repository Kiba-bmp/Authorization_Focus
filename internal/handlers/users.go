package handlers

import (
	"net/http"
	"strings"
	"time"

	"focus/account-cabinet/internal/api"
	"focus/account-cabinet/internal/auth"
	"focus/account-cabinet/internal/store"

	"github.com/gin-gonic/gin"
)

func (s *Server) GetMe(c *gin.Context) {
	uid, ok := userIDFromCtx(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, api.ErrorBody{Error: "unauthorized"})
		return
	}
	u, err := s.store.GetUserByID(c.Request.Context(), uid)
	if err != nil {
		if err == store.ErrNotFound {
			c.JSON(http.StatusUnauthorized, api.ErrorBody{Error: "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, api.ErrorBody{Error: "lookup failed"})
		return
	}
	c.JSON(http.StatusOK, api.UserMeResponse{
		UserID:       u.ID,
		UserName:     u.UserName,
		Email:        u.Email,
		AvatarEmoji:  u.AvatarEmoji,
		InviteCode:   u.InviteCode,
		TodayMinutes: u.TodayMinutes,
		WeekMinutes:  u.WeekMinutes,
	})
}

func (s *Server) UpdateProfile(c *gin.Context) {
	uid, ok := userIDFromCtx(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, api.ErrorBody{Error: "unauthorized"})
		return
	}
	var req api.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, api.ErrorBody{Error: "invalid request body"})
		return
	}
	u, err := s.store.GetUserByID(c.Request.Context(), uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, api.ErrorBody{Error: "lookup failed"})
		return
	}
	if !auth.CheckPassword(u.PasswordHash, req.CurrentPassword) {
		c.JSON(http.StatusUnauthorized, api.ErrorBody{Error: "invalid current password"})
		return
	}
	var newHash *string
	if req.NewPassword != nil && strings.TrimSpace(*req.NewPassword) != "" {
		if len(strings.TrimSpace(*req.NewPassword)) < 8 {
			c.JSON(http.StatusBadRequest, api.ErrorBody{Error: "new password must be at least 8 characters"})
			return
		}
		h, err := auth.HashPassword(strings.TrimSpace(*req.NewPassword))
		if err != nil {
			c.JSON(http.StatusInternalServerError, api.ErrorBody{Error: "could not hash password"})
			return
		}
		newHash = &h
	}
	oldEmail := u.Email
	var emailPtr *string
	if req.Email != nil {
		e := strings.TrimSpace(strings.ToLower(*req.Email))
		if e == "" {
			c.JSON(http.StatusBadRequest, api.ErrorBody{Error: "email cannot be empty"})
			return
		}
		emailPtr = &e
	}
	err = s.store.UpdateProfile(c.Request.Context(), uid, req.UserName, emailPtr, newHash, req.AvatarEmoji)
	if err != nil {
		if err == store.ErrEmailTaken {
			c.JSON(http.StatusConflict, api.ErrorBody{Error: "email already in use"})
			return
		}
		c.JSON(http.StatusInternalServerError, api.ErrorBody{Error: "update failed"})
		return
	}
	if s.cfg.SkipEmailVerification && emailPtr != nil && *emailPtr != oldEmail {
		nu, _ := s.store.GetUserByID(c.Request.Context(), uid)
		_ = s.store.SetEmailVerified(c.Request.Context(), nu.Email, true)
	} else if emailPtr != nil && *emailPtr != oldEmail {
		nu, err := s.store.GetUserByID(c.Request.Context(), uid)
		if err == nil && !nu.EmailVerified {
			code := sixDigitCode()
			expires := time.Now().Add(s.cfg.VerificationCodeTTL)
			if err := s.store.SaveEmailCode(c.Request.Context(), nu.Email, code, expires); err == nil {
				s.mail.SendVerificationCode(nu.Email, code)
			}
		}
	}
	c.Status(http.StatusNoContent)
}

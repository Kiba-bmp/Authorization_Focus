package handlers

import (
	"errors"
	"net/http"
	"time"

	"focus/account-cabinet/internal/api"
	"focus/account-cabinet/internal/repository"
	"focus/account-cabinet/internal/service"

	"github.com/gin-gonic/gin"
)

func toAuthResponse(result *service.AuthResult) api.AuthResponse {
	return api.AuthResponse{
		UserID:      result.UserID,
		Email:       result.Email,
		AccessToken: result.AccessToken,
		ExpiresAt:   result.ExpiresAt.UTC().Format(time.RFC3339Nano),
	}
}

func (s *Server) Register(c *gin.Context) {
	var req api.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, api.ErrorBody{Error: "invalid request body"})
		return
	}

	result, err := s.auth.Register(c.Request.Context(), req.UserName, req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrEmailTaken):
			c.JSON(http.StatusConflict, api.ErrorBody{Error: "email already registered"})
		default:
			c.JSON(http.StatusInternalServerError, api.ErrorBody{Error: "registration failed"})
		}
		return
	}

	c.JSON(http.StatusOK, toAuthResponse(result))
}

func (s *Server) Login(c *gin.Context) {
	var req api.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, api.ErrorBody{Error: "invalid request body"})
		return
	}

	result, err := s.auth.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidCredentials):
			c.JSON(http.StatusUnauthorized, api.ErrorBody{Error: "invalid email or password"})
		default:
			c.JSON(http.StatusInternalServerError, api.ErrorBody{Error: "login failed"})
		}
		return
	}

	c.JSON(http.StatusOK, toAuthResponse(result))
}

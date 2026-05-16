package handlers

import (
	"net/http"
	"time"

	"focus/account-cabinet/internal/api"
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
		c.JSON(http.StatusBadRequest, api.ErrorBody{Error: "Некорректное тело запроса."})
		return
	}

	result, err := s.auth.Register(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		if appErr, ok := service.AsAppError(err); ok {
			c.JSON(appErr.StatusCode, api.ErrorBody{Error: appErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, api.ErrorBody{Error: "Не удалось выполнить регистрацию."})
		return
	}

	c.JSON(http.StatusOK, toAuthResponse(result))
}

func (s *Server) Login(c *gin.Context) {
	var req api.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, api.ErrorBody{Error: "Некорректное тело запроса."})
		return
	}

	result, err := s.auth.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		if appErr, ok := service.AsAppError(err); ok {
			c.JSON(appErr.StatusCode, api.ErrorBody{Error: appErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, api.ErrorBody{Error: "Не удалось выполнить вход."})
		return
	}

	c.JSON(http.StatusOK, toAuthResponse(result))
}

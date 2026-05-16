package handlers

import (
	"net/http"
	"strings"

	"focus/account-cabinet/internal/api"
	"focus/account-cabinet/internal/auth"

	"github.com/gin-gonic/gin"
)

const ctxUserID = "userID"

func (s *Server) RequireAuth(c *gin.Context) {
	h := c.GetHeader("Authorization")
	if h == "" || !strings.HasPrefix(strings.ToLower(h), "bearer ") {
		c.JSON(http.StatusUnauthorized, api.ErrorBody{Error: "Отсутствует или некорректен заголовок авторизации."})
		c.Abort()
		return
	}
	raw := strings.TrimSpace(h[7:])
	uid, err := auth.ParseToken(s.cfg.JWTSecret, raw)
	if err != nil {
		c.JSON(http.StatusUnauthorized, api.ErrorBody{Error: "Токен недействителен или срок его действия истёк."})
		c.Abort()
		return
	}
	c.Set(ctxUserID, uid)
	c.Next()
}

func userIDFromCtx(c *gin.Context) (string, bool) {
	v, ok := c.Get(ctxUserID)
	if !ok {
		return "", false
	}
	id, ok := v.(string)
	return id, ok && id != ""
}

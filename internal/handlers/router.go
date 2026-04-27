package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Mount registers HTTP routes on r.
func (s *Server) Mount(r *gin.Engine) {
	r.GET("/health", health)

	auth := r.Group("/api/auth")
	{
		auth.POST("/register", s.Register)
		auth.POST("/login", s.Login)
	}

	// Спецификация не под /swagger/* — иначе catch-all конфликтует с фиксированным путём в дереве Gin.
	r.GET("/openapi.json", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json", swaggerSpecJSON)
	})
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/openapi.json")))
}

func health(c *gin.Context) {
	c.String(http.StatusOK, "ok")
}

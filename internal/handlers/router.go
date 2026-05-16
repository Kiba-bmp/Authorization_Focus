package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Mount registers HTTP routes on r.
func (s *Server) Mount(r *gin.Engine) {
	r.GET("/api/auth/health", health)

	auth := r.Group("/api/auth")
	{
		auth.POST("/register", s.Register)
		auth.POST("/login", s.Login)
	}

	// Документация auth под тем же префиксом /api/auth
	r.GET("/api/auth/openapi.json", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json", swaggerSpecJSON)
	})
	r.GET("/api/auth/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/api/auth/openapi.json")))
}

func health(c *gin.Context) {
	c.String(http.StatusOK, "ok")
}

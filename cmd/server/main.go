package main

import (
	"log"

	"focus/account-cabinet/internal/appclient"
	"focus/account-cabinet/internal/config"
	"focus/account-cabinet/internal/handlers"
	"focus/account-cabinet/internal/repository"
	"focus/account-cabinet/internal/service"
	"focus/account-cabinet/internal/store"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	st, err := store.Open(cfg)
	if err != nil {
		log.Fatalf("Ошибка запуска хранилища: %v", err)
	}
	defer func() {
		if err := st.Close(); err != nil {
			log.Printf("Ошибка при закрытии соединения с базой данных: %v", err)
		}
	}()

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	db := st.DB()
	userRepo := repository.NewUserRepository(db)
	backendClient := appclient.New(cfg.BackendURL, cfg.BackendInternalKey)

	authService := service.NewAuthService(cfg, userRepo, backendClient)

	srv := handlers.NewServer(authService, cfg)
	srv.Mount(r)

	log.Printf("Сервис auth запущен на %s. Swagger: http://127.0.0.1%s/swagger/index.html", cfg.Addr, cfg.Addr)
	if err := r.Run(cfg.Addr); err != nil {
		log.Fatalf("Ошибка запуска HTTP-сервера: %v", err)
	}
}

package main

import (
	"log"
	"os"

	"focus/account-cabinet/internal/config"
	"focus/account-cabinet/internal/handlers"
	"focus/account-cabinet/internal/mailer"
	"focus/account-cabinet/internal/store"

	"github.com/gin-gonic/gin"
)

// Сервис личного кабинета: JWT, регистрация/логин, подтверждение почты, профиль. Платежи не реализуются.
func main() {
	cfg := config.Load()
	st, err := store.Open(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	defer st.Close()

	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	m := &mailer.Console{}
	srv := handlers.NewServer(st, cfg, m)
	srv.Mount(r)

	log.Printf("account-cabinet listening on %s  swagger: http://127.0.0.1%s/swagger/index.html\n", cfg.Addr, cfg.Addr)
	if err := r.Run(cfg.Addr); err != nil {
		log.Fatal(err)
	}
}

package handlers

import (
	"focus/account-cabinet/internal/config"
	"focus/account-cabinet/internal/service"
)

type Server struct {
	auth *service.AuthService
	cfg  config.Config
}

func NewServer(authSvc *service.AuthService, cfg config.Config) *Server {
	return &Server{auth: authSvc, cfg: cfg}
}

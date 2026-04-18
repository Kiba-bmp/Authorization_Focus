package handlers

import (
	"focus/account-cabinet/internal/config"
	"focus/account-cabinet/internal/mailer"
	"focus/account-cabinet/internal/store"
)

type Server struct {
	store *store.Store
	cfg   config.Config
	mail  *mailer.Console
}

func NewServer(st *store.Store, cfg config.Config, mail *mailer.Console) *Server {
	return &Server{store: st, cfg: cfg, mail: mail}
}

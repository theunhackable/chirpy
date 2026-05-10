package handler

import "github.com/theunhackable/chirpy/internal/config"

type Handler struct {
	Cfg *config.ApiConfig
}

func NewHandler(cfg *config.ApiConfig) *Handler {
	return &Handler{Cfg: cfg}
}

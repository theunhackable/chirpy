package handler

import (
	"context"
	"net/http"

	"github.com/theunhackable/chirpy/internal/config"
	helper "github.com/theunhackable/chirpy/internal/helpers"
	model "github.com/theunhackable/chirpy/internal/models"
)

func Reset(w http.ResponseWriter, r *http.Request) {
	cfg := config.GetConf()

	if cfg.Platform != "dev" {
		helper.RespondWithError(w, 403, "you dont have access")
	}

	if err := cfg.Q.DeleteChirps(context.Background()); err != nil {
		helper.RespondWithError(w, 500, "Error deleting chirps")
		return
	}

	if err := cfg.Q.DeleteUsers(context.Background()); err != nil {
		helper.RespondWithError(w, 500, "Error deleting users")
		return
	}

	if err := cfg.Q.DeleteRefreshTokens(context.Background()); err != nil {
		helper.RespondWithError(w, 500, "Error deleting refresh tokens")
		return
	}

	msg := model.Message{
		Message: "Cleared tables successfully.",
		Type:    "message",
	}
	helper.RespondWithJson(w, 200, msg)
}

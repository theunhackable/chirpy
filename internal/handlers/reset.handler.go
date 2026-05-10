package handler

import (
	"net/http"

	helper "github.com/theunhackable/chirpy/internal/helpers"
	model "github.com/theunhackable/chirpy/internal/models"
)

func (h *Handler) Reset(w http.ResponseWriter, r *http.Request) {
	if h.Cfg.Platform != "dev" {
		helper.RespondWithError(w, 403, "you dont have access")
		return
	}

	if err := h.Cfg.Q.DeleteChirps(r.Context()); err != nil {
		helper.RespondWithError(w, 500, "Error deleting chirps")
		return
	}

	if err := h.Cfg.Q.DeleteUsers(r.Context()); err != nil {
		helper.RespondWithError(w, 500, "Error deleting users")
		return
	}

	if err := h.Cfg.Q.DeleteRefreshTokens(r.Context()); err != nil {
		helper.RespondWithError(w, 500, "Error deleting refresh tokens")
		return
	}

	helper.RespondWithJson(w, 200, model.Message{
		Message: "Cleared tables successfully.",
		Type:    "message",
	})
}

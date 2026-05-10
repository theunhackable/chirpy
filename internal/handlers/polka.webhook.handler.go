package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/theunhackable/chirpy/internal/auth"
	"github.com/theunhackable/chirpy/internal/database"
	helper "github.com/theunhackable/chirpy/internal/helpers"
	model "github.com/theunhackable/chirpy/internal/models"
)

func (h *Handler) PolkaWebhook(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetAPIKey(r.Header)

	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	if apiKey != h.Cfg.PolkaKey {
		helper.RespondWithError(w, http.StatusUnauthorized, "invalid API key")
		return
	}

	params := model.PolkaReq{}

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&params); err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	userId, err := uuid.Parse(params.Data.UserId)

	if err != nil {
		helper.RespondWithJson(w, http.StatusNoContent, nil)
		return
	}

	if params.Event != "user.upgraded" {
		helper.RespondWithJson(w, http.StatusNoContent, nil)
		return
	}

	if err := h.Cfg.Q.UpdateUserSubscriptionById(r.Context(),
		database.UpdateUserSubscriptionByIdParams{
			IsChirpyRed: true,
			ID:          userId,
			UpdatedAt:   time.Now(),
		}); err != nil {
		helper.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	helper.RespondWithJson(w, http.StatusNoContent, nil)
}

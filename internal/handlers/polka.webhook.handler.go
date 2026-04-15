package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/theunhackable/chirpy/internal/auth"
	"github.com/theunhackable/chirpy/internal/config"
	"github.com/theunhackable/chirpy/internal/database"
	helper "github.com/theunhackable/chirpy/internal/helpers"
	model "github.com/theunhackable/chirpy/internal/models"
)

func PoolkaWebhook(w http.ResponseWriter, r *http.Request) {

	apiKey, err := auth.GetAPIKey(r.Header)

	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	if apiKey != os.Getenv("POLKA_KEY") {
		helper.RespondWithError(w, http.StatusUnauthorized, err.Error())
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

	cfg := config.GetConf()

	if err := cfg.Q.UpdateUserSubscriptionById(context.Background(),
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

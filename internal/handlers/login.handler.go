package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/theunhackable/chirpy/internal/auth"
	"github.com/theunhackable/chirpy/internal/config"
	"github.com/theunhackable/chirpy/internal/database"
	helper "github.com/theunhackable/chirpy/internal/helpers"
	model "github.com/theunhackable/chirpy/internal/models"
)

func Login(w http.ResponseWriter, r *http.Request) {

	var params = model.UserReq{}

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&params); err != nil {
		helper.RespondWithError(w, http.StatusUnprocessableEntity,
			fmt.Sprintf("Something went wrong loging in user: %v", err))
		return
	}
	cfg := config.GetConf()
	user, err := cfg.Q.GetUserByEmail(context.Background(), params.Email)

	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized,
			fmt.Sprintf("Unable to get user with email %s", params.Email))
		return
	}

	isSame, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)

	if err != nil || isSame == false {
		helper.RespondWithError(w, http.StatusUnauthorized,
			fmt.Sprintf("Unable to get user with email %s", params.Email))
		return
	}

	token, err := auth.MakeJWT(user.ID,
		cfg.JWTSecret,
		time.Hour,
	)

	if err != nil || isSame == false {
		helper.RespondWithError(w, http.StatusInternalServerError,
			err.Error(),
		)
		return
	}
	refreshToken := auth.MakeRefreshToken()

	if err := cfg.Q.CreateRefreshToken(context.Background(),
		database.CreateRefreshTokenParams{
			Token:     refreshToken,
			UserID:    user.ID,
			ExpiresAt: time.Now().AddDate(0, 0, 60).UTC(),
		}); err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError,
			err.Error(),
		)

	}

	helper.RespondWithJson(w, http.StatusOK, model.User{
		ID:           user.ID,
		Email:        user.Email,
		IsChirpyRed:  user.IsChirpyRed,
		Token:        token,
		RefreshToken: refreshToken,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	})

}

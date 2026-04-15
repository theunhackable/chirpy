package handler

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/theunhackable/chirpy/internal/auth"
	"github.com/theunhackable/chirpy/internal/config"
	"github.com/theunhackable/chirpy/internal/database"
	helper "github.com/theunhackable/chirpy/internal/helpers"
	model "github.com/theunhackable/chirpy/internal/models"
)

func RefreshAccessToken(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)

	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	cfg := config.GetConf()

	tokenInfo, err := cfg.Q.GetRefreshTokenByToken(context.Background(), token)

	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	if tokenInfo.ExpiresAt.Before(time.Now()) || tokenInfo.RevokedAt.Valid {
		helper.RespondWithError(w, http.StatusUnauthorized, "Token is expired or revoked")
		return
	}

	newToken, err := auth.MakeJWT(tokenInfo.UserID,
		cfg.JWTSecret,
		time.Hour,
	)

	helper.RespondWithJson(w, http.StatusOK, model.AccessToken{
		Token: newToken,
	})

}

func RevokeRefreshToken(w http.ResponseWriter, r *http.Request) {

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	cfg := config.GetConf()

	editTime := time.Now().UTC()
	revokedTime := sql.NullTime{
		Time:  editTime,
		Valid: true,
	}

	if err := cfg.Q.UpdateRefreshTokenByToken(context.Background(),
		database.UpdateRefreshTokenByTokenParams{
			Token:     token,
			UpdatedAt: editTime,
			RevokedAt: revokedTime}); err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.RespondWithJson(w, http.StatusNoContent, nil)
}

package handler

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/theunhackable/chirpy/internal/auth"
	"github.com/theunhackable/chirpy/internal/database"
	helper "github.com/theunhackable/chirpy/internal/helpers"
	model "github.com/theunhackable/chirpy/internal/models"
)

func (h *Handler) RefreshAccessToken(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)

	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	tokenInfo, err := h.Cfg.Q.GetRefreshTokenByToken(r.Context(), token)

	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	if tokenInfo.ExpiresAt.Before(time.Now()) || tokenInfo.RevokedAt.Valid {
		helper.RespondWithError(w, http.StatusUnauthorized, "Token is expired or revoked")
		return
	}

	newToken, err := auth.MakeJWT(tokenInfo.UserID,
		h.Cfg.JWTSecret,
		time.Hour,
	)

	helper.RespondWithJson(w, http.StatusOK, model.AccessToken{
		Token: newToken,
	})
}

func (h *Handler) RevokeRefreshToken(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	editTime := time.Now().UTC()
	revokedTime := sql.NullTime{
		Time:  editTime,
		Valid: true,
	}

	if err := h.Cfg.Q.UpdateRefreshTokenByToken(r.Context(),
		database.UpdateRefreshTokenByTokenParams{
			Token:     token,
			UpdatedAt: editTime,
			RevokedAt: revokedTime}); err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.RespondWithJson(w, http.StatusNoContent, nil)
}

package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/theunhackable/chirpy/internal/auth"
	"github.com/theunhackable/chirpy/internal/database"
	helper "github.com/theunhackable/chirpy/internal/helpers"
	model "github.com/theunhackable/chirpy/internal/models"
)

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	params := model.UserReq{}

	if err := decoder.Decode(&params); err != nil {
		helper.RespondWithError(w, http.StatusUnprocessableEntity,
			fmt.Sprintf("Something went wrong decoding user create params: %v", err))
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)

	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError,
			fmt.Sprintf("Something went wrong while hashing user password: %v", err))
		return
	}

	newUser := database.CreateUserParams{
		ID:             uuid.New(),
		Email:          params.Email,
		HashedPassword: hashedPassword,
	}
	user, err := h.Cfg.Q.CreateUser(r.Context(), newUser)

	if err != nil {
		helper.RespondWithError(w, 500, fmt.Sprintf("Something went wrong: %v", err))
		return
	}

	helper.RespondWithJson(w, 201, model.User{
		ID:          user.ID,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	})
}

func (h *Handler) EditUser(w http.ResponseWriter, r *http.Request) {
	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	userInfo, err := auth.ValidateJWT(accessToken, h.Cfg.JWTSecret)

	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	decoder := json.NewDecoder(r.Body)

	params := model.UserReq{}
	if err := decoder.Decode(&params); err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	updatedUser, err := h.Cfg.Q.UpdateUserById(r.Context(), database.UpdateUserByIdParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
		ID:             userInfo,
		UpdatedAt:      time.Now().UTC(),
	})

	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.RespondWithJson(w, http.StatusOK, model.User{
		ID:          updatedUser.ID,
		Email:       updatedUser.Email,
		IsChirpyRed: updatedUser.IsChirpyRed,
		CreatedAt:   updatedUser.CreatedAt,
		UpdatedAt:   updatedUser.UpdatedAt,
	})
}

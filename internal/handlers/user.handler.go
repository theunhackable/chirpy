package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/theunhackable/chirpy/internal/auth"
	"github.com/theunhackable/chirpy/internal/config"
	"github.com/theunhackable/chirpy/internal/database"
	helper "github.com/theunhackable/chirpy/internal/helpers"
	model "github.com/theunhackable/chirpy/internal/models"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {

	cfg := config.GetConf()
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
	user, err := cfg.Q.CreateUser(r.Context(), newUser)

	if err != nil {
		helper.RespondWithError(w, 500, fmt.Sprintf("Something went wrong: %v", err))
		return
	}

	res := model.User{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	helper.RespondWithJson(w, 201, res)
}

func EditUser(w http.ResponseWriter, r *http.Request) {

	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	cfg := config.GetConf()

	userInfo, err := auth.ValidateJWT(accessToken, cfg.JWTSecret)

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

	updatedUser, err := cfg.Q.UpdateUserById(context.Background(), database.UpdateUserByIdParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
		ID:             userInfo,
		UpdatedAt:      time.Now().UTC(),
	})

	if err != nil {
		fmt.Println("error updating the user...", updatedUser)
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.RespondWithJson(w, http.StatusOK, model.User{
		ID:        updatedUser.ID,
		Email:     updatedUser.Email,
		CreatedAt: updatedUser.CreatedAt,
		UpdatedAt: updatedUser.UpdatedAt,
	})
}

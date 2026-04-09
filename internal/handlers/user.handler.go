package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

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

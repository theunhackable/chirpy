package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
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
		helper.RespondWithError(w, 422, "Something went wrong")
		return
	}
	newUser := database.CreateUserParams{
		ID:    uuid.New(),
		Email: params.Email,
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

package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/theunhackable/chirpy/internal/auth"
	"github.com/theunhackable/chirpy/internal/config"
	"github.com/theunhackable/chirpy/internal/database"
	helper "github.com/theunhackable/chirpy/internal/helpers"
	model "github.com/theunhackable/chirpy/internal/models"
)

func GetChirps(w http.ResponseWriter, r *http.Request) {
	cfg := config.GetConf()

	chirps, err := cfg.Q.GetAllChirps(context.Background())

	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
	}

	respChirps := make([]model.ChirpResp, len(chirps))
	for i, chirp := range chirps {
		respChirps[i] = model.ChirpResp{
			Id:        chirp.ID,
			UserId:    chirp.UserID,
			Body:      chirp.Body,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
		}
	}
	helper.RespondWithJson(w, http.StatusOK, respChirps)
}

func GetChirpById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("chirpID")
	cfg := config.GetConf()
	parsedID, err := uuid.Parse(id)
	if err != nil {
		helper.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	chirp, err := cfg.Q.GetChirpById(context.Background(), parsedID)
	if err != nil {
		helper.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	helper.RespondWithJson(w, http.StatusOK, model.ChirpResp{
		Id:        chirp.ID,
		UserId:    chirp.UserID,
		Body:      chirp.Body,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
	})

}

func PostChirp(w http.ResponseWriter, r *http.Request) {

	cfg := config.GetConf()

	token, err := auth.GetBearerToken(r.Header)

	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	userId, err := auth.ValidateJWT(token, cfg.JWTSecret)

	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := model.ChirpReq{}

	if err := decoder.Decode(&params); err != nil {
		helper.RespondWithError(w, http.StatusUnprocessableEntity, "Request body parsing failed.")
		return
	}

	if len(params.Body) > 140 {
		helper.RespondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	cleanBody := helper.Clean(params.Body)

	chirp, err := cfg.Q.CreateChip(context.Background(), database.CreateChipParams{
		ID:     uuid.New(),
		UserID: userId,
		Body:   cleanBody,
	})

	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
	}

	helper.RespondWithJson(w, http.StatusCreated, model.ChirpResp{
		Id:        chirp.ID,
		UserId:    chirp.UserID,
		Body:      chirp.Body,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
	})
}

func DeleteChirpById(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("chirpID")
	parsedID, err := uuid.Parse(id)
	if err != nil {
		helper.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	cfg := config.GetConf()

	token, err := auth.GetBearerToken(r.Header)

	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	userId, err := auth.ValidateJWT(token, cfg.JWTSecret)

	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// i need to verify if the chirpid belong to the user or not

	chirp, err := cfg.Q.GetChirpById(context.Background(), parsedID)

	if err != nil {
		helper.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	if chirp.UserID != userId {
		helper.RespondWithError(w, http.StatusForbidden, "You are not authorized delete this chirp.")
		return
	}

	if err := cfg.Q.DeleteChirpById(context.Background(), chirp.ID); err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.RespondWithJson(w, http.StatusNoContent, nil)

}

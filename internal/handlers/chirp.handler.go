package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/theunhackable/chirpy/internal/config"
	"github.com/theunhackable/chirpy/internal/database"
	helper "github.com/theunhackable/chirpy/internal/helpers"
	model "github.com/theunhackable/chirpy/internal/models"
)

func GetChirps(w http.ResponseWriter, r *http.Request) {
	cfg := config.GetConf()

	chirps, err := cfg.Q.GetAllChirps(context.Background())

	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Unalbe to get chirps from the db: %v", err))
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
		helper.RespondWithError(w, http.StatusNotFound, fmt.Sprintf("Unable to parse the given chirp id to get the chirp: %v", err))
		return
	}

	chirp, err := cfg.Q.GetChirpById(context.Background(), parsedID)
	if err != nil {
		helper.RespondWithError(w, http.StatusNotFound, fmt.Sprintf("Unable get chirp with given id: %v", err))
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

	cfg := config.GetConf()

	cleanBody := helper.Clean(params.Body)

	chirp, err := cfg.Q.CreateChip(context.Background(), database.CreateChipParams{
		ID:     uuid.New(),
		UserID: params.UserId,
		Body:   cleanBody,
	})

	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Something went wrong while creating chirp: %v", err))
	}

	helper.RespondWithJson(w, http.StatusCreated, model.ChirpResp{
		Id:        chirp.ID,
		UserId:    chirp.UserID,
		Body:      chirp.Body,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
	})
}

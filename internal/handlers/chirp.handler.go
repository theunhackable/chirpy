package handler

import (
	"encoding/json"
	"net/http"
	"sort"

	"github.com/google/uuid"
	"github.com/theunhackable/chirpy/internal/auth"
	"github.com/theunhackable/chirpy/internal/database"
	helper "github.com/theunhackable/chirpy/internal/helpers"
	model "github.com/theunhackable/chirpy/internal/models"
)

func (h *Handler) GetChirps(w http.ResponseWriter, r *http.Request) {
	var filteredChirps []database.Chirp

	authorId := r.URL.Query().Get("author_id")
	sortQ := r.URL.Query().Get("sort")

	if authorId != "" {
		parsedAuthorId, err := uuid.Parse(authorId)
		if err != nil {
			helper.RespondWithError(w, http.StatusNotFound, err.Error())
			return
		}

		chirps, err := h.Cfg.Q.GetChirpsByUserId(r.Context(), parsedAuthorId)
		if err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		filteredChirps = chirps
	} else {
		chirps, err := h.Cfg.Q.GetAllChirps(r.Context())

		if err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		filteredChirps = chirps
	}

	respChirps := make([]model.ChirpResp, len(filteredChirps))
	for i, chirp := range filteredChirps {
		respChirps[i] = model.ChirpResp{
			Id:        chirp.ID,
			UserId:    chirp.UserID,
			Body:      chirp.Body,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
		}
	}

	if sortQ == "desc" {
		sort.Slice(respChirps, func(i, j int) bool {
			return respChirps[i].CreatedAt.After(respChirps[j].CreatedAt)
		})
	}

	helper.RespondWithJson(w, http.StatusOK, respChirps)
}

func (h *Handler) GetChirpById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("chirpID")
	parsedID, err := uuid.Parse(id)
	if err != nil {
		helper.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	chirp, err := h.Cfg.Q.GetChirpById(r.Context(), parsedID)
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

func (h *Handler) PostChirp(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)

	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	userId, err := auth.ValidateJWT(token, h.Cfg.JWTSecret)

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

	chirp, err := h.Cfg.Q.CreateChirp(r.Context(), database.CreateChirpParams{
		ID:     uuid.New(),
		UserID: userId,
		Body:   cleanBody,
	})

	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.RespondWithJson(w, http.StatusCreated, model.ChirpResp{
		Id:        chirp.ID,
		UserId:    chirp.UserID,
		Body:      chirp.Body,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
	})
}

func (h *Handler) DeleteChirpById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("chirpID")
	parsedID, err := uuid.Parse(id)
	if err != nil {
		helper.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	token, err := auth.GetBearerToken(r.Header)

	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	userId, err := auth.ValidateJWT(token, h.Cfg.JWTSecret)

	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	chirp, err := h.Cfg.Q.GetChirpById(r.Context(), parsedID)

	if err != nil {
		helper.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	if chirp.UserID != userId {
		helper.RespondWithError(w, http.StatusForbidden, "You are not authorized delete this chirp.")
		return
	}

	if err := h.Cfg.Q.DeleteChirpById(r.Context(), chirp.ID); err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.RespondWithJson(w, http.StatusNoContent, nil)
}

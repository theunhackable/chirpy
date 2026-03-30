package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func clean(unclean string) string {
	badWords := []string{"kerfuffle", "sharbert", "fornax"}
	fmt.Print(badWords)

	uncleanWords := strings.Split(unclean, " ")

	for _, badWord := range badWords {
		for uncleanInd, uncleanWord := range uncleanWords {
			if strings.ToLower(uncleanWord) == badWord {
				uncleanWords[uncleanInd] = "****"
			}
		}
	}
	return strings.Join(uncleanWords, " ")
}
func validateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, 422, "Something went wrong")
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}
	type successBody struct {
		CleanedBody string `json:"cleaned_body"`
	}

	cleanedText := clean(params.Body)
	respondWithJson(w, 200, successBody{CleanedBody: cleanedText})
}

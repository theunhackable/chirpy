package helper

import (
	"encoding/json"
	"net/http"
)

func RespondWithJson(w http.ResponseWriter, code int, payload any) {
	type retVals struct {
		Valid bool `json:"valid"`
	}
	w.Header().Set("Content-Type", "application/json")

	body, err := json.Marshal(payload)

	if err != nil {
		RespondWithError(w, 500, "Something went wrong")
		return
	}
	w.WriteHeader(code)
	w.Write(body)

}

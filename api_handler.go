package main

import (
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/theunhackable/chirpy/internal/database"
)

type apiConfig struct {
	fileServerHits atomic.Int32
	q              *database.Queries
}

func (cfg *apiConfig) middlewareMetricInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("added 1 to the hits")
		cfg.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

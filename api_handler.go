package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileServerHits atomic.Int32
}

func metricsHandler(cfg *apiConfig) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		hitCount := fmt.Sprintf(`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, cfg.fileServerHits.Load())
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(hitCount))
	}
}

func resetHandler(cfg *apiConfig) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("reset")
		w.WriteHeader(http.StatusOK)
		cfg.fileServerHits.Store(0)
	}
}

func (cfg *apiConfig) middlewareMetricInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("added 1 to the hits")
		cfg.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

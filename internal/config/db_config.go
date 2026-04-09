package config

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/theunhackable/chirpy/internal/database"
)

type ApiConfig struct {
	FileServerHits atomic.Int32
	Q              *database.Queries
	Platform       string
	JWTSecret      string
}

func (cfg *ApiConfig) MiddlewareMetricInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("added 1 to the hits")
		cfg.FileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
func MetricsHandler(cfg *ApiConfig) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		hitCount := fmt.Sprintf(`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, cfg.FileServerHits.Load())
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(hitCount))
	}
}

func ResetHandler(cfg *ApiConfig) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("reset")
		w.WriteHeader(http.StatusOK)
		cfg.FileServerHits.Store(0)
	}
}

func GetConf() *ApiConfig {
	dbURL := os.Getenv("DB_URL")

	jwtSecret := os.Getenv("JWT_SECRET")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		panic(err)
	}

	dbQueries := database.New(db)
	cfg := &ApiConfig{
		Q:         dbQueries,
		Platform:  os.Getenv("PLATFORM"),
		JWTSecret: jwtSecret,
	}

	return cfg

}

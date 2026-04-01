package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/theunhackable/chirpy/internal/database"
)

func okBody(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	godotenv.Load()

	fileServerPath := "."
	dbURL := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		panic(err)
	}

	dbQueries := database.New(db)

	mux := http.NewServeMux()

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	cfg := &apiConfig{
		q: dbQueries,
	}

	mux.Handle("GET /app/", http.StripPrefix("/app", cfg.middlewareMetricInc(http.FileServer(http.Dir(fileServerPath)))))
	mux.HandleFunc("GET /admin/metrics", metricsHandler(cfg))
	mux.HandleFunc("POST /admin/reset", resetHandler(cfg))
	mux.HandleFunc("GET /api/healthz", okBody)
	mux.HandleFunc("POST /api/validate_chirp", validateChirp)

	fmt.Printf("running server on http://localhost%v", server.Addr)
	server.ListenAndServe()
}

package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/theunhackable/chirpy/internal/config"
	handler "github.com/theunhackable/chirpy/internal/handlers"
)

func main() {
	godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	} else {
		port = ":" + port
	}

	fileServerPath := "."
	mux := http.NewServeMux()

	cfg := config.GetConf()
	h := handler.NewHandler(cfg)

	mux.Handle("GET /app/", http.StripPrefix("/app", cfg.MiddlewareMetricInc(http.FileServer(http.Dir(fileServerPath)))))
	mux.HandleFunc("GET /admin/metrics", cfg.MetricsHandler())
	mux.HandleFunc("GET /api/healthz", handler.HealthzHandler)
	mux.HandleFunc("GET /api/chirps", h.GetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", h.GetChirpById)
	mux.HandleFunc("POST /api/login", h.Login)
	mux.HandleFunc("POST /admin/reset", h.Reset)
	mux.HandleFunc("POST /api/users", h.CreateUser)
	mux.HandleFunc("PUT /api/users", h.EditUser)
	mux.HandleFunc("POST /api/chirps", h.PostChirp)
	mux.HandleFunc("POST /api/refresh", h.RefreshAccessToken)
	mux.HandleFunc("POST /api/revoke", h.RevokeRefreshToken)
	mux.HandleFunc("POST /api/polka/webhooks", h.PolkaWebhook)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", h.DeleteChirpById)

	server := http.Server{
		Addr:    port,
		Handler: mux,
	}

	fmt.Printf("running server on http://localhost%v\n", port)
	server.ListenAndServe()
}

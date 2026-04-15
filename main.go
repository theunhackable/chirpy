package main

import (
	"fmt"
	"net/http"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/theunhackable/chirpy/internal/config"
	handler "github.com/theunhackable/chirpy/internal/handlers"
)

func main() {
	godotenv.Load()

	fileServerPath := "."
	mux := http.NewServeMux()

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	cfg := config.GetConf()
	mux.Handle("GET /app/", http.StripPrefix("/app", cfg.MiddlewareMetricInc(http.FileServer(http.Dir(fileServerPath)))))
	mux.HandleFunc("GET /admin/metrics", config.MetricsHandler(cfg))
	mux.HandleFunc("GET /api/healthz", handler.HealthzHandler)
	mux.HandleFunc("GET /api/chirps", handler.GetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", handler.GetChirpById)

	mux.HandleFunc("POST /api/login", handler.Login)
	mux.HandleFunc("POST /admin/reset", handler.Reset)
	mux.HandleFunc("POST /api/users", handler.CreateUser)
	mux.HandleFunc("PUT /api/users", handler.EditUser)
	mux.HandleFunc("POST /api/chirps", handler.PostChirp)
	mux.HandleFunc("POST /api/refresh", handler.RefreshAccessToken)
	mux.HandleFunc("POST /api/revoke", handler.RevokeRefreshToken)

	// webhook
	mux.HandleFunc("POST /api/polka/webhooks", handler.PoolkaWebhook)

	mux.HandleFunc("DELETE /api/chirps/{chirpID}", handler.DeleteChirpById)

	fmt.Printf("running server on http://localhost%v", server.Addr)
	server.ListenAndServe()
}

package main

import (
	"fmt"
	"net/http"
)

func okBody(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {

	fileServerPath := "."

	mux := http.NewServeMux()

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	cfg := &apiConfig{}

	mux.Handle("GET /app/", http.StripPrefix("/app", cfg.middlewareMetricInc(http.FileServer(http.Dir(fileServerPath)))))
	mux.HandleFunc("GET /admin/metrics", metricsHandler(cfg))
	mux.HandleFunc("POST /admin/reset", resetHandler(cfg))
	mux.HandleFunc("GET /api/healthz", okBody)

	fmt.Printf("running server on http://localhost%v", server.Addr)
	server.ListenAndServe()
}

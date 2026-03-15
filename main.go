package main

import (
	"fmt"
	"net/http"
)

func okBody(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	w.Write([]byte("OK"))
}

func main() {

	fileServerPath := "."
	mux := http.NewServeMux()

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	mux.Handle("GET /app/", http.StripPrefix("/app", http.FileServer(http.Dir(fileServerPath))))
	mux.HandleFunc("GET /healthz", okBody)

	fmt.Printf("running server on http://localhost%v", server.Addr)
	server.ListenAndServe()
}

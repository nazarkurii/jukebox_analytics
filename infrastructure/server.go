package infra

import (
	"net/http"
	"time"
)

func NewServer() (*http.Server, *http.ServeMux) {
	mux := http.NewServeMux()
	return &http.Server{
		Addr:         ":8080",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      mux,
	}, mux
}

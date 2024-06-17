package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/ri/223", withTimeout(handleRequest))

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	time.Sleep(10 * time.Second)
	fmt.Println("Starting server on :8080")
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		fmt.Printf("Error starting server: %s\n", err)
	}
}

func withTimeout(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		r = r.WithContext(ctx)
		next(w, r)
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	select {
	case <-r.Context().Done():
		fmt.Fprintln(w, "Request was canceled or timed out")
	case <-time.After(2 * time.Second):
		fmt.Fprintln(w, "Request processed successfully")
	}
}

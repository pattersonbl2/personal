package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"ark31/backend/internal/handler"
	"ark31/backend/internal/middleware"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	_ = godotenv.Load()
	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)
	mux.Handle("/api/contact", middleware.RateLimit(http.HandlerFunc(handler.ContactHandler), 5, time.Hour))
	mux.HandleFunc("/api/resume", handler.ResumeHandler)

	var c *cors.Cors
	if origins := os.Getenv("ALLOWED_ORIGINS"); origins != "" {
		parts := strings.Split(origins, ",")
		for i := range parts {
			parts[i] = strings.TrimSpace(parts[i])
		}
		c = cors.New(cors.Options{
			AllowedOrigins:   parts,
			AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
			AllowedHeaders:   []string{"Content-Type"},
			AllowCredentials: false,
		})
	} else {
		c = cors.Default()
	}
	h := middleware.SecurityHeaders(c.Handler(mux))

	log.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", h))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

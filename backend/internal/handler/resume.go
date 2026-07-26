package handler

import (
	"crypto/subtle"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// ResumeHandler accepts GET or POST /api/resume with resume token + Turnstile response.
// Streams PDF if the resume token is valid, Turnstile siteverify succeeds, and the file exists.
func ResumeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.Method == http.MethodPost {
		r.Body = http.MaxBytesReader(w, r.Body, 64*1024)
		if err := r.ParseForm(); err != nil {
			log.Printf("resume: reject reason=validation detail=bad_request ip=%s", clientIP(r))
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
	}

	if !requireTurnstile(w, r, "resume") {
		return
	}

	token := strings.TrimSpace(r.FormValue("token"))
	if token == "" {
		token = strings.TrimSpace(r.URL.Query().Get("token"))
	}

	expected := os.Getenv("RESUME_TOKEN")
	if expected == "" {
		log.Printf("resume: RESUME_TOKEN not set, refusing to serve")
		http.Error(w, "service unavailable", http.StatusServiceUnavailable)
		return
	}

	// Always run constant-time compare even for empty token to avoid timing differences.
	if subtle.ConstantTimeCompare([]byte(token), []byte(expected)) != 1 {
		log.Printf("resume: reject reason=auth detail=bad_token ip=%s", clientIP(r))
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	path := os.Getenv("RESUME_PDF_PATH")
	if path == "" {
		path = "resume.pdf"
	}
	path, err := filepath.Abs(path)
	if err != nil {
		log.Printf("resume: invalid path: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("resume: file not found: %s", path)
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		log.Printf("resume: open failed: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		log.Printf("resume: stat failed: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	log.Printf("resume: accepted ip=%s", clientIP(r))
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=\"resume.pdf\"")
	http.ServeContent(w, r, "resume.pdf", info.ModTime(), f)
}

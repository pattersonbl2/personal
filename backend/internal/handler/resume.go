package handler

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// ResumeHandler accepts GET /api/resume?token=XXX.
// Streams PDF if token is valid and file exists. Returns 401 if invalid, 404/500 if file missing.
func ResumeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	token := r.URL.Query().Get("token")
	expected := os.Getenv("RESUME_TOKEN")
	if expected == "" {
		log.Printf("resume: RESUME_TOKEN not set, refusing to serve")
		http.Error(w, "service unavailable", http.StatusServiceUnavailable)
		return
	}

	if token == "" || token != expected {
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

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=\"resume.pdf\"")
	http.ServeContent(w, r, "resume.pdf", info.ModTime(), f)
}

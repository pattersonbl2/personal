package handler

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestResumeHandlerUnauthorized(t *testing.T) {
	t.Setenv("RESUME_TOKEN", "resume-secret")
	req := httptest.NewRequest(http.MethodGet, "/api/resume?token=wrong", nil)
	rr := httptest.NewRecorder()
	ResumeHandler(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d body=%s", rr.Code, rr.Body.String())
	}
}

func TestResumeHandlerSuccess(t *testing.T) {
	t.Setenv("RESUME_TOKEN", "resume-secret")

	dir := t.TempDir()
	pdfPath := filepath.Join(dir, "resume.pdf")
	if err := os.WriteFile(pdfPath, []byte("%PDF-1.4 test"), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("RESUME_PDF_PATH", pdfPath)

	req := httptest.NewRequest(http.MethodGet, "/api/resume?token=resume-secret", nil)
	rr := httptest.NewRecorder()
	ResumeHandler(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", rr.Code, rr.Body.String())
	}
	if ct := rr.Header().Get("Content-Type"); ct != "application/pdf" {
		t.Fatalf("content-type = %q", ct)
	}
}

func TestResumeHandlerMethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/resume", nil)
	rr := httptest.NewRecorder()
	ResumeHandler(rr, req)
	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d", rr.Code)
	}
}

package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"ark31/backend/internal/turnstile"
)

func TestResumeHandlerRequiresTurnstile(t *testing.T) {
	t.Setenv("TURNSTILE_SECRET", "test-secret")
	t.Setenv("RESUME_TOKEN", "resume-secret")

	prev := turnstileVerifier
	turnstileVerifier = func(ctx context.Context, token, remoteIP string) (*turnstile.Result, error) {
		return &turnstile.Result{Success: false, ErrorCodes: []string{"invalid-input-response"}}, nil
	}
	t.Cleanup(func() { turnstileVerifier = prev })

	body := strings.NewReader("token=resume-secret&cf-turnstile-response=bad")
	req := httptest.NewRequest(http.MethodPost, "/api/resume", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	ResumeHandler(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("status = %d body=%s", rr.Code, rr.Body.String())
	}
}

func TestResumeHandlerSuccess(t *testing.T) {
	t.Setenv("TURNSTILE_SECRET", "test-secret")
	t.Setenv("RESUME_TOKEN", "resume-secret")

	dir := t.TempDir()
	pdfPath := filepath.Join(dir, "resume.pdf")
	if err := os.WriteFile(pdfPath, []byte("%PDF-1.4 test"), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("RESUME_PDF_PATH", pdfPath)

	prev := turnstileVerifier
	turnstileVerifier = func(ctx context.Context, token, remoteIP string) (*turnstile.Result, error) {
		if token != "good" {
			t.Fatalf("token = %q", token)
		}
		if remoteIP == "" {
			t.Fatal("expected remote IP")
		}
		return &turnstile.Result{Success: true}, nil
	}
	t.Cleanup(func() { turnstileVerifier = prev })

	body := strings.NewReader("token=resume-secret&cf-turnstile-response=good")
	req := httptest.NewRequest(http.MethodPost, "/api/resume", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("CF-Connecting-IP", "203.0.113.9")
	rr := httptest.NewRecorder()
	ResumeHandler(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", rr.Code, rr.Body.String())
	}
	if ct := rr.Header().Get("Content-Type"); ct != "application/pdf" {
		t.Fatalf("content-type = %q", ct)
	}
}

func TestResumeHandlerSecretMissing(t *testing.T) {
	_ = os.Unsetenv("TURNSTILE_SECRET")
	t.Setenv("RESUME_TOKEN", "resume-secret")
	req := httptest.NewRequest(http.MethodPost, "/api/resume", strings.NewReader("token=resume-secret"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	ResumeHandler(rr, req)
	if rr.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d", rr.Code)
	}
}

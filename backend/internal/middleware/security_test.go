package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSecurityHeaders_OmitsCSPForResumePDF(t *testing.T) {
	h := SecurityHeaders(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/pdf")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("%PDF-1.5"))
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/resume?token=x", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if csp := rr.Header().Get("Content-Security-Policy"); csp != "" {
		t.Fatalf("expected no CSP on resume PDF response, got %q", csp)
	}
	if rr.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Fatalf("expected nosniff still set")
	}
}

func TestSecurityHeaders_KeepsCSPForAPIJSON(t *testing.T) {
	h := SecurityHeaders(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))

	req := httptest.NewRequest(http.MethodPost, "/api/contact", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if csp := rr.Header().Get("Content-Security-Policy"); csp == "" {
		t.Fatal("expected CSP on non-PDF API responses")
	}
}

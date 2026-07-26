package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"ark31/backend/internal/turnstile"
)

func TestIsTimingSuspicious(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)

	tests := []struct {
		name   string
		formTS string
		want   bool
		reason string
	}{
		{"missing", "", true, "timing_missing"},
		{"invalid", "abc", true, "timing_invalid"},
		{"too_fast", "1699999999", true, "timing_too_fast"}, // 1s before
		{"ok", "1699999980", false, ""},                     // 20s before
		{"stale", "1699990000", true, "timing_stale_or_future"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, reason := isTimingSuspicious(tt.formTS, now)
			if got != tt.want || reason != tt.reason {
				t.Fatalf("isTimingSuspicious(%q) = (%v, %q), want (%v, %q)", tt.formTS, got, reason, tt.want, tt.reason)
			}
		})
	}
}

func TestIsSpam(t *testing.T) {
	if !isSpam("ok", "zekisuquc419@gmail.com", "hello there friend") {
		t.Fatal("expected blocklist hit")
	}
	if !isSpam("Иван", "a@b.com", "hello there friend") {
		t.Fatal("expected cyrillic name")
	}
	if !isSpam("Bob", "a@b.com", "see https://spam.example") {
		t.Fatal("expected URL heuristic")
	}
	if isSpam("Bob", "bob@example.com", "hello there friend") {
		t.Fatal("legitimate message flagged")
	}
}

func TestWantsJSON(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/api/contact", nil)
	r.Header.Set("Accept", "application/json")
	if !wantsJSON(r) {
		t.Fatal("expected JSON for Accept")
	}
	r2 := httptest.NewRequest(http.MethodPost, "/api/contact", nil)
	r2.Header.Set("X-Requested-With", "XMLHttpRequest")
	if !wantsJSON(r2) {
		t.Fatal("expected JSON for X-Requested-With")
	}
	r3 := httptest.NewRequest(http.MethodPost, "/api/contact", nil)
	if wantsJSON(r3) {
		t.Fatal("expected HTML mode by default")
	}
}

func TestContactHandlerJSONConfigError(t *testing.T) {
	prev := turnstileVerifier
	turnstileVerifier = func(ctx context.Context, token, remoteIP string) (*turnstile.Result, error) {
		return &turnstile.Result{Success: true}, nil
	}
	t.Cleanup(func() { turnstileVerifier = prev })

	t.Setenv("TURNSTILE_SECRET_KEY", "test-secret")
	t.Setenv("RESEND_API_KEY", "")
	t.Setenv("CONTACT_EMAIL", "")

	formTS := time.Now().Add(-5 * time.Second).Unix()
	body := strings.NewReader(
		"name=Brandon&email=brandon@example.com&message=Hello there friend&" +
			"form_ts=" + strconv.FormatInt(formTS, 10) + "&cf-turnstile-response=tok&website=",
	)
	req := httptest.NewRequest(http.MethodPost, "/api/contact", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	rr := httptest.NewRecorder()

	ContactHandler(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, body = %s", rr.Code, rr.Body.String())
	}
	var payload map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &payload); err != nil {
		t.Fatal(err)
	}
	if payload["ok"] != false {
		t.Fatalf("payload = %#v", payload)
	}
}

func TestContactHandlerTurnstileFail(t *testing.T) {
	prev := turnstileVerifier
	turnstileVerifier = func(ctx context.Context, token, remoteIP string) (*turnstile.Result, error) {
		return &turnstile.Result{Success: false, ErrorCodes: []string{"invalid-input-response"}}, nil
	}
	t.Cleanup(func() { turnstileVerifier = prev })
	t.Setenv("TURNSTILE_SECRET_KEY", "test-secret")

	formTS := time.Now().Add(-5 * time.Second).Unix()
	body := strings.NewReader(
		"name=Brandon&email=brandon@example.com&message=Hello there friend&" +
			"form_ts=" + strconv.FormatInt(formTS, 10) + "&cf-turnstile-response=bad",
	)
	req := httptest.NewRequest(http.MethodPost, "/api/contact", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	rr := httptest.NewRecorder()
	ContactHandler(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status = %d body=%s", rr.Code, rr.Body.String())
	}
}

func TestContactHandlerHoneypotSilent(t *testing.T) {
	body := strings.NewReader("name=Bot&email=bot@example.com&message=Hello there friend&website=http://spam")
	req := httptest.NewRequest(http.MethodPost, "/api/contact", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	rr := httptest.NewRecorder()
	ContactHandler(rr, req)
	if rr.Code != http.StatusAccepted {
		t.Fatalf("status = %d", rr.Code)
	}
	var payload map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &payload)
	if payload["ok"] != true {
		t.Fatalf("payload = %#v", payload)
	}
}

func TestContactHandlerSecretMissing(t *testing.T) {
	_ = os.Unsetenv("TURNSTILE_SECRET_KEY")
	formTS := time.Now().Add(-5 * time.Second).Unix()
	body := strings.NewReader(
		"name=Brandon&email=brandon@example.com&message=Hello there friend&form_ts=" + strconv.FormatInt(formTS, 10),
	)
	req := httptest.NewRequest(http.MethodPost, "/api/contact", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	rr := httptest.NewRecorder()
	ContactHandler(rr, req)
	if rr.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d body=%s", rr.Code, rr.Body.String())
	}
}

func TestContactHandlerTooFastSilent(t *testing.T) {
	t.Setenv("TURNSTILE_SECRET_KEY", "test-secret")
	formTS := time.Now().Unix() // too fast
	body := strings.NewReader(
		"name=Brandon&email=brandon@example.com&message=Hello there friend&" +
			"form_ts=" + strconv.FormatInt(formTS, 10) + "&cf-turnstile-response=tok",
	)
	req := httptest.NewRequest(http.MethodPost, "/api/contact", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	rr := httptest.NewRecorder()
	ContactHandler(rr, req)
	if rr.Code != http.StatusAccepted {
		t.Fatalf("status = %d body=%s", rr.Code, rr.Body.String())
	}
}

package handler

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"

	"ark31/backend/internal/turnstile"
)

// turnstileVerifier is overridable in tests.
var turnstileVerifier = func(ctx context.Context, token, remoteIP string) (*turnstile.Result, error) {
	client := &turnstile.Client{Secret: os.Getenv("TURNSTILE_SECRET")}
	return client.Verify(ctx, token, remoteIP)
}

// requireTurnstile verifies cf-turnstile-response via canonical siteverify.
// Returns true when the request may proceed. On failure it writes an HTTP error
// and returns false. Secret is read from TURNSTILE_SECRET (fail-closed if unset).
func requireTurnstile(w http.ResponseWriter, r *http.Request, surface string) bool {
	secret := os.Getenv("TURNSTILE_SECRET")
	if secret == "" {
		log.Printf("%s: reject reason=turnstile detail=secret_missing ip=%s", surface, clientIP(r))
		http.Error(w, "service unavailable", http.StatusServiceUnavailable)
		return false
	}

	token := strings.TrimSpace(r.FormValue("cf-turnstile-response"))
	if token == "" {
		token = strings.TrimSpace(r.URL.Query().Get("cf-turnstile-response"))
	}

	tsResult, err := turnstileVerifier(r.Context(), token, clientIP(r))
	if err != nil {
		log.Printf("%s: reject reason=turnstile detail=verify_error err=%v ip=%s", surface, err, clientIP(r))
		http.Error(w, "verification failed", http.StatusBadGateway)
		return false
	}
	if tsResult == nil || !tsResult.Success {
		codes := []string{"failed"}
		if tsResult != nil && len(tsResult.ErrorCodes) > 0 {
			codes = tsResult.ErrorCodes
		}
		log.Printf("%s: reject reason=turnstile detail=%s ip=%s", surface, strings.Join(codes, ","), clientIP(r))
		http.Error(w, "forbidden", http.StatusForbidden)
		return false
	}
	return true
}

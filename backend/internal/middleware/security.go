package middleware

import (
	"net/http"
	"strings"
)

// SecurityHeaders adds common security headers to responses.
// CSP is omitted for /api/resume so browsers can open/download the PDF without
// a blank spinning tab (default-src 'none' breaks built-in PDF viewers).
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
		if !strings.HasPrefix(r.URL.Path, "/api/resume") {
			w.Header().Set("Content-Security-Policy", "default-src 'none'; style-src 'unsafe-inline'; frame-ancestors 'none'")
		}
		w.Header().Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
		w.Header().Set("Server", "") // suppress Go version disclosure
		next.ServeHTTP(w, r)
	})
}

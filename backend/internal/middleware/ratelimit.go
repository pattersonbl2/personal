package middleware

import (
	"net/http"
	"strings"
	"sync"
	"time"
)

// RateLimit returns a handler that limits requests per IP.
// limit = max requests allowed, window = time window (e.g. time.Hour).
func RateLimit(next http.Handler, limit int, window time.Duration) http.Handler {
	mu := sync.Mutex{}
	ips := make(map[string][]time.Time)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		// Prefer Cloudflare's real IP when behind CF; otherwise X-Forwarded-For (first = client)
		if cf := r.Header.Get("CF-Connecting-IP"); cf != "" {
			ip = strings.TrimSpace(cf)
		} else if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
			if idx := strings.Index(forwarded, ","); idx >= 0 {
				ip = strings.TrimSpace(forwarded[:idx])
			} else {
				ip = strings.TrimSpace(forwarded)
			}
		}

		mu.Lock()
		now := time.Now()
		cutoff := now.Add(-window)
		times := ips[ip]

		// Drop requests older than window
		var recent []time.Time
		for _, t := range times {
			if t.After(cutoff) {
				recent = append(recent, t)
			}
		}

		if len(recent) >= limit {
			mu.Unlock()
			http.Error(w, "too many requests", http.StatusTooManyRequests)
			return
		}

		ips[ip] = append(recent, now)
		mu.Unlock()

		next.ServeHTTP(w, r)
	})
}

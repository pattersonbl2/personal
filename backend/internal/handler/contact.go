package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/mail"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/resend/resend-go/v3"
)

const (
	maxNameLen       = 100
	maxEmailLen      = 254
	maxMessageLen    = 5000
	minMessageLen    = 10
	minFormFillDelay = 2 * time.Second
	maxFormAge       = 2 * time.Hour
)

// ContactHandler accepts POST to /api/contact, validates input, sends email via Resend.
// Form fields: name, email, message, form_ts, cf-turnstile-response.
// Honeypot: website (must be empty).
// Responds with JSON when Accept includes application/json (or X-Requested-With is set);
// otherwise returns a branded HTML page.
func ContactHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	jsonMode := wantsJSON(r)

	r.Body = http.MaxBytesReader(w, r.Body, 64*1024) // 64KB limit
	if err := r.ParseForm(); err != nil {
		log.Printf("contact: reject reason=validation detail=bad_request ip=%s", clientIP(r))
		sendContactError(w, jsonMode, "Bad request.", http.StatusBadRequest)
		return
	}

	// Honeypot: bots often fill hidden fields — silent success so they learn nothing.
	if r.FormValue("website") != "" {
		log.Printf("contact: reject reason=spam detail=honeypot ip=%s", clientIP(r))
		sendContactSuccess(w, jsonMode)
		return
	}

	name := strings.TrimSpace(r.FormValue("name"))
	email := strings.TrimSpace(r.FormValue("email"))
	message := strings.TrimSpace(r.FormValue("message"))
	formTS := strings.TrimSpace(r.FormValue("form_ts"))

	if name == "" || email == "" || message == "" {
		log.Printf("contact: reject reason=validation detail=required_fields ip=%s", clientIP(r))
		sendContactError(w, jsonMode, "Name, email, and message are required.", http.StatusBadRequest)
		return
	}
	if utf8.RuneCountInString(name) > maxNameLen {
		log.Printf("contact: reject reason=validation detail=name_too_long ip=%s", clientIP(r))
		sendContactError(w, jsonMode, "Name is too long.", http.StatusBadRequest)
		return
	}
	if utf8.RuneCountInString(email) > maxEmailLen {
		log.Printf("contact: reject reason=validation detail=email_too_long ip=%s", clientIP(r))
		sendContactError(w, jsonMode, "Email is too long.", http.StatusBadRequest)
		return
	}
	if _, err := mail.ParseAddress(email); err != nil {
		log.Printf("contact: reject reason=validation detail=invalid_email ip=%s", clientIP(r))
		sendContactError(w, jsonMode, "Invalid email address.", http.StatusBadRequest)
		return
	}
	if utf8.RuneCountInString(message) > maxMessageLen {
		log.Printf("contact: reject reason=validation detail=message_too_long ip=%s", clientIP(r))
		sendContactError(w, jsonMode, "Message is too long.", http.StatusBadRequest)
		return
	}
	if utf8.RuneCountInString(message) < minMessageLen {
		log.Printf("contact: reject reason=validation detail=message_too_short ip=%s", clientIP(r))
		sendContactError(w, jsonMode, "Message is too short.", http.StatusBadRequest)
		return
	}

	if tooFast, reason := isTimingSuspicious(formTS, time.Now()); tooFast {
		log.Printf("contact: reject reason=spam detail=%s ip=%s", reason, clientIP(r))
		sendContactSuccess(w, jsonMode)
		return
	}

	// Canonical Turnstile siteverify gate (TURNSTILE_SECRET).
	if os.Getenv("TURNSTILE_SECRET") == "" {
		log.Printf("contact: reject reason=turnstile detail=secret_missing ip=%s", clientIP(r))
		sendContactError(w, jsonMode, "Email service is temporarily unavailable.", http.StatusServiceUnavailable)
		return
	}
	token := strings.TrimSpace(r.FormValue("cf-turnstile-response"))
	tsResult, err := turnstileVerifier(r.Context(), token, clientIP(r))
	if err != nil {
		log.Printf("contact: reject reason=turnstile detail=verify_error err=%v ip=%s", err, clientIP(r))
		sendContactError(w, jsonMode, "Verification failed. Please try again.", http.StatusBadGateway)
		return
	}
	if tsResult == nil || !tsResult.Success {
		codes := []string{"failed"}
		if tsResult != nil && len(tsResult.ErrorCodes) > 0 {
			codes = tsResult.ErrorCodes
		}
		log.Printf("contact: reject reason=turnstile detail=%s ip=%s", strings.Join(codes, ","), clientIP(r))
		sendContactError(w, jsonMode, "Verification failed. Please try again.", http.StatusForbidden)
		return
	}

	// Anti-spam: silent-discard so bots don't know they were blocked.
	if isSpam(name, email, message) {
		log.Printf("contact: reject reason=spam detail=heuristics ip=%s", clientIP(r))
		sendContactSuccess(w, jsonMode)
		return
	}

	apiKey := os.Getenv("RESEND_API_KEY")
	contactEmail := os.Getenv("CONTACT_EMAIL")
	if apiKey == "" || contactEmail == "" {
		log.Printf("contact: reject reason=config detail=email_unset ip=%s", clientIP(r))
		sendContactError(w, jsonMode, "Email service is temporarily unavailable.", http.StatusInternalServerError)
		return
	}

	client := resend.NewClient(apiKey)

	subject := fmt.Sprintf("[Contact] from %s", name)
	html := fmt.Sprintf(
		"<p><strong>From:</strong> %s &lt;%s&gt;</p><p><strong>Message:</strong></p><pre>%s</pre>",
		escapeHTML(name), escapeHTML(email), escapeHTML(message),
	)

	params := &resend.SendEmailRequest{
		From:    os.Getenv("RESEND_FROM"),
		To:      []string{contactEmail},
		Subject: subject,
		Html:    html,
		ReplyTo: email,
	}

	if params.From == "" {
		params.From = "Contact Form <onboarding@resend.dev>"
	}

	_, err = client.Emails.Send(params)
	if err != nil {
		log.Printf("contact: reject reason=email detail=send_failed err=%v ip=%s", err, clientIP(r))
		sendContactError(w, jsonMode, "Failed to send message. Please try again later.", http.StatusInternalServerError)
		return
	}

	log.Printf("contact: accepted ip=%s", clientIP(r))
	sendContactSuccess(w, jsonMode)
}

func wantsJSON(r *http.Request) bool {
	if strings.EqualFold(r.Header.Get("X-Requested-With"), "XMLHttpRequest") {
		return true
	}
	accept := r.Header.Get("Accept")
	return strings.Contains(accept, "application/json")
}

func clientIP(r *http.Request) string {
	if cf := r.Header.Get("CF-Connecting-IP"); cf != "" {
		return strings.TrimSpace(cf)
	}
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		if idx := strings.Index(forwarded, ","); idx >= 0 {
			return strings.TrimSpace(forwarded[:idx])
		}
		return strings.TrimSpace(forwarded)
	}
	return r.RemoteAddr
}

// isTimingSuspicious reports whether form_ts indicates a bot-like fill time.
// Missing/invalid timestamps are treated as spam (Turnstile already requires JS).
func isTimingSuspicious(formTS string, now time.Time) (bool, string) {
	if formTS == "" {
		return true, "timing_missing"
	}
	sec, err := strconv.ParseInt(formTS, 10, 64)
	if err != nil {
		return true, "timing_invalid"
	}
	started := time.Unix(sec, 0)
	elapsed := now.Sub(started)
	if elapsed < minFormFillDelay {
		return true, "timing_too_fast"
	}
	if elapsed > maxFormAge || started.After(now.Add(30*time.Second)) {
		return true, "timing_stale_or_future"
	}
	return false, ""
}

// isSpam returns true if the submission looks like spam.
// Rules are intentionally simple and conservative to avoid false positives.
func isSpam(name, email, message string) bool {
	lowerEmail := strings.ToLower(email)
	if lowerEmail == "zekisuquc419@gmail.com" {
		return true
	}
	if containsScript(name, 0x0400, 0x052F) || containsScript(message, 0x0400, 0x052F) {
		return true
	}
	msg := strings.ToLower(message)
	if strings.Contains(msg, "http://") || strings.Contains(msg, "https://") || strings.Contains(msg, "www.") {
		return true
	}
	return false
}

func containsScript(s string, lo, hi rune) bool {
	for _, r := range s {
		if r >= lo && r <= hi {
			return true
		}
	}
	return false
}

func sendContactSuccess(w http.ResponseWriter, jsonMode bool) {
	if jsonMode {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusAccepted)
		_ = json.NewEncoder(w).Encode(map[string]any{"ok": true, "message": "Message received. Thanks!"})
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusAccepted)
	sendStyledHTML(w, "Message received. Thanks!")
}

func sendContactError(w http.ResponseWriter, jsonMode bool, msg string, code int) {
	if jsonMode {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(code)
		_ = json.NewEncoder(w).Encode(map[string]any{"ok": false, "error": msg})
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(code)
	sendStyledHTML(w, escapeHTML(msg))
}

func sendStyledHTML(w http.ResponseWriter, message string) {
	const tpl = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <meta name="color-scheme" content="dark light">
  <title>Contact · ark31.info</title>
  <style>
    :root {
      --bg: #0d1117;
      --bg-card: #161b22;
      --border: #2a2f37;
      --text: #e6edf3;
      --text-muted: #9aa4b2;
      --accent: #5aa2ff;
      --on-accent: #08152e;
      --radius: 6px;
      --font-body: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
      --font-mono: 'Cascadia Code', 'Fira Code', 'JetBrains Mono', ui-monospace, SFMono-Regular, monospace;
    }
    @media (prefers-color-scheme: light) {
      :root {
        --bg: #ffffff;
        --bg-card: #f6f8fa;
        --border: #d0d7de;
        --text: #1f2328;
        --text-muted: #656d76;
        --accent: #2563eb;
        --on-accent: #ffffff;
      }
    }
    * { box-sizing: border-box; }
    body {
      margin: 0; min-height: 100vh; display: flex; align-items: center; justify-content: center;
      font-family: var(--font-body); background: var(--bg); color: var(--text);
      padding: 1.5rem;
    }
    .card {
      background: var(--bg-card); border: 1px solid var(--border); border-radius: var(--radius);
      padding: 2rem 2.25rem; text-align: center; max-width: 420px; width: 100%;
    }
    .mark {
      font-family: var(--font-mono); font-size: 0.8rem; letter-spacing: 0.04em;
      color: var(--accent); margin-bottom: 1rem;
    }
    .card p { margin: 0 0 1.5rem; font-size: 1.1rem; line-height: 1.5; }
    .card a {
      display: inline-block; padding: 0.6rem 1.25rem; background: var(--accent); color: var(--on-accent);
      text-decoration: none; border-radius: var(--radius); font-weight: 500;
    }
    .card a:hover { opacity: .92; }
  </style>
</head>
<body>
  <div class="card">
    <div class="mark">// ark31.info</div>
    <p>{{MESSAGE}}</p>
    <a href="https://ark31.info/contact/">Back to Contact</a>
  </div>
</body>
</html>`
	_, _ = w.Write([]byte(strings.Replace(tpl, "{{MESSAGE}}", message, 1)))
}

func escapeHTML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, `"`, "&quot;")
	s = strings.ReplaceAll(s, "'", "&#39;")
	return s
}

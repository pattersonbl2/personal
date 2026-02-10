package handler

import (
	"fmt"
	"log"
	"net/http"
	"net/mail"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/resend/resend-go/v3"
)

const (
	maxNameLen    = 100
	maxEmailLen   = 254
	maxMessageLen = 5000
)

// ContactHandler accepts POST to /api/contact, validates input, sends email via Resend.
// Form fields: name, email, message. Honeypot: website (must be empty).
func ContactHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 64*1024) // 64KB limit
	if err := r.ParseForm(); err != nil {
		sendErrorHTML(w, "Bad request.", http.StatusBadRequest)
		return
	}

	// Honeypot: bots often fill hidden fields
	if r.FormValue("website") != "" {
		sendErrorHTML(w, "Invalid submission.", http.StatusBadRequest)
		return
	}

	name := strings.TrimSpace(r.FormValue("name"))
	email := strings.TrimSpace(r.FormValue("email"))
	message := strings.TrimSpace(r.FormValue("message"))

	if name == "" || email == "" || message == "" {
		sendErrorHTML(w, "Name, email, and message are required.", http.StatusBadRequest)
		return
	}
	if utf8.RuneCountInString(name) > maxNameLen {
		sendErrorHTML(w, "Name is too long.", http.StatusBadRequest)
		return
	}
	if utf8.RuneCountInString(email) > maxEmailLen {
		sendErrorHTML(w, "Email is too long.", http.StatusBadRequest)
		return
	}
	if _, err := mail.ParseAddress(email); err != nil {
		sendErrorHTML(w, "Invalid email address.", http.StatusBadRequest)
		return
	}
	if utf8.RuneCountInString(message) > maxMessageLen {
		sendErrorHTML(w, "Message is too long.", http.StatusBadRequest)
		return
	}

	apiKey := os.Getenv("RESEND_API_KEY")
	contactEmail := os.Getenv("CONTACT_EMAIL")
	if apiKey == "" || contactEmail == "" {
		log.Println("contact: RESEND_API_KEY or CONTACT_EMAIL not set")
		sendErrorHTML(w, "Email service is temporarily unavailable.", http.StatusInternalServerError)
		return
	}

	client := resend.NewClient(apiKey)

	subject := fmt.Sprintf("[Contact] from %s", name)
	html := fmt.Sprintf(
		"<p><strong>From:</strong> %s &lt;%s&gt;</p><p><strong>Message:</strong></p><pre>%s</pre>",
		escapeHTML(name), escapeHTML(email), escapeHTML(message),
	)

	params := &resend.SendEmailRequest{
		From:    os.Getenv("RESEND_FROM"), // e.g. "Contact <onboarding@resend.dev>"
		To:      []string{contactEmail},
		Subject: subject,
		Html:    html,
		ReplyTo: email,
	}

	if params.From == "" {
		params.From = fmt.Sprintf("Contact Form <onboarding@resend.dev>")
	}

	_, err := client.Emails.Send(params)
	if err != nil {
		log.Printf("contact: failed to send email: %v", err)
		sendErrorHTML(w, "Failed to send message. Please try again later.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusAccepted)
	sendStyledHTML(w, "Message received. Thanks!")
}

func sendErrorHTML(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(code)
	sendStyledHTML(w, escapeHTML(msg))
}

func sendStyledHTML(w http.ResponseWriter, message string) {
	html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Contact</title>
  <style>
    * { box-sizing: border-box; }
    body { font-family: system-ui, -apple-system, sans-serif; margin: 0; min-height: 100vh; display: flex; align-items: center; justify-content: center; background: #f5f5f5; color: #1a1a1a; }
    @media (prefers-color-scheme: dark) { body { background: #1a1a1a; color: #e5e5e5; } }
    .card { background: #fff; padding: 2rem 2.5rem; border-radius: 12px; text-align: center; box-shadow: 0 2px 12px rgba(0,0,0,.08); max-width: 420px; }
    @media (prefers-color-scheme: dark) { .card { background: #2a2a2a; box-shadow: 0 2px 12px rgba(0,0,0,.3); } }
    .card p { margin: 0 0 1.25rem; font-size: 1.125rem; line-height: 1.5; }
    .card a { display: inline-block; padding: 0.6rem 1.25rem; background: #1a1a1a; color: #fff; text-decoration: none; border-radius: 8px; font-weight: 500; transition: opacity .2s; }
    .card a:hover { opacity: .9; }
    @media (prefers-color-scheme: dark) { .card a { background: #e5e5e5; color: #1a1a1a; } }
  </style>
</head>
<body>
  <div class="card">
    <p>%s</p>
    <a href="https://ark31.info/">Back to Site</a>
  </div>
</body>
</html>`, message)
	w.Write([]byte(html))
}

func escapeHTML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, `"`, "&quot;")
	s = strings.ReplaceAll(s, "'", "&#39;")
	return s
}

package turnstile

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestVerifySuccess(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s", r.Method)
		}
		if err := r.ParseForm(); err != nil {
			t.Fatal(err)
		}
		if r.Form.Get("secret") != "test-secret" {
			t.Errorf("secret = %q", r.Form.Get("secret"))
		}
		if r.Form.Get("response") != "tok" {
			t.Errorf("response = %q", r.Form.Get("response"))
		}
		if r.Form.Get("remoteip") != "1.2.3.4" {
			t.Errorf("remoteip = %q", r.Form.Get("remoteip"))
		}
		_ = json.NewEncoder(w).Encode(Result{Success: true, Hostname: "ark31.info"})
	}))
	defer srv.Close()

	c := &Client{
		Secret:     "test-secret",
		HTTPClient: srv.Client(),
		VerifyURL:  srv.URL,
	}
	res, err := c.Verify(context.Background(), "tok", "1.2.3.4")
	if err != nil {
		t.Fatal(err)
	}
	if !res.Success {
		t.Fatalf("expected success, got %#v", res)
	}
}

func TestVerifyMissingToken(t *testing.T) {
	c := &Client{Secret: "test-secret"}
	res, err := c.Verify(context.Background(), "  ", "")
	if err != nil {
		t.Fatal(err)
	}
	if res.Success {
		t.Fatal("expected failure for empty token")
	}
}

func TestVerifyMissingSecret(t *testing.T) {
	c := &Client{}
	_, err := c.Verify(context.Background(), "tok", "")
	if err == nil {
		t.Fatal("expected error when secret missing")
	}
}

func TestVerifySiteverifyFailure(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(Result{
			Success:    false,
			ErrorCodes: []string{"invalid-input-response"},
		})
	}))
	defer srv.Close()

	c := &Client{Secret: "s", HTTPClient: srv.Client(), VerifyURL: srv.URL}
	res, err := c.Verify(context.Background(), "bad", "")
	if err != nil {
		t.Fatal(err)
	}
	if res.Success {
		t.Fatal("expected success=false")
	}
}

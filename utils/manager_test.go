package utils

import (
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func newTestManager(rt roundTripFunc) *MailManager {
	return &MailManager{
		client: &http.Client{
			Timeout:   time.Second,
			Transport: rt,
		},
		color: &Color{},
	}
}

func jsonResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

func TestCreateAccountAPIStatusError(t *testing.T) {
	manager := newTestManager(func(req *http.Request) (*http.Response, error) {
		if req.URL.Path != "/accounts" {
			t.Fatalf("unexpected path: %s", req.URL.Path)
		}
		return jsonResponse(http.StatusUnprocessableEntity, `{"message":"address is invalid"}`), nil
	})

	id, err := manager.createAccountAPI("bad", "password")
	if err == nil {
		t.Fatal("expected error for non-201 response")
	}
	if id != "" {
		t.Fatalf("expected empty id, got %q", id)
	}
	if !strings.Contains(err.Error(), "status 422") {
		t.Fatalf("expected status in error, got: %v", err)
	}
	if !strings.Contains(err.Error(), "address is invalid") {
		t.Fatalf("expected API message in error, got: %v", err)
	}
}

func TestCreateAccountAPIMissingID(t *testing.T) {
	manager := newTestManager(func(req *http.Request) (*http.Response, error) {
		return jsonResponse(http.StatusCreated, `{"address":"x@example.com"}`), nil
	})

	_, err := manager.createAccountAPI("x@example.com", "password")
	if err == nil {
		t.Fatal("expected error when id is missing")
	}
	if !strings.Contains(err.Error(), "missing id") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetTokenMissingToken(t *testing.T) {
	manager := newTestManager(func(req *http.Request) (*http.Response, error) {
		if req.URL.Path != "/token" {
			t.Fatalf("unexpected path: %s", req.URL.Path)
		}
		return jsonResponse(http.StatusOK, `{"foo":"bar"}`), nil
	})

	_, err := manager.getToken("x@example.com", "password")
	if err == nil {
		t.Fatal("expected error when token is missing")
	}
	if !strings.Contains(err.Error(), "missing token") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestFetchMessagesAPIStatusError(t *testing.T) {
	manager := newTestManager(func(req *http.Request) (*http.Response, error) {
		return jsonResponse(http.StatusUnauthorized, `{"hydra:description":"JWT Token not found"}`), nil
	})

	_, err := manager.fetchMessagesAPI("bad-token")
	if err == nil {
		t.Fatal("expected error for unauthorized response")
	}
	if !strings.Contains(err.Error(), "status 401") {
		t.Fatalf("expected status code in error, got: %v", err)
	}
	if !strings.Contains(err.Error(), "JWT Token not found") {
		t.Fatalf("expected hydra message in error, got: %v", err)
	}
}

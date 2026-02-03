package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/CloudyKit/jet/v6"
	"github.com/suryansh74/sketlo/internal/config"
)

// Test Cases
// ==================================================
// TestAppRoutes
func TestAppRoutes(t *testing.T) {
	views := jet.NewSet(
		jet.NewOSFileSystemLoader("../../views"),
	)

	cfg := config.NewConfig(views)
	server := NewServer(cfg)

	t.Run("if user hit /check_health, then sc:200 rb:json", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/check_health", nil)
		res := httptest.NewRecorder()
		server.router.ServeHTTP(res, req)

		assertStatusCode(t, res.Code, http.StatusOK)

		var gotResponseBody map[string]string
		err := json.NewDecoder(res.Body).Decode(&gotResponseBody)
		assertNotError(t, err)

		assertResponseBody(t, gotResponseBody["message"], "working fine")
		assertHeader(t, res.Header().Get("Content-Type"), "application/json")
	})

	t.Run("if user hit /, then sc:200, rb:<form and name", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		res := httptest.NewRecorder()
		server.router.ServeHTTP(res, req)

		assertStatusCode(t, res.Code, http.StatusOK)
		body := res.Body.String()
		assertStringContains(t, body, "<form")
		assertStringContains(t, body, `name="username"`)
	})

	t.Run("Given a valid username is submitted", func(t *testing.T) {
		formData := url.Values{"username": {"Ronak"}}
		req := httptest.NewRequest(http.MethodPost, "/join", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		res := httptest.NewRecorder()
		server.router.ServeHTTP(res, req)

		assertStatusCode(t, res.Code, http.StatusSeeOther)
		responseLocation := res.Header().Get("Location")
		assertLocation(t, responseLocation, "/game?username=Ronak")
	})
}

// Assert Helper Function
// ==================================================
// assertStatusCode
func assertStatusCode(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}

func assertResponseBody(t *testing.T, got, want any) {
	t.Helper()
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func assertHeader(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func assertNotError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("got error:%s", err.Error())
	}
}

func assertStringContains(t *testing.T, body, match string) {
	t.Helper()
	if !strings.Contains(body, match) {
		t.Errorf("got %q, want %q", body, match)
	}
}

func assertLocation(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

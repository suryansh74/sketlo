package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAppRoutes(t *testing.T) {
	handler := SetupRoutes()
	t.Run("if user hit /check_health, sc:200 rb:json", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/check_health", nil)
		res := httptest.NewRecorder()
		handler.ServeHTTP(res, req)

		assertStatusCode(t, res.Code, http.StatusOK)

		var gotResponseBody map[string]string
		err := json.NewDecoder(res.Body).Decode(&gotResponseBody)
		assertNotError(t, err)

		assertResponseBody(t, gotResponseBody["message"], "working fine")
		assertHeader(t, res.Header().Get("Content-Type"), "application/json")
	})
}

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
		t.Errorf("got %s, want %s", got, want)
	}
}

func assertNotError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("got error:%s", err.Error())
	}
}

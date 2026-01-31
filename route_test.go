package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAppRoutes(t *testing.T) {
	t.Run("user hit home page endpoint", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		res := httptest.NewRecorder()
		SayHello(res, req)
		got := res.Code

		// checking if status code is expected
		want := http.StatusOK
		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}

		// checking if response body is expected
		gotBody := res.Body.String()
		wantBody := "hello world"
		if gotBody != wantBody {
			t.Errorf("got %s, want %s", gotBody, wantBody)
		}
	})
}

package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/CloudyKit/jet/v6"
	"github.com/gorilla/websocket"
	"github.com/suryansh74/sketlo/internal/chat"
	"github.com/suryansh74/sketlo/internal/config"
)

// Test Cases
// ==================================================
// TestAppRoutes
func TestAppRoutes(t *testing.T) {
	views := jet.NewSet(
		jet.NewOSFileSystemLoader("../../views"),
	)

	hub := chat.NewHub()
	go hub.Run()

	cfg := config.NewConfig(views, hub)
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
	t.Run("game endpoint should display username", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, `/game?username=Ronak`, nil)
		res := httptest.NewRecorder()
		server.router.ServeHTTP(res, req)

		assertStatusCode(t, res.Code, http.StatusOK)
		assertStringContains(t, res.Body.String(), "Ronak")
	})
	t.Run("check websocket upgrade for /ws endpoint", func(t *testing.T) {
		ts := httptest.NewServer(server.router)
		defer ts.Close()

		wsURL := "ws" + strings.TrimPrefix(ts.URL, "http") + "/ws"

		dialer := websocket.Dialer{}
		conn, res, err := dialer.Dial(wsURL, nil)
		assertNotError(t, err)
		defer conn.Close()

		assertStatusCode(t, res.StatusCode, http.StatusSwitchingProtocols)
	})
	t.Run("hub broadcasting working", func(t *testing.T) {
		ts := httptest.NewServer(server.router)
		defer ts.Close()

		wsURL := "ws" + strings.TrimPrefix(ts.URL, "http") + "/ws"

		dialer := websocket.Dialer{}

		aliceConn, _, err := dialer.Dial(wsURL, nil)
		assertNotError(t, err)
		defer aliceConn.Close()

		bobConn, _, err := dialer.Dial(wsURL, nil)
		assertNotError(t, err)
		defer bobConn.Close()

		message := []byte("Hello Sketlo")
		err = aliceConn.WriteMessage(websocket.TextMessage, message)
		assertNotError(t, err)

		_, got, err := bobConn.ReadMessage()
		assertNotError(t, err)

		assertEqual(t, string(got), string(message))
	})
	t.Run("hub broadcasting json payload", func(t *testing.T) {
		ts := httptest.NewServer(server.router)
		defer ts.Close()

		wsURL := "ws" + strings.TrimPrefix(ts.URL, "http") + "/ws"

		dialer := websocket.Dialer{}

		aliceConn, _, err := dialer.Dial(wsURL, nil)
		assertNotError(t, err)
		defer aliceConn.Close()

		bobConn, _, err := dialer.Dial(wsURL, nil)
		assertNotError(t, err)
		defer bobConn.Close()

		// alice sending data
		wsPayload := chat.WsPayload{
			Action:   "broadcast",
			Username: "Alice",
			Message:  "Hello JSON",
		}

		err = aliceConn.WriteJSON(wsPayload)
		assertNotError(t, err)

		// bob reads json same data as sended
		var received chat.WsPayload
		err = bobConn.ReadJSON(&received)
		assertNotError(t, err)

		assertEqual(t, received.Action, "broadcast")
		assertEqual(t, received.Message, "Hello JSON")
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
		t.Fatalf("got error:%s", err.Error())
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

func assertEqual(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

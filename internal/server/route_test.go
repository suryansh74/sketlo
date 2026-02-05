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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

		assert.Equal(t, res.Code, http.StatusOK)

		var gotResponseBody map[string]string
		err := json.NewDecoder(res.Body).Decode(&gotResponseBody)
		require.NoError(t, err)

		assert.Equal(t, gotResponseBody["message"], "working fine")
		assert.Equal(t, res.Header().Get("Content-Type"), "application/json")
	})

	t.Run("if user hit /, then sc:200, rb:<form and name", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		res := httptest.NewRecorder()
		server.router.ServeHTTP(res, req)

		assert.Equal(t, res.Code, http.StatusOK)
		body := res.Body.String()
		assert.Contains(t, body, "<form")
		assert.Contains(t, body, `name="username"`)
	})

	t.Run("Given a valid username is submitted", func(t *testing.T) {
		formData := url.Values{"username": {"Ronak"}}
		req := httptest.NewRequest(http.MethodPost, "/join", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		res := httptest.NewRecorder()
		server.router.ServeHTTP(res, req)

		assert.Equal(t, res.Code, http.StatusSeeOther)
		responseLocation := res.Header().Get("Location")
		assert.Equal(t, responseLocation, "/game?username=Ronak")
	})
	t.Run("game endpoint should display username", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, `/game?username=Ronak`, nil)
		res := httptest.NewRecorder()
		server.router.ServeHTTP(res, req)

		assert.Equal(t, res.Code, http.StatusOK)
		assert.Contains(t, res.Body.String(), "Ronak")
	})
	t.Run("check websocket upgrade for /ws endpoint", func(t *testing.T) {
		ts := httptest.NewServer(server.router)
		defer ts.Close()

		wsURL := "ws" + strings.TrimPrefix(ts.URL, "http") + "/ws"

		dialer := websocket.Dialer{}
		conn, res, err := dialer.Dial(wsURL, nil)
		require.NoError(t, err)
		defer conn.Close()

		assert.Equal(t, res.StatusCode, http.StatusSwitchingProtocols)
	})

	t.Run("hub broadcasting json payload", func(t *testing.T) {
		ts := httptest.NewServer(server.router)
		defer ts.Close()

		wsURL := "ws" + strings.TrimPrefix(ts.URL, "http") + "/ws"

		dialer := websocket.Dialer{}

		aliceConn, _, err := dialer.Dial(wsURL, nil)
		require.NoError(t, err)
		defer aliceConn.Close()

		bobConn, _, err := dialer.Dial(wsURL, nil)
		require.NoError(t, err)
		defer bobConn.Close()

		// alice sending data
		wsPayload := chat.WsPayload{
			Action:   "broadcast",
			Username: "Alice",
			Message:  "Hello JSON",
		}

		err = aliceConn.WriteJSON(wsPayload)
		require.NoError(t, err)

		// bob reads json same data as sended
		var received chat.WsPayload
		err = bobConn.ReadJSON(&received)
		require.NoError(t, err)

		assert.Equal(t, received.Action, "broadcast")
		assert.Equal(t, received.Message, "Hello JSON")
	})
}

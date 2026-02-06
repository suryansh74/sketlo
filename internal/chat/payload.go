package chat

type WsPayload struct {
	Action   string   `json:"action"`
	Username string   `json:"username,omitempty"`
	Message  string   `json:"message,omitempty"`
	Users    []string `json:"users,omitempty"`
}

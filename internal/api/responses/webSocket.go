package responses

type WSMessage struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

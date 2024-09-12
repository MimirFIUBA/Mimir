package models

type ItemsResponse struct {
	Status int   `json:"status"`
	Items  []any `json:"items"`
}

type WSMessage struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

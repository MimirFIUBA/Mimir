package models

type ItemsResponse struct {
	Status int   `json:"status"`
	Items  []any `json:"items"`
}

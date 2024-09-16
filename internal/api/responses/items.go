package responses

type ItemsResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Items   interface{} `json:"items"`
}

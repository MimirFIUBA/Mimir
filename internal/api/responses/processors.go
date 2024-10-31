package responses

type Handler struct {
	Name           string        `json:"name"`
	Topic          string        `json:"topic"`
	HandlerType    string        `json:"type"`
	SensorId       string        `json:"sensorId,omitempty"`
	Configurations []interface{} `json:"configurations,omitempty"`
}

type JsonConfiguration struct {
	Path string `json:"path"`
}

type BytesConfiguration struct {
	DataType   string `json:"dataType"`
	Endianness string `json:"endianness"`
	Size       int    `json:"size"`
}

func NewHandlerResponse(name, topic, handlerType string) *Handler {
	return &Handler{Name: name, Topic: topic, HandlerType: handlerType}
}

package responses

type Handler struct {
	Name           string        `json:"name"`
	Topic          string        `json:"topic"`
	HandlerType    string        `json:"type"`
	SensorId       string        `json:"sensorId,omitempty"`
	NodeId         string        `json:"nodeId,omitempty"`
	Configurations []interface{} `json:"configurations,omitempty"`
}

type JsonConfiguration struct {
	IdPosition     string                            `json:"idPosition"`
	Path           string                            `json:"path"`
	AdditionalData []JsonAdditionalDataConfiguration `json:"additionalData"`
}

type JsonAdditionalDataConfiguration struct {
	Name string `json:"name"`
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

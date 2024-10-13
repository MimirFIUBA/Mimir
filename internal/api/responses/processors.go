package responses

type Processor struct {
	Name           string        `json:"name"`
	Topic          string        `json:"topic"`
	ProcessorType  string        `json:"type"`
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

func NewProcessorResponse(name, topic, processorType string) *Processor {
	return &Processor{Name: name, Topic: topic, ProcessorType: processorType}
}

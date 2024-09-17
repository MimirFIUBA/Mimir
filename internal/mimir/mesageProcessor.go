package mimir

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io"
	"math"
	"strings"
	"time"

	mimir "mimir/internal/mimir/models"
)

type MessageProcessor interface {
	ProcessMessage(topic string, payload []byte) error
}

type ProcessorRegistry struct {
	processors map[string]MessageProcessor
}

func NewProcessorRegistry() *ProcessorRegistry {
	return &ProcessorRegistry{processors: make(map[string]MessageProcessor)}
}

func (r *ProcessorRegistry) RegisterProcessor(topic string, processor MessageProcessor) {
	r.processors[topic] = processor
}

func (r *ProcessorRegistry) GetProcessor(topic string) (MessageProcessor, bool) {
	processor, exists := r.processors[topic]
	return processor, exists
}
func (r *ProcessorRegistry) GetProcessors() map[string]MessageProcessor {
	return r.processors
}

type JSONProcessor struct {
	SensorId                string
	jsonValueConfigurations []JSONValueConfiguration
}

type JSONValueConfiguration struct {
	idPosition    string
	valuePosition string
}

func NewJSONValueConfiguration(idPath, valuePath string) *JSONValueConfiguration {
	return &JSONValueConfiguration{idPath, valuePath}
}

func NewJSONProcessor() *JSONProcessor {
	return &JSONProcessor{}
}

func (p *JSONProcessor) AddValueConfiguration(configuration *JSONValueConfiguration) {
	p.jsonValueConfigurations = append(p.jsonValueConfigurations, *configuration)
}

func (p *JSONProcessor) ProcessMessage(topic string, payload []byte) error {
	var jsonPayload = string(payload)
	jsonDataReader := strings.NewReader(jsonPayload)
	decoder := json.NewDecoder(jsonDataReader)
	var jsonMap map[string]interface{}
	for {
		err := decoder.Decode(&jsonMap)
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
	}

	var sensorId string
	if p.SensorId != "" {
		sensorId = p.SensorId
	}

	for _, configuration := range p.jsonValueConfigurations {
		if configuration.idPosition != "" {
			idInterface, ok := getValueFromJSON(jsonMap, configuration.idPosition)
			if !ok {
				return ValueNotFoundError{"sensorId"}
			}
			sensorIdValue, ok := idInterface.(string)
			if !ok {
				return WrongFormatError{"sensorId"}
			}
			sensorId = sensorIdValue
		}

		valueInterface, ok := getValueFromJSON(jsonMap, configuration.valuePosition)
		if !ok {
			return ValueNotFoundError{"sensorId"}
		}

		sensorReading := mimir.SensorReading{SensorID: sensorId, Value: valueInterface, Time: time.Now()}
		Manager.readingsChannel <- sensorReading
	}
	return nil
}

type BytesProcessor struct {
	SensorId            string
	BytesConfigurations []BytesConfiguration
}

func NewBytesProcessor() *BytesProcessor {
	return &BytesProcessor{"", nil}
}

type BytesConfiguration struct {
	DataType   string
	Endianness binary.ByteOrder
	Size       int //size in bytes
}

func NewBytesConfiguration(dataType string, endianess binary.ByteOrder, size int) *BytesConfiguration {
	return &BytesConfiguration{dataType, endianess, size}
}

func (p *BytesProcessor) AddBytesConfiguration(configuration BytesConfiguration) {
	p.BytesConfigurations = append(p.BytesConfigurations, configuration)
}

func readBool(stream *bytes.Reader) (bool, error) {
	var value byte
	err := binary.Read(stream, binary.BigEndian, &value)
	if err != nil {
		return false, err
	}
	return value != 0, nil
}

func (p *BytesProcessor) ProcessMessage(topic string, payload []byte) error {
	var sensorId string
	if p.SensorId != "" {
		sensorId = p.SensorId
	}

	i := 0
	for _, configuration := range p.BytesConfigurations {
		dataBytes := payload[i : configuration.Size+i]
		var data interface{}
		switch configuration.DataType {
		case "id":
			sensorId = string(dataBytes)
		case "string":
			data = string(dataBytes)
		case "int":
			data = configuration.Endianness.Uint32(dataBytes)
		case "float":
			data = math.Float32frombits(configuration.Endianness.Uint32(dataBytes))
		case "bool":
			stream := bytes.NewReader(dataBytes)
			value, err := readBool(stream)
			if err != nil {
				panic("Fail to read bool")
			}
			data = value
		}

		if sensorId != "" && configuration.DataType != "id" {
			sensorReading := mimir.SensorReading{SensorID: sensorId, Value: data, Time: time.Now()}
			Manager.readingsChannel <- sensorReading
		}

		i += configuration.Size
	}
	return nil
}

type XMLProcessor struct{}

func NewXMLProcessor() *XMLProcessor {
	return &XMLProcessor{}
}

func (p *XMLProcessor) ProcessMessage(topic string, payload []byte) error {
	//TODO
	panic("Missing implementation")
}

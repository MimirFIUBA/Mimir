package processors

import (
	"bytes"
	"encoding/binary"
	"math"
	"mimir/internal/consts"
	"mimir/internal/mimir/models"
	"strings"
	"time"
)

type BytesProcessor struct {
	SensorId            string
	Name                string
	Topic               string
	Type                string
	BytesConfigurations []BytesConfiguration
	ReadingsChannel     chan models.SensorReading
}

func NewBytesProcessor() *BytesProcessor {
	return &BytesProcessor{}
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

func (p *BytesProcessor) SetReadingsChannel(readingsChannel chan models.SensorReading) {
	p.ReadingsChannel = readingsChannel
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
			sensorReading := models.SensorReading{SensorID: sensorId, Value: data, Time: time.Now(), Topic: topic}
			p.ReadingsChannel <- sensorReading
		}

		i += configuration.Size
	}
	return nil
}

func (p *BytesProcessor) GetConfigFilename() string {
	return strings.ReplaceAll(p.Topic, "/", "_") + consts.PROCESSORS_FILE_SUFFIX
}

func (p *BytesProcessor) GetTopic() string {
	return p.Topic
}

func (p *BytesProcessor) GetType() ProcessorType {
	return BYTES_PROCESSOR
}

func (p *BytesProcessor) UpdateFields(map[string]interface{}) error {
	return nil
}

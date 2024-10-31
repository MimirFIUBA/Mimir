package handlers

import (
	"bytes"
	"encoding/binary"
	"math"
	"mimir/internal/consts"
	"mimir/internal/models"
	"strings"
	"time"
)

type BytesHandler struct {
	SensorId            string
	Name                string
	Topic               string
	Type                string
	BytesConfigurations []BytesConfiguration
	ReadingsChannel     chan models.SensorReading
}

func NewBytesHandler() *BytesHandler {
	return &BytesHandler{}
}

type BytesConfiguration struct {
	DataType   string
	Endianness binary.ByteOrder
	Size       int //size in bytes
}

func NewBytesConfiguration(dataType string, endianess binary.ByteOrder, size int) *BytesConfiguration {
	return &BytesConfiguration{dataType, endianess, size}
}

func (p *BytesHandler) AddBytesConfiguration(configuration BytesConfiguration) {
	p.BytesConfigurations = append(p.BytesConfigurations, configuration)
}

func (p *BytesHandler) SetReadingsChannel(readingsChannel chan models.SensorReading) {
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

func (p *BytesHandler) HandleMessage(msg Message) error {
	var sensorId string
	if p.SensorId != "" {
		sensorId = p.SensorId
	}

	i := 0
	for _, configuration := range p.BytesConfigurations {
		dataBytes := msg.Payload[i : configuration.Size+i]
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
			sensorReading := models.SensorReading{SensorID: sensorId, Value: data, Time: time.Now(), Topic: msg.Topic}
			p.ReadingsChannel <- sensorReading
		}

		i += configuration.Size
	}
	return nil
}

func (p *BytesHandler) GetConfigFilename() string {
	return strings.ReplaceAll(p.Topic, "/", "_") + consts.HANDLERS_FILE_SUFFIX
}

func (p *BytesHandler) GetTopic() string {
	return p.Topic
}

func (p *BytesHandler) GetType() HandlerType {
	return BYTES_HANDLER
}

func (p *BytesHandler) UpdateFields(map[string]interface{}) error {
	return nil
}
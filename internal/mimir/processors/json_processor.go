package processors

import (
	"encoding/json"
	"fmt"
	"io"
	"mimir/internal/consts"
	"mimir/internal/models"
	"mimir/internal/utils"
	"strings"
	"time"
)

type JSONProcessor struct {
	SensorId                string                    `json:"sensorId,omitempty"`
	Name                    string                    `json:"name"`
	Topic                   string                    `json:"topic"`
	Type                    string                    `json:"type"`
	JsonValueConfigurations []JSONValueConfiguration  `json:"configurations"`
	ReadingsChannel         chan models.SensorReading `json:"-"`
}

type JSONValueConfiguration struct {
	IdPosition string `json:"idPosition,omitempty"`
	ValuePath  string `json:"path,omitempty"`
}

func NewJSONValueConfiguration(idPath, valuePath string) *JSONValueConfiguration {
	return &JSONValueConfiguration{idPath, valuePath}
}

func NewJSONProcessor() *JSONProcessor {
	return &JSONProcessor{}
}

func (p *JSONProcessor) SetReadingsChannel(readingsChannel chan models.SensorReading) {
	p.ReadingsChannel = readingsChannel
}

func (p *JSONProcessor) AddValueConfiguration(configuration *JSONValueConfiguration) {
	p.JsonValueConfigurations = append(p.JsonValueConfigurations, *configuration)
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

	for _, configuration := range p.JsonValueConfigurations {
		if configuration.IdPosition != "" {
			idInterface, ok := utils.GetValueFromJSON(jsonMap, configuration.IdPosition)
			if !ok {
				return ValueNotFoundError{"sensorId"}
			}
			sensorIdValue, ok := idInterface.(string)
			if !ok {
				return WrongFormatError{"sensorId"}
			}
			sensorId = sensorIdValue
		}

		valueInterface, ok := utils.GetValueFromJSON(jsonMap, configuration.ValuePath)
		if !ok {
			return ValueNotFoundError{configuration.ValuePath}
		}

		sensorReading := models.SensorReading{SensorID: sensorId, Value: valueInterface, Time: time.Now(), Topic: topic}
		p.ReadingsChannel <- sensorReading
	}
	return nil
}

func (p *JSONProcessor) GetConfigFilename() string {
	return strings.ReplaceAll(p.Topic, "/", "_") + consts.PROCESSORS_FILE_SUFFIX
}

func (p *JSONProcessor) GetTopic() string {
	return p.Topic
}

func (p *JSONProcessor) GetType() ProcessorType {
	return JSON_PROCESSOR
}

func (p *JSONProcessor) UpdateFields(fieldsToUpdate map[string]interface{}) error {
	for k, v := range fieldsToUpdate {
		switch k {
		case "name":
			name, ok := v.(string)
			if !ok {
				return fmt.Errorf("name is not a string")
			}
			p.Name = name
		}
	}
	return nil
}

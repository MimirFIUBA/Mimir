package api

import (
	"encoding/binary"
	mimir "mimir/internal/mimir"
)

func jsonToProcessor(processorType string, jsonMap map[string]interface{}) (mimir.MessageProcessor, error) {

	switch processorType {
	case "bytes":
		return jsonMapToBytesProcessor(jsonMap)
	case "json":
		return jsonToJsonProcessor()
	case "xml":
		return jsonToXMLProcessor()
	default:
		return nil, nil
	}
}

type BytesProcessor struct {
	*mimir.BytesProcessor
}

func (p *BytesProcessor) setSensorId(jsonMap map[string]interface{}) bool {
	sensorIdValue, exists := jsonMap["sensorId"]
	if exists {
		sensorId, ok := sensorIdValue.(string)
		if !ok {
			return false
		}
		p.SensorId = sensorId
	}

	return true
}

func (p *BytesProcessor) setConfigurations(jsonMap map[string]interface{}) error {
	configurationsValue, exists := jsonMap["configurations"]
	if !exists {
		return RequiredFieldError{"configurations"}
	}

	configurations, ok := configurationsValue.([]interface{})
	if !ok {
		return WrongFormatError{"configurations"}
	}

	for _, configurationInterface := range configurations {
		configurationValue, ok := configurationInterface.(map[string]interface{})
		if !ok {
			return WrongFormatError{"byteConfiguration"}
		}
		configuration, err := jsonMapToByteConfiguration(configurationValue)
		if err != nil {
			return err
		}
		p.BytesConfigurations = append(p.BytesConfigurations, *configuration)
	}

	return nil
}

func jsonMapToByteConfiguration(jsonMap map[string]interface{}) (*mimir.BytesConfiguration, error) {
	dataTypeInterface, exists := jsonMap["dataType"]
	if !exists {
		return nil, RequiredFieldError{"dataType"}
	}
	dataTypeValue, ok := dataTypeInterface.(string)
	if !ok {
		return nil, WrongFormatError{"dataType"}
	}

	endiannessInterface, exists := jsonMap["endianness"]
	if !exists {
		return nil, RequiredFieldError{"endianness"}
	}
	endiannessValue, ok := endiannessInterface.(string)
	if !ok {
		return nil, WrongFormatError{"endianness"}
	}

	var byteOrder binary.ByteOrder
	switch endiannessValue {
	case "littleEndian":
		byteOrder = binary.LittleEndian
	case "bigEndian":
		byteOrder = binary.BigEndian
	case "nativeEndian":
		byteOrder = binary.NativeEndian
	default:
		return nil, WrongFormatError{"endianness"}
	}

	sizeInterface, exists := jsonMap["size"]
	if !exists {
		return nil, RequiredFieldError{"size"}
	}
	sizeValue, ok := sizeInterface.(float64)
	if !ok {
		return nil, WrongFormatError{"size"}
	}
	//TODO: validate that size is an int (not float)

	return mimir.NewBytesConfiguration(dataTypeValue, byteOrder, int(sizeValue)), nil
}

func jsonMapToBytesProcessor(jsonMap map[string]interface{}) (mimir.MessageProcessor, error) {
	bytesProcessor := BytesProcessor{mimir.NewBytesProcessor()}

	ok := bytesProcessor.setSensorId(jsonMap)
	if !ok {
		return nil, WrongFormatError{"sensorId"}
	}

	err := bytesProcessor.setConfigurations(jsonMap)
	if err != nil {
		return nil, err
	}

	return bytesProcessor, nil

}

func jsonToJsonProcessor() (mimir.MessageProcessor, error) {
	processor := mimir.NewJSONProcessor()
	return processor, nil
}

func jsonToXMLProcessor() (mimir.MessageProcessor, error) {
	processor := mimir.NewXMLProcessor()
	return processor, nil
}

package processors

import (
	"encoding/binary"
)

func JsonToProcessor(jsonMap map[string]interface{}) (MessageProcessor, error) {
	processorType, ok := jsonMap["type"].(string)
	if !ok {
		panic("bad configuration")
	}

	switch processorType {
	case "bytes":
		return jsonMapToBytesProcessor(jsonMap)
	case "json":
		return jsonToJsonProcessor(jsonMap)
	case "xml":
		return jsonToXMLProcessor()
	default:
		return nil, nil
	}
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

func jsonMapToByteConfiguration(jsonMap map[string]interface{}) (*BytesConfiguration, error) {
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

	return NewBytesConfiguration(dataTypeValue, byteOrder, int(sizeValue)), nil
}

func jsonMapToBytesProcessor(jsonMap map[string]interface{}) (MessageProcessor, error) {
	bytesProcessor := NewBytesProcessor()
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

func jsonToJsonProcessor(jsonMap map[string]interface{}) (MessageProcessor, error) {
	processor := NewJSONProcessor()

	configurationsValue, exists := jsonMap["configurations"]
	if !exists {
		return nil, RequiredFieldError{"configurations"}
	}

	configurations, ok := configurationsValue.([]interface{})
	if !ok {
		return nil, WrongFormatError{"configurations"}
	}

	for _, configurationInterface := range configurations {
		configurationValue, ok := configurationInterface.(map[string]interface{})
		if !ok {
			return nil, WrongFormatError{"byteConfiguration"}
		}
		configuration, err := jsonMapToJsonConfiguration(configurationValue)
		if err != nil {
			return nil, err
		}
		processor.AddValueConfiguration(configuration)
	}
	return processor, nil
}

func jsonMapToJsonConfiguration(jsonMap map[string]interface{}) (*JSONValueConfiguration, error) {
	configuration := &JSONValueConfiguration{}
	pathInterface, exists := jsonMap["path"]
	if exists {
		path, ok := pathInterface.(string)
		if !ok {
			return nil, WrongFormatError{"Path is not a string"}
		}
		configuration.ValuePath = path
	}

	return configuration, nil

}

func jsonToXMLProcessor() (MessageProcessor, error) {
	processor := NewXMLProcessor()
	return processor, nil
}

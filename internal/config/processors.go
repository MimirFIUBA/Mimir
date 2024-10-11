package config

import (
	"fmt"
	"mimir/internal/db"
	"mimir/internal/mimir"
	"mimir/internal/mimir/models"
	"mimir/internal/utils"

	"github.com/gookit/config"
)

func BuildProcessors(mimirProcessor *mimir.MimirProcessor) {
	configuration := config.Data()
	processors, ok := configuration["processors"].([]interface{})
	if !ok {
		panic("Wrong processors format, no processors")
	}

	sensors := make([]*models.Sensor, 0)
	for _, processorInterface := range processors {
		processorMap, ok := processorInterface.(map[string]interface{})
		if !ok {
			panic("bad configuration")
		}

		topic, ok := processorMap["topic"].(string)
		if !ok {
			panic("bad configuration")
		}

		processorType, ok := processorMap["type"].(string)
		if !ok {
			panic("bad configuration")
		}

		processor, err := utils.JsonToProcessor(processorType, processorMap)
		if err != nil {
			fmt.Println(err)
			panic("bad configuration")
		}

		mimir.MessageProcessors.RegisterProcessor(topic, processor)
		sensor := models.NewSensor(topic)
		sensor.Topic = topic
		mimirProcessor.RegisterSensor(sensor)
		sensors = append(sensors, sensor)
		// db.SensorsData.CreateSensor(sensor)
	}
	db.SensorsData.LoadSensors(sensors)

}

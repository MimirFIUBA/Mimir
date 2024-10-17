package config

import (
	"encoding/json"
	"fmt"
	"log"
	"mimir/internal/consts"
	"mimir/internal/db"
	"mimir/internal/mimir"
	"mimir/internal/mimir/models"
	"mimir/internal/mimir/processors"
	"mimir/internal/utils"
	"os"

	"github.com/gookit/ini/v2"
)

func BuildProcessors(mimirProcessor *mimir.MimirProcessor) {
	dir := ini.String(consts.PROCESSORS_DIR_CONFIG_NAME)
	files := utils.ListFilesWithSuffix(dir, "*"+consts.PROCESSORS_FILE_SUFFIX)
	sensors := make([]models.Sensor, 0)
	for _, v := range files {
		byteValue, err := os.ReadFile(v)
		if err != nil {
			log.Fatal(err)
			return
		}
		var jsonMap map[string]interface{}
		json.Unmarshal(byteValue, &jsonMap)

		topic, ok := jsonMap["topic"].(string)
		if !ok {
			panic("bad configuration")
		}

		processor, err := processors.JsonToProcessor(jsonMap)
		if err != nil {
			fmt.Println(err)
			panic("bad configuration")
		}

		processor.SetReadingsChannel(mimirProcessor.ReadingChannel)
		mimir.MessageProcessors.RegisterProcessor(topic, processor)
		sensor := models.NewSensor(topic)
		sensor.Topic = topic
		mimirProcessor.RegisterSensor(sensor)
		sensors = append(sensors, *sensor)
	}
	db.SensorsData.LoadSensors(sensors)
}

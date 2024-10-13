package config

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"mimir/internal/db"
	"mimir/internal/mimir"
	"mimir/internal/mimir/models"
	"mimir/internal/mimir/processors"
	"os"
	"path"

	"github.com/gookit/ini/v2"
)

func listFiles(dir string) []string {
	root := os.DirFS(dir)

	mdFiles, err := fs.Glob(root, "*.json")

	if err != nil {
		log.Fatal(err)
	}

	var files []string
	for _, v := range mdFiles {
		files = append(files, path.Join(dir, v))
	}
	return files
}

func BuildProcessors(mimirProcessor *mimir.MimirProcessor) {

	dir := ini.String("processors_dir")

	files := listFiles(dir)

	sensors := make([]*models.Sensor, 0)
	for _, v := range files {
		byteValue, err := os.ReadFile(v)
		if err != nil {
			log.Fatal(err)
			return
		}
		var jsonMap map[string]interface{}
		json.Unmarshal(byteValue, &jsonMap)
		fmt.Println(jsonMap)

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
		sensors = append(sensors, sensor)
	}

	// configuration := config.Data()
	// processorsArray, ok := configuration["processors"].([]interface{})
	// if !ok {
	// 	panic("Wrong processors format, no processors")
	// }

	// sensors := make([]*models.Sensor, 0)
	// for _, processorInterface := range processorsArray {
	// 	processorMap, ok := processorInterface.(map[string]interface{})
	// 	if !ok {
	// 		panic("bad configuration")
	// 	}

	// 	topic, ok := processorMap["topic"].(string)
	// 	if !ok {
	// 		panic("bad configuration")
	// 	}

	// 	processor, err := processors.JsonToProcessor(processorMap)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		panic("bad configuration")
	// 	}

	// 	processor.SetReadingsChannel(mimirProcessor.ReadingChannel)
	// 	mimir.MessageProcessors.RegisterProcessor(topic, processor)
	// 	sensor := models.NewSensor(topic)
	// 	sensor.Topic = topic
	// 	mimirProcessor.RegisterSensor(sensor)
	// 	sensors = append(sensors, sensor)
	// }
	db.SensorsData.LoadSensors(sensors)
}

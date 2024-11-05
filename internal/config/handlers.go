package config

import (
	"encoding/json"
	"log/slog"
	"mimir/internal/consts"
	"mimir/internal/db"
	"mimir/internal/handlers"
	"mimir/internal/mimir"
	"mimir/internal/models"
	"mimir/internal/utils"
	"os"

	"github.com/gookit/ini/v2"
)

func BuildHandlers(mimirEngine *mimir.MimirEngine) {
	dir := ini.String(consts.HANDLERS_DIR_CONFIG_NAME)
	files := utils.ListFilesWithSuffix(dir, "*"+consts.HANDLERS_FILE_SUFFIX)
	sensors := make([]models.Sensor, 0)
	for _, v := range files {
		slog.Info("building handler", "file", v)
		byteValue, err := os.ReadFile(v)
		if err != nil {
			slog.Error("error building handler", "file", v, "error", err)
			continue
		}
		var jsonMap map[string]interface{}
		json.Unmarshal(byteValue, &jsonMap)

		topic, ok := jsonMap["topic"].(string)
		if !ok {
			slog.Error("error building handler", "file", v, "error", "could not convert topic to string")
			continue
		}

		processor, err := handlers.JsonToHandler(jsonMap)
		if err != nil {
			slog.Error("error building handler", "file", v, "error", err)
			continue
		}

		processor.SetReadingsChannel(mimirEngine.ReadingChannel)
		mimir.Mimir.MsgProcessor.RegisterHandler(topic, processor)
		sensor := models.NewSensor(topic)
		sensor.Topic = topic
		mimirEngine.RegisterSensor(sensor)
		sensors = append(sensors, *sensor)
	}

	db.SensorsData.LoadSensors(sensors)
}

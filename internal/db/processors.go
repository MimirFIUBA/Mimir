package db

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"mimir/internal/consts"
	"mimir/internal/mimir/processors"
	"os"
	"strings"

	"github.com/gookit/ini/v2"
)

func (d *DatabaseManager) SaveProcessor(processor processors.MessageProcessor) {
	jsonString, err := json.MarshalIndent(processor, "", "    ")
	if err != nil {
		fmt.Println("Error ", err)
	}

	fileName := ini.String(consts.PROCESSORS_DIR_CONFIG_NAME) + "/" + processor.GetConfigFilename()

	os.WriteFile(fileName, jsonString, os.ModePerm)
}

func (d *DatabaseManager) DeleteProcessor(processor processors.MessageProcessor) {
	fileName := ini.String(consts.PROCESSORS_DIR_CONFIG_NAME) + "/" + processor.GetConfigFilename()
	newName := strings.Replace(fileName, ".json", "_deleted.json", 1)
	err := os.Rename(fileName, newName)
	if err != nil {
		slog.Error("error renaming file for deletion", "error", err)
	}
}

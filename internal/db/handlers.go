package db

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"mimir/internal/consts"
	"mimir/internal/handlers"
	"os"
	"strings"

	"github.com/gookit/ini/v2"
)

func (d *DatabaseManager) SaveHandler(handler handlers.MessageHandler) {
	jsonString, err := json.MarshalIndent(handler, "", "    ")
	if err != nil {
		fmt.Println("Error ", err)
	}

	fileName := ini.String(consts.HANDLERS_DIR_CONFIG_NAME) + "/" + handler.GetConfigFilename()

	os.WriteFile(fileName, jsonString, os.ModePerm)
}

func (d *DatabaseManager) DeleteHandler(handler handlers.MessageHandler) {
	fileName := ini.String(consts.HANDLERS_DIR_CONFIG_NAME) + "/" + handler.GetConfigFilename()
	newName := strings.Replace(fileName, ".json", "_deleted.json", 1)
	err := os.Rename(fileName, newName)
	if err != nil {
		slog.Error("error renaming file for deletion", "error", err)
	}
}

package config

import (
	"encoding/json"
	"log/slog"
	"mimir/internal/consts"
	"mimir/internal/db"
	"mimir/internal/utils"
	"os"

	"github.com/gookit/ini/v2"
)

func LoadVariables() {
	dir := ini.String(consts.VARIABLES_DIR_CONFIG_NAME)
	files := utils.ListFilesWithSuffix(dir, "*"+consts.VARIABLES_FILE_SUFFIX)
	variablesToLoad := make([]*db.UserVariable, 0)
	if len(files) > 0 {
		for _, filename := range files {
			slog.Info("loading variables from file", "file", filename)
			byteValue, err := os.ReadFile(filename)
			if err != nil {
				slog.Error("error reading file", "file", filename)
			}
			var variables []*db.UserVariable
			json.Unmarshal(byteValue, &variables)
			for _, variable := range variables {
				variable.Filename = filename
			}
			variablesToLoad = append(variablesToLoad, variables...)
		}
	}

	slog.Info("variables to load", "variables", variablesToLoad)
	db.AddUserVariables(variablesToLoad...)
}

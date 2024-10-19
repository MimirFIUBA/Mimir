package config

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"mimir/internal/consts"
	"mimir/internal/db"
	"mimir/internal/mimir"
	"mimir/internal/utils"
	"mimir/triggers"
	"os"

	"github.com/gookit/ini/v2"
)

func BuildTriggers(mimirProcessor *mimir.MimirProcessor) error {
	dir := ini.String(consts.TRIGGERS_DIR_CONFIG_NAME)
	files := utils.ListFilesWithSuffix(dir, "*"+consts.TRIGGERS_FILE_SUFFIX)

	if len(files) > 0 {
		triggersToUpsert := make([]db.Trigger, 0)

		for _, filename := range files {
			byteValue, err := os.ReadFile(filename)
			if err != nil {
				log.Fatal(err)
				return fmt.Errorf("error reading trigger file %s", filename)
			}
			var triggerData db.Trigger
			json.Unmarshal(byteValue, &triggerData)
			triggerData.Filename = filename
			triggersToUpsert = append(triggersToUpsert, triggerData)

			db.RegisterTrigger(&triggerData)
		}

		_, err := db.Database.UpsertTriggers(triggersToUpsert)
		if err != nil {
			slog.Error("error upserting triggers", "error", err)
			return err
		}

	}

	return nil
}

func BuildTriggerFromMap(triggerMap map[string]interface{}, mimirProcessor *mimir.MimirProcessor) *triggers.Trigger {
	trigger := buildTrigger(triggerMap)
	condition, exists := buildCondition(triggerMap)
	if exists {
		trigger.Condition = condition
	}
	actions := buildActions(triggerMap, mimirProcessor)
	for _, action := range actions {
		trigger.AddAction(action)
	}

	return trigger
}

func buildCondition(triggerMap map[string]interface{}) (triggers.Condition, bool) {
	conditionConfiguration, exists := triggerMap["condition"]
	if exists {
		var condition triggers.Condition
		switch conditionValue := conditionConfiguration.(type) {
		case string:
			condition = triggers.BuildConditionFromString(conditionValue)
		case map[string]interface{}:
			condition = buildConditionFromMap(conditionValue)
		}
		return condition, true
	}
	return nil, false
}

func buildConditionFromMap(_ map[string]interface{}) triggers.Condition {
	panic("missing implementation")
}

func buildTrigger(triggerMap map[string]interface{}) *triggers.Trigger {
	nameValue, exists := triggerMap["name"]
	if !exists {
		panic("Missing name for trigger")
	}

	triggerName, ok := nameValue.(string)
	if !ok {
		panic("Trigger name is not a string")
	}
	return triggers.NewTrigger(triggerName)
}

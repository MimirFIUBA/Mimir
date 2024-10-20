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

			trigger := BuildTriggerObserver(triggerData, mimirProcessor)
			db.RegisterTrigger(trigger, triggerData.Topics)
		}

		_, err := db.Database.UpsertTriggers(triggersToUpsert)
		if err != nil {
			slog.Error("error upserting triggers", "error", err)
			return err
		}

	}

	return nil
}

func BuildTriggerObserver(t db.Trigger, mimirProcessor *mimir.MimirProcessor) triggers.TriggerObserver {
	trigger := triggers.NewTrigger(t.Name)
	trigger.Condition = triggers.BuildConditionFromString(string(t.Condition))
	for _, action := range t.Actions {
		triggerAction := ToTriggerAction(action, mimirProcessor)
		trigger.AddAction(triggerAction)
	}

	return trigger
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

func ToTriggerAction(a db.Action, mimirProcessor *mimir.MimirProcessor) triggers.Action {
	var triggerAction triggers.Action

	switch a.Type {
	case "print":
		action := triggers.NewPrintAction()
		action.Message = a.Message
		action.Name = a.Name
		triggerAction = action
	case "alert":
		action := mimirProcessor.NewSendMQTTMessageAction(a.Message)
		action.Message = a.Message
		triggerAction = &action
	case "webSocket":
		action := mimirProcessor.NewSendWebSocketMessageAction(a.Message)
		action.Message = a.Message
		triggerAction = &action
	default:
		slog.Warn("action type not recognized while creating trigger action", "type", a.Type)
	}

	return triggerAction
}

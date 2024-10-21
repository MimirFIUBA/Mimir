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
	"go.mongodb.org/mongo-driver/bson"
)

func BuildTriggers(mimirProcessor *mimir.MimirProcessor) error {
	dir := ini.String(consts.TRIGGERS_DIR_CONFIG_NAME)
	files := utils.ListFilesWithSuffix(dir, "*"+consts.TRIGGERS_FILE_SUFFIX)

	if len(files) > 0 {
		triggersByFilename := make(map[string]triggers.TriggerObserver)
		topicsByFilename := make(map[string][]string)
		triggersToUpsert := make([]db.Trigger, 0)

		for _, filename := range files {
			slog.Info("Building trigger", "trigger", filename)
			byteValue, err := os.ReadFile(filename)
			if err != nil {
				log.Fatal(err)
				return fmt.Errorf("error reading trigger file %s", filename)
			}
			var triggerData db.Trigger
			json.Unmarshal(byteValue, &triggerData)
			triggerData.Filename = filename
			triggersToUpsert = append(triggersToUpsert, triggerData)

			trigger, err := BuildTriggerObserver(triggerData, mimirProcessor)
			if err == nil {
				triggersByFilename[filename] = trigger
				topicsByFilename[filename] = triggerData.Topics
			} else {
				slog.Error("Error creating trigger", "error", err)
			}
		}

		_, err := db.Database.UpsertTriggers(triggersToUpsert)
		if err != nil {
			slog.Error("error upserting triggers", "error", err)
			return err
		}

		//TODO: doing this to set id on trigger. See if we can improve (or at least take it somewhere else)
		values := bson.A{}
		for filename := range triggersByFilename {
			values = append(values, filename)
		}
		filter := bson.D{{Key: "filename", Value: bson.D{{Key: "$in", Value: values}}}}
		dbTriggers, err := db.Database.FindTriggers(filter)
		if err != nil {
			return err
		}
		for _, dbTriggerData := range dbTriggers {
			triggerToUpdate, exists := triggersByFilename[dbTriggerData.Filename]
			if exists {
				triggerToUpdate.SetID(dbTriggerData.ID.Hex())
				topics, exists := topicsByFilename[dbTriggerData.Filename]
				db.TriggerFilenamesById[dbTriggerData.ID.Hex()] = dbTriggerData.Filename
				if exists {
					db.RegisterTrigger(triggerToUpdate, topics)
				}
			}
		}
	}
	return nil
}

func BuildTriggerObserver(t db.Trigger, mimirProcessor *mimir.MimirProcessor) (triggers.TriggerObserver, error) {
	trigger := triggers.NewTrigger(t.Name)
	trigger.SetID(t.ID.Hex())
	trigger.IsActive = t.IsActive
	condition, err := triggers.BuildConditionFromString(string(t.Condition))
	if err != nil {
		return nil, err
	}
	trigger.Condition = condition
	BuildActions(t, trigger, mimirProcessor)

	return trigger, nil
}

func BuildActions(triggerData db.Trigger, trigger triggers.TriggerObserver, mimirProcessor *mimir.MimirProcessor) {
	for _, action := range triggerData.Actions {
		triggerAction := ToTriggerAction(action, mimirProcessor)
		trigger.AddAction(triggerAction)
	}
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
			conditionBuilt, err := triggers.BuildConditionFromString(conditionValue)
			if err != nil {
				return nil, false //TODO ver si tenemos que devolver el error
			}
			condition = conditionBuilt
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
		//TODO see if returning nil is fine or we need some error here
		slog.Warn("action type not recognized while creating trigger action", "type", a.Type)
	}

	return triggerAction
}

package config

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"mimir/internal/consts"
	"mimir/internal/db"
	"mimir/internal/mimir"
	"mimir/internal/models"
	"mimir/internal/utils"
	"mimir/triggers"
	"os"
	"time"

	"github.com/gookit/ini/v2"
	"go.mongodb.org/mongo-driver/bson"
)

var triggerTypeByName = map[string]triggers.TriggerType{
	"event":     triggers.EVENT_TRIGGER,
	"timer":     triggers.TIMER_TRIGGER,
	"frequency": triggers.FREQUENCY_TRIGGER,
}

func BuildTriggers(mimirProcessor *mimir.MimirProcessor) error {
	dir := ini.String(consts.TRIGGERS_DIR_CONFIG_NAME)
	files := utils.ListFilesWithSuffix(dir, "*"+consts.TRIGGERS_FILE_SUFFIX)

	if len(files) > 0 {
		triggersByFilename := make(map[string]triggers.Trigger)
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

			trigger, err := BuildTrigger(triggerData, mimirProcessor)
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

func BuildTrigger(t db.Trigger, mimirProcessor *mimir.MimirProcessor) (triggers.Trigger, error) {

	triggerType, exists := triggerTypeByName[t.Type]
	if !exists {
		return nil, fmt.Errorf("trigger type is missing")
	}
	trigger, err := mimir.TriggerFactory.BuildTrigger(models.TriggerOptions{
		Name:        t.Name,
		TriggerType: triggerType,
		Timeout:     time.Duration(t.Timeout) * time.Second,
		Frequency:   time.Duration(t.Frequency) * time.Second,
	})

	if err != nil {
		return nil, err
	}
	trigger.SetID(t.ID.Hex())
	trigger.Activate()
	err = trigger.UpdateCondition(string(t.Condition))
	if err != nil {
		return nil, err
	}
	BuildActions(t, trigger, mimirProcessor)

	return trigger, nil
}

func BuildActions(triggerData db.Trigger, trigger triggers.Trigger, mimirProcessor *mimir.MimirProcessor) {
	for _, action := range triggerData.Actions {
		triggerAction := ToTriggerAction(action)
		trigger.AddAction(triggerAction)
	}
}

func ToTriggerAction(a db.Action) triggers.Action {
	var triggerAction triggers.Action

	switch a.Type {
	case "print":
		action := triggers.NewPrintAction()
		action.Message = a.Message
		action.Name = a.Name
		triggerAction = action
	case "alert":
		action := mimir.ActionFactory.NewSendMQTTMessageAction(a.Message)
		action.Message = a.Message
		triggerAction = &action
	case "webSocket":
		action := mimir.ActionFactory.NewSendWebSocketMessageAction(a.Message)
		action.Message = a.Message
		triggerAction = &action
	default:
		//TODO see if returning nil is fine or we need some error here
		slog.Warn("action type not recognized while creating trigger action", "type", a.Type)
	}

	return triggerAction
}

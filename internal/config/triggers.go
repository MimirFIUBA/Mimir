package config

import (
	"encoding/json"
	"fmt"
	"log"
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
	for _, v := range files {
		byteValue, err := os.ReadFile(v)
		if err != nil {
			log.Fatal(err)
			return fmt.Errorf("error reading trigger file %s", v)
		}
		var triggerMap map[string]interface{}
		json.Unmarshal(byteValue, &triggerMap)
		trigger := BuildTriggerFromMap(triggerMap, mimirProcessor)

		registerTrigger(trigger, triggerMap)
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

func registerTrigger(trigger *triggers.Trigger, triggerMap map[string]interface{}) {
	topicsInterface, exists := triggerMap["topics"]
	if exists {
		topics, ok := topicsInterface.([]interface{})
		if !ok {
			panic("Bad format for topics")
		}

		for _, topicInterface := range topics {
			topic, ok := topicInterface.(string)
			if !ok {
				panic("Topic is not a string")
			}
			sensor, err := db.SensorsData.GetSensorByTopic(topic)
			if err == nil {
				sensor.Register(trigger)
			}
		}

	}
}

func buildCondition(triggerMap map[string]interface{}) (triggers.Condition, bool) {
	conditionConfiguration, exists := triggerMap["condition"]
	if exists {
		var condition triggers.Condition
		switch conditionValue := conditionConfiguration.(type) {
		case string:
			condition = buildConditionFromString(conditionValue)
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

func buildConditionFromString(conditionString string) triggers.Condition {
	if conditionString != "" {
		tokens := Tokenize(conditionString)
		condition, err := ParseCondition(tokens)
		if err != nil {
			fmt.Println(err)
		} else {
			return condition
		}
	}
	return &triggers.TrueCondition{}
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

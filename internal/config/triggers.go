package config

import (
	"fmt"
	"mimir/internal/db"
	"mimir/internal/mimir"
	"mimir/triggers"

	"github.com/gookit/config"
)

func BuildTriggers(mimirProcessor *mimir.MimirProcessor) {
	configuration := config.Data()
	triggersConfiguration, ok := configuration["triggers"].([]interface{})
	if !ok {
		panic("Wrong processors format, no triggers")
	}

	for _, triggerInterface := range triggersConfiguration {
		triggerMap, ok := triggerInterface.(map[string]interface{})
		if !ok {
			panic("bad configuration")
		}

		trigger := buildTrigger(triggerMap)

		fmt.Println("build condition for ", trigger.Name)
		condition, exists := buildCondition(triggerMap)
		if exists {
			trigger.Condition = condition
		}

		actions := buildActions(triggerMap, mimirProcessor)
		for _, action := range actions {
			trigger.AddAction(action)
		}

		registerTrigger(trigger, triggerMap)
	}
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
	//TODO: implement
	return &triggers.TrueCondition{}
}

func buildConditionFromString(conditionString string) triggers.Condition {
	//TODO: implement
	if conditionString != "" {
		tokens := Tokenize(conditionString)
		condition, err := ParseCondition(tokens)
		if err != nil {
			fmt.Println("err")
			fmt.Println(err)
		} else {
			fmt.Println("tokens", tokens)
			fmt.Println("condition", condition)
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

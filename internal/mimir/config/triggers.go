package config

import (
	"fmt"
	"mimir/internal/db"
	"mimir/triggers"

	"github.com/gookit/config"
)

func BuildTriggers() {
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

		fmt.Println("Building conditions for trigger: ", trigger.Name)
		condition, exists := buildCondition(triggerMap)
		if exists {
			trigger.Condition = condition
		}

		fmt.Println("Building actions for trigger: ", trigger.Name)
		actions := buildActions(triggerMap)
		for _, action := range actions {
			trigger.AddAction(action)
		}

		fmt.Println("Registering trigger: ", trigger.Name)
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

func buildActions(triggerMap map[string]interface{}) []triggers.Action {
	actionsInterface, exists := triggerMap["actions"]
	var actions []triggers.Action
	if exists {
		switch actionsValue := actionsInterface.(type) {
		case []interface{}:
			for _, actionInterface := range actionsValue {
				actionMap, ok := actionInterface.(map[string]interface{})
				if ok {
					action := buildAction(actionMap)
					actions = append(actions, action)
				} else {
					fmt.Println("error parsing action")
				}
			}
		default:
			panic("Error parsing configuration for actions")
		}
	}
	return actions
}

func buildAction(actionMap map[string]interface{}) triggers.Action {
	actionType, exists := actionMap["type"]
	if !exists {
		panic("No type for action")
	}

	var action triggers.Action

	switch actionType {
	case "print":
		action = buildPrintAction(actionMap)
	default:
		fmt.Println("Action type not recognized")
	}

	return action
}

func buildPrintAction(actionMap map[string]interface{}) triggers.Action {
	action := triggers.NewPrintAction()
	messageInterface, exists := actionMap["message"]
	if !exists {
		fmt.Println("no message")
		//TODO should try to look message fmt function
	}

	message, ok := messageInterface.(string)
	if !ok {
		panic("Message is not a string")
	}

	action.SetMessage(message)
	return action
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

func buildConditionFromString(_ string) triggers.Condition {
	//TODO: implement
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

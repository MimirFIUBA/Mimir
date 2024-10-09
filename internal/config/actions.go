package config

import (
	"fmt"
	"mimir/internal/mimir"
	"mimir/triggers"
)

func buildActions(triggerMap map[string]interface{}, mimirProcessor *mimir.MimirProcessor) []triggers.Action {
	actionsInterface, exists := triggerMap["actions"]
	var actions []triggers.Action
	if exists {
		switch actionsValue := actionsInterface.(type) {
		case []interface{}:
			for _, actionInterface := range actionsValue {
				actionMap, ok := actionInterface.(map[string]interface{})
				if ok {
					action := buildAction(actionMap, mimirProcessor)
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

func buildAction(actionMap map[string]interface{}, mimirProcessor *mimir.MimirProcessor) triggers.Action {
	actionType, exists := actionMap["type"]
	if !exists {
		panic("No type for action")
	}

	var action triggers.Action

	switch actionType {
	case "print":
		action = buildPrintAction(actionMap)
	case "alert":
		action = buildAlertAction(actionMap, mimirProcessor.OutgoingMessagesChannel)
	case "webSocket":
		action = buildAlertAction(actionMap, mimirProcessor.WsChannel)
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

func buildAlertAction(actionMap map[string]interface{}, channel chan string) triggers.Action {
	action := triggers.NewSendMessageThroughChannel(channel)
	messageInterface, exists := actionMap["message"]
	if !exists {
		fmt.Println("no message")
		//TODO should try to look message fmt function
	}

	message, ok := messageInterface.(string)
	if !ok {
		panic("Message is not a string")
	}

	action.Message = message
	return action
}

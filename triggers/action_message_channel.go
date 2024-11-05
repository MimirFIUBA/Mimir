package triggers

import (
	"time"

	"github.com/google/uuid"
)

type SendMessageThroughChannel struct {
	MessageContructor       func(Event) string
	Message                 string
	OutgoingMessagesChannel chan string
	NextAction              Action
}

func NewSendMessageThroughChannel(channel chan string) *SendMessageThroughChannel {
	return &SendMessageThroughChannel{}
}

func (action *SendMessageThroughChannel) Execute(event Event) {
	if action.MessageContructor != nil {
		action.OutgoingMessagesChannel <- action.MessageContructor(event)
	} else {
		action.OutgoingMessagesChannel <- action.Message
	}

	if action.NextAction != nil {
		nextEvent := Event{
			Id:        uuid.NewString(),
			Timestamp: time.Now(),
			Data:      map[string]interface{}{"previousEvent": event},
			SenderId:  event.SenderId,
			Type:      CHANNEL_MESSAGE_SENT,
		}

		action.NextAction.Execute(nextEvent)
	}
}

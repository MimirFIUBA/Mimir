package triggers

import (
	"time"

	"github.com/google/uuid"
)

type SendMessageThroughChannel[T any] struct {
	MessageContructor       func(Event) T
	Message                 T
	OutgoingMessagesChannel chan T
	NextAction              Action
}

func NewSendMessageThroughChannel[T any](channel chan T) *SendMessageThroughChannel[T] {
	return &SendMessageThroughChannel[T]{
		OutgoingMessagesChannel: channel,
	}
}

func (action *SendMessageThroughChannel[T]) Execute(event Event) {
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

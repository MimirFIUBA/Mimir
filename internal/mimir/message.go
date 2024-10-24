package mimir

type Message struct {
	Topic   string
	Payload []byte
}

type MessageChannel chan Message

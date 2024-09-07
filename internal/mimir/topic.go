package mimir

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Topic struct {
	Name         string
	IsSubscribed bool
}

type TopicManager struct {
	Topics          map[string]Topic
	MQTTClient      mqtt.Client
	newTopicChannel chan string
}

func NewTopicManager(mqttClient mqtt.Client, topicChannel chan string) *TopicManager {
	return &TopicManager{make(map[string]Topic), mqttClient, topicChannel}
}

func (tm *TopicManager) AddTopic(name string) {
	topic, ok := tm.Topics[name]
	if !ok || !topic.IsSubscribed {
		topic.Name = name
		if token := tm.MQTTClient.Subscribe(topic.Name, 0, onMessageReceived); token.Wait() && token.Error() != nil {
			panic(fmt.Sprintf("Error subscribing to topic: %s", token.Error()))
		}
		tm.Topics[name] = topic
	}
}

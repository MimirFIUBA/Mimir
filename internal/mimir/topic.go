package mimir

import (
	"fmt"
)

type Topic struct {
	Name         string
	IsSubscribed bool
}

func (m *MQTTManager) AddTopic(name string) {
	topic, ok := m.Topics[name]
	if !ok || !topic.IsSubscribed {
		topic.Name = name
		if token := m.MQTTClient.Subscribe(topic.Name, 0, onMessageReceived); token.Wait() && token.Error() != nil {
			panic(fmt.Sprintf("Error subscribing to topic: %s", token.Error()))
		}
		m.Topics[name] = topic
	}
}

func (m *MQTTManager) GetSubscribedTopics() []string {
	var topics []string
	for topicName, topic := range m.Topics {
		if topic.IsSubscribed {
			topics = append(topics, topicName)
		}
	}
	return topics
}

package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand/v2"
	"mimir/internal/config"
	"mimir/internal/consts"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"github.com/gookit/ini/v2"
)

type mqttGenerator struct {
	topic  string
	client mqtt.Client
	c      chan int
}

func (g mqttGenerator) GenerateIntData(n int, id int) {
	for i := 1; i <= n; i++ {
		message := fmt.Sprintf(`{"sensorId": %d, "data": %d, "time": "%s"}`, id, rand.IntN(100), time.Now())
		token := g.client.Publish(g.topic, 0, false, message)
		token.Wait()

		fmt.Printf("Published topic %s: %s\n", g.topic, message)
		time.Sleep(1 * time.Second)
	}

	g.c <- 0
}

func (g mqttGenerator) GenerateFloatData(n int, id string, mps int, multiplier float64, offset float64) {
	for i := 1; i <= n; i++ {
		message := fmt.Sprintf(`{"id": "%s", "data": %.2f, "time": "%s"}`, id, rand.Float64()*multiplier+offset, time.Now())
		token := g.client.Publish(g.topic, 0, false, message)
		token.Wait()

		fmt.Printf("Published topic %s: %s\n", g.topic, message)
		time.Sleep(time.Duration(1000/mps) * time.Millisecond)
	}

	g.c <- 0
}

func (g mqttGenerator) GenerateBytes(id string, numbers ...uint8) {
	buf := new(bytes.Buffer)
	for _, n := range numbers {
		err := binary.Write(buf, binary.BigEndian, n)
		if err != nil {
			log.Fatalf("Failed to encode int: %v", err)
		}
	}

	payload := buf.Bytes()

	token := g.client.Publish(g.topic, 0, false, payload)
	token.Wait()
	if token.Error() != nil {
		log.Printf("Failed to publish message: %v\n", token.Error())
	}

	fmt.Printf("Published topic %s: %08b\n", g.topic, payload)

	g.c <- 0

}

func onMessageReceived(client mqtt.Client, message mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", message.Payload(), message.Topic())

	elapsedTime := time.Since(LastMessageTime)
	fmt.Println("Elapsed time:", elapsedTime)
}

var LastMessageTime time.Time

func mainb() {
	config.LoadIni()
	opts := mqtt.NewClientOptions()
	opts.AddBroker(ini.String(consts.MQTT_BROKER_CONFIG_NAME))

	c := make(chan string)
	fmt.Println("creating new client")
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(fmt.Sprintf("Error connecting to MQTT broker: %s", token.Error()))
	}
	fmt.Println("client connected")

	if token := client.Subscribe("mimir/benchmark-alert", 0, onMessageReceived); token.Wait() && token.Error() != nil {
		panic(fmt.Sprintf("Error subscribing to topic: %s", token.Error()))
	}
	fmt.Println("subscribed to topic")

	multiplier := 50.0
	offset := 50.0
	message := fmt.Sprintf(`{"value": %.2f, "id": "%s"}`, rand.Float64()*multiplier+offset, uuid.New())
	token := client.Publish("mimir/benchmark-test", 0, false, message)
	token.Wait()
	LastMessageTime = time.Now()
	fmt.Println("message sent")
	<-c

	client.Disconnect(250)
}

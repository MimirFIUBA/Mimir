package main

import (
	"fmt"
	"math/rand/v2"
	"mimir/src/consts"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type mqttGenerator struct {
	topic  string
	client mqtt.Client
	c      chan int
}

func (g mqttGenerator) generateIntData(n int, id int) {
	for i := 1; i <= n; i++ {
		message := fmt.Sprintf(`{"sensorId": %d, "data": %d, "time": "%s"}`, id, rand.IntN(100), time.Now())
		token := g.client.Publish(g.topic, 0, false, message)
		token.Wait()

		fmt.Printf("Published topic %s: %s\n", g.topic, message)
		time.Sleep(1 * time.Second)
	}

	g.c <- 0
}

func (g mqttGenerator) generateFloatData(n int, id int) {
	for i := 1; i <= n; i++ {
		message := fmt.Sprintf(`{"sensorId": %d, "data": %.2f, "time": "%s"}`, id, rand.Float64()*40, time.Now())
		token := g.client.Publish(g.topic, 0, false, message)
		token.Wait()

		fmt.Printf("Published topic %s: %s\n", g.topic, message)
		time.Sleep(1 * time.Second)
	}

	g.c <- 0
}

func main() {

	opts := mqtt.NewClientOptions()
	opts.AddBroker(consts.Broker)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(fmt.Sprintf("Error connecting to MQTT broker:", token.Error()))
	}

	c := make(chan int)

	generatorTemp := mqttGenerator{consts.TopicTemp, client, c}
	generatorPH := mqttGenerator{consts.TopicPH, client, c}

	go generatorTemp.generateIntData(10, 2)
	go generatorPH.generateFloatData(10, 1)
	_, _ = <-c, <-c

	client.Disconnect(250)
}

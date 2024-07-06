package main

import (
	"fmt"
	"mimir/src/consts"
	"time"
	"math/rand"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func generateData(n int, client mqtt.Client, topic string, c chan int) {
	for i := 1; i <= n; i++ {
		message := fmt.Sprintf("{data: %.2f, time: %s}", rand.Float64()*40, time.Now())
		token := client.Publish(topic, 0, false, message)
		token.Wait()

		fmt.Println("Published topic %s: %s", topic, message)
		time.Sleep(1 * time.Second)
	}

	c <- 1
}

func main() {

	opts := mqtt.NewClientOptions()
	opts.AddBroker(consts.Broker)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(fmt.Sprintf("Error connecting to MQTT broker:", token.Error()))
	}
	
	c := make(chan int)
	go generateData(10, client, consts.TopicTemp, c)
	go generateData(10, client, consts.TopicPH, c)
	temp, ph := <- c, <- c
	fmt.Println(temp)
	fmt.Println(ph)
	
	client.Disconnect(250)
}
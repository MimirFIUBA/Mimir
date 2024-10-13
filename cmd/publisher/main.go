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
	"github.com/gookit/ini/v2"
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

func (g mqttGenerator) generateFloatData(n int, id string, mps int) {
	for i := 1; i <= n; i++ {
		message := fmt.Sprintf(`{"id": "%s", "data": %.2f, "time": "%s"}`, id, rand.Float64()*40+50, time.Now())
		token := g.client.Publish(g.topic, 0, false, message)
		token.Wait()

		fmt.Printf("Published topic %s: %s\n", g.topic, message)
		time.Sleep(time.Duration(1000/mps) * time.Millisecond)
	}

	g.c <- 0
}

func (g mqttGenerator) generateBytes(id string, numbers []uint8) {
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

func main() {
	config.LoadIni()
	opts := mqtt.NewClientOptions()
	opts.AddBroker(ini.String(consts.MQTT_BROKER_CONFIG_NAME))

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(fmt.Sprintf("Error connecting to MQTT broker: %s", token.Error()))
	}

	c := make(chan int)

	generator := mqttGenerator{"mimir/test", client, c}
	go generator.generateFloatData(20, "1", 2)

	// message := fmt.Sprintf(`{"id": "%s", "data": %.2f, "time": "%s"}`, "0", 15.0, time.Now())
	// token := client.Publish("mimir/esp32/waterTemp", 0, false, message)
	// token.Wait()

	// message = fmt.Sprintf(`{"id": "%s", "data": %.2f, "time": "%s"}`, "0", 9.0, time.Now())
	// token = client.Publish("mimir/esp32/waterTemp", 0, false, message)
	// token.Wait()

	// message = fmt.Sprintf(`{"id": "%s", "data": %.2f, "time": "%s"}`, "0", 11.0, time.Now())
	// token = client.Publish("mimir/esp32/waterTemp", 0, false, message)
	// token.Wait()

	// message = fmt.Sprintf(`{"id": "%s", "data": %.2f, "time": "%s"}`, "0", 55.0, time.Now())
	// token = client.Publish("mimir/esp32/waterTemp", 0, false, message)
	// token.Wait()

	// numbers := []uint8{65, 1, 50, 65, 35, 51, 51}
	// go generator.generateBytes("1", numbers)
	<-c
	// generatorTemp := mqttGenerator{consts.TopicTemp, client, c}
	// generatorPH := mqttGenerator{consts.TopicPH, client, c}

	// go generatorTemp.generateIntData(10, 2)
	// go generatorPH.generateFloatData(10, 1)
	// _, _ = <-c, <-c

	client.Disconnect(250)
}

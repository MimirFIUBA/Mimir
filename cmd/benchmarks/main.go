package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"mimir/internal/config"
	"mimir/internal/consts"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"github.com/gookit/ini/v2"
)

// Estructura del mensaje de envío
type SendMessage struct {
	Value float64 `json:"value"`
	ID    string  `json:"id"`
}

// Estructura del mensaje de respuesta
type ResponseMessage struct {
	Message           string `json:"message"`
	OriginalMessageID string `json:"originalMessageId"`
}

var (
	// Almacena el tiempo de envío de cada mensaje por ID
	sendTimestamps sync.Map
	responseTimes  []time.Duration
	c              = make(chan struct{})
)

func main() {
	// Configuración del cliente MQTT
	config.LoadIni()
	opts := mqtt.NewClientOptions()
	opts.AddBroker(ini.String(consts.MQTT_BROKER_CONFIG_NAME))
	opts.SetClientID("mqtt-benchmark-client")
	opts.OnConnect = func(c mqtt.Client) {
		log.Println("Conectado a MQTT broker")
	}
	opts.OnConnectionLost = func(c mqtt.Client, err error) {
		log.Printf("Conexión perdida: %v", err)
	}

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Error conectando al broker: %v", token.Error())
	}
	defer client.Disconnect(250)

	// Suscribirse al topic de respuestas
	responseTopic := "mimir/benchmark-alert"
	client.Subscribe(responseTopic, 0, func(client mqtt.Client, msg mqtt.Message) {
		handleResponse(msg)
	})

	// Enviar mensajes y calcular el tiempo de respuesta
	messageAmmount := 5
	sendMessages(client, "mimir/benchmark-test", messageAmmount)

	for range messageAmmount {
		<-c
	}

	fmt.Println("Done")

}

// Envía varios mensajes al tema especificado
func sendMessages(client mqtt.Client, topic string, count int) {
	for i := 0; i < count; i++ {
		// Genera un ID único para cada mensaje
		id := uuid.New().String()
		message := SendMessage{
			Value: rand.Float64() * 100, // Valor aleatorio
			ID:    id,
		}

		// Serializa el mensaje a JSON
		payload, err := json.Marshal(message)
		if err != nil {
			log.Printf("Error serializando mensaje: %v", err)
			continue
		}

		// Guarda el tiempo de envío
		sendTimestamps.Store(id, time.Now())

		// Publica el mensaje
		token := client.Publish(topic, 0, false, payload)
		token.Wait()

		log.Printf("Mensaje enviado: %s", payload)
		time.Sleep(1 * time.Second) // Intervalo entre envíos (opcional)
	}
}

// Maneja los mensajes de respuesta recibidos y calcula el tiempo de respuesta
func handleResponse(msg mqtt.Message) {
	var response ResponseMessage
	if err := json.Unmarshal(msg.Payload(), &response); err != nil {
		log.Printf("Error al parsear el mensaje de respuesta: %v", err)
		return
	}

	// Obtiene el timestamp de envío correspondiente al mensaje original
	if startTime, ok := sendTimestamps.Load(response.OriginalMessageID); ok {
		sendTime := startTime.(time.Time)
		elapsedTime := time.Since(sendTime)
		responseTimes = append(responseTimes, elapsedTime)
		log.Printf("Tiempo de respuesta para ID %s: %d", response.OriginalMessageID, elapsedTime)
		fmt.Println("Elapsed time:", elapsedTime)

		// Remueve el ID del mapa una vez calculado
		sendTimestamps.Delete(response.OriginalMessageID)
	} else {
		log.Printf("ID de mensaje %s no encontrado en sendTimestamps", response.OriginalMessageID)
	}

	c <- struct{}{}
}

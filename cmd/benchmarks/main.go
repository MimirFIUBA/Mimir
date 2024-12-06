package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math"
	"math/rand"
	"mimir/internal/config"
	"mimir/internal/consts"
	"os"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"github.com/gookit/ini/v2"

	// chart "github.com/wcarczuk/go-chart/v2"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

var (
	messageCountFlag int
	sleepTime        int
	filename         string
	regression       bool
)

var (
	// Almacena el tiempo de envío de cada mensaje por ID
	sendTimestamps sync.Map
	responseTimes  []time.Duration
	c              = make(chan struct{})
)

// Estructura del mensaje de envío
type SendMessage struct {
	Value float64 `json:"data"`
	ID    string  `json:"id"`
}

// Estructura del mensaje de respuesta
type ResponseMessage struct {
	Message           string `json:"message"`
	OriginalMessageID string `json:"originalMessageId"`
}

func main() {
	flag.IntVar(&messageCountFlag, "count", 1, "ammount of messages that will be send")
	flag.IntVar(&sleepTime, "time", 1000, "time between messages sent in miliseconds")
	flag.StringVar(&filename, "o", "chart.png", "name of output filename for chart")
	flag.BoolVar(&regression, "r", false, "graph linear regression")
	flag.Parse()

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

	fmt.Printf("Sending %d messages\n", messageCountFlag)

	// Enviar mensajes y calcular el tiempo de respuesta
	go sendMessages(client, "mimir/benchmark-test", messageCountFlag)

	for range messageCountFlag {
		<-c
	}

	calculateStats(responseTimes)
	plotDurations(responseTimes, filename)
	fmt.Println("Done")

}

// Envía varios mensajes al tema especificado
func sendMessages(client mqtt.Client, topic string, count int) {
	for i := 0; i < count; i++ {
		// Genera un ID único para cada mensaje
		id := uuid.New().String()
		message := SendMessage{
			Value: rand.Float64()*100 + 50, // Valor aleatorio
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
		time.Sleep(time.Duration(sleepTime) * time.Millisecond) // Intervalo entre envíos (opcional)
	}
}

// Maneja los mensajes de respuesta recibidos y calcula el tiempo de respuesta
func handleResponse(msg mqtt.Message) {
	now := time.Now()
	var response ResponseMessage
	if err := json.Unmarshal(msg.Payload(), &response); err != nil {
		log.Printf("Error al parsear el mensaje de respuesta: %v", err)
		return
	}

	// Obtiene el timestamp de envío correspondiente al mensaje original
	if startTime, ok := sendTimestamps.Load(response.OriginalMessageID); ok {
		sendTime := startTime.(time.Time)
		elapsedTime := now.Sub(sendTime)
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

func calculateStats(responseTimes []time.Duration) {
	var sum time.Duration = 0
	var max time.Duration = responseTimes[0]
	var min time.Duration = responseTimes[0]
	for _, responseTime := range responseTimes {
		sum += responseTime
		if responseTime > max {
			max = responseTime
		}
		if responseTime < min {
			min = responseTime
		}
	}

	avg := sum / time.Duration(len(responseTimes))
	fmt.Println("Average response time: ", avg)
	fmt.Println("Max response time: ", max)
	fmt.Println("Min response time: ", min)
}

// Función que convierte un slice de time.Duration a segundos (float64)
func convertDurationsToFloat(durations []time.Duration) []float64 {
	var result []float64
	for _, duration := range durations {
		// Convertir cada duración a segundos
		result = append(result, float64(duration.Milliseconds()))
	}
	return result
}

// Función para generar las etiquetas dinámicamente, asegurando que haya 10 etiquetas
func generateDynamicLabels(durations []time.Duration) []string {
	var labels []string
	for i := range durations {
		labels = append(labels, fmt.Sprintf("%d", i))
	}
	return labels
}

func plotDurations(durations []time.Duration, filename string) {

	graph := charts.NewLine()

	miliseconds := convertDurationsToFloat(durations)

	labels := generateDynamicLabels(durations)

	graph.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: fmt.Sprintf("Latencia(#%d-%dms)", messageCountFlag, sleepTime),
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Name: "# Mensaje",
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Name: "Latencia (ms)",
			Min:  0,
		}),
	)

	var total float64
	for _, d := range durations {
		total += float64(d.Milliseconds())
	}

	// Calcular el promedio.
	average := total / float64(len(durations))

	graph.AddSeries("Duración", generateData(miliseconds), charts.WithLabelOpts(opts.Label{Show: opts.Bool(true), Position: "bottom"}))
	graph.AddSeries("Promedio", generateAvgData(miliseconds, average)).
		SetSeriesOptions(
			charts.WithLineChartOpts(opts.LineChart{
				ShowSymbol: opts.Bool(false),
			}),
			charts.WithLabelOpts(opts.Label{
				Show: opts.Bool(true),
			}),
		)

	if regression {
		graph.AddSeries("Regresión\nlinear", generateLinearRegressionData(miliseconds)).
			SetSeriesOptions(
				charts.WithLineChartOpts(opts.LineChart{
					ShowSymbol: opts.Bool(false),
				}),
				charts.WithLabelOpts(opts.Label{
					Show: opts.Bool(true),
				}),
			)

		graph.AddSeries("Regresión\nlogarítmica", generateLogarithmicRegressionData(miliseconds)).
			SetSeriesOptions(
				charts.WithLineChartOpts(opts.LineChart{
					ShowSymbol: opts.Bool(false),
				}),
				charts.WithLabelOpts(opts.Label{
					Show: opts.Bool(true),
				}),
			)
	}

	graph.SetXAxis(labels)

	f, err := os.Create(filename + ".html")
	if err != nil {
		fmt.Println("Error creando archivo:", err)
		return
	}
	defer f.Close()

	if err := graph.Render(f); err != nil {
		fmt.Println("Error renderizando gráfico:", err)
	}
}

// Convierte las duraciones en un formato adecuado para las barras
func generateData(durations []float64) []opts.LineData {
	var data []opts.LineData
	for _, d := range durations {
		data = append(data, opts.LineData{Value: d})
	}
	return data
}

func generateLinearRegressionData(durations []float64) []opts.LineData {
	var data []opts.LineData
	var xData []float64
	for i := range durations {
		xData = append(xData, float64(i+1))
	}

	a, b, err := linearRegression(xData, durations)
	if err != nil {
		return data
	}

	fmt.Printf("f(x): %.2fx + %.2f\n", a, b)

	for _, x := range xData {
		data = append(data, opts.LineData{Value: a*x + b})
	}
	return data
}

func generateLogarithmicRegressionData(durations []float64) []opts.LineData {
	var data []opts.LineData
	var xData []float64
	for i := range durations {
		xData = append(xData, float64(i+1))
	}

	a, b, err := logarithmicRegression(xData, durations)
	if err != nil {
		fmt.Println(err)
		return data
	}

	fmt.Printf("f(x): %.2f + %.2fln(x)\n", a, b)

	for _, x := range xData {
		data = append(data, opts.LineData{Value: a + b*math.Log(x)})
	}
	return data
}

func generateAvgData(durations []float64, avg float64) []opts.LineData {
	var data []opts.LineData
	for range durations {
		data = append(data, opts.LineData{Value: avg})
	}
	return data
}

// Función para calcular la regresión lineal
func linearRegression(x, y []float64) (float64, float64, error) {
	// Verificar que ambas listas tengan la misma longitud
	if len(x) != len(y) {
		return 0, 0, fmt.Errorf("las longitudes de x e y no coinciden")
	}

	n := float64(len(x))
	if n == 0 {
		return 0, 0, fmt.Errorf("los datos están vacíos")
	}

	// Calcular sumas necesarias
	var sumX, sumY, sumXY, sumX2 float64
	for i := 0; i < len(x); i++ {
		sumX += x[i]
		sumY += y[i]
		sumXY += x[i] * y[i]
		sumX2 += x[i] * x[i]
	}

	// Calcular pendiente (m) y ordenada al origen (b)
	denominator := n*sumX2 - sumX*sumX
	if denominator == 0 {
		return 0, 0, fmt.Errorf("los datos no permiten calcular una regresión válida")
	}
	m := (n*sumXY - sumX*sumY) / denominator
	b := (sumY*sumX2 - sumX*sumXY) / denominator

	return m, b, nil
}

func logarithmicRegression(x, y []float64) (float64, float64, error) {
	// Verificar que ambas listas tengan la misma longitud
	if len(x) != len(y) {
		return 0, 0, fmt.Errorf("las longitudes de x e y no coinciden")
	}

	n := float64(len(x))
	if n == 0 {
		return 0, 0, fmt.Errorf("los datos están vacíos")
	}

	var sumLogX, sumY, sumLogX2, sumLogXY float64

	for i := 0; i < len(x); i++ {
		// Verificar que x[i] > 0 (para calcular logaritmo)
		if x[i] <= 0 {
			return 0, 0, fmt.Errorf("todos los valores de x deben ser mayores a 0")
		}

		logX := math.Log(x[i]) // Logaritmo natural de x[i]
		sumLogX += logX
		sumY += y[i]
		sumLogX2 += logX * logX
		sumLogXY += logX * y[i]
	}

	// Calcular los coeficientes a y b
	denominator := n*sumLogX2 - sumLogX*sumLogX
	if denominator == 0 {
		return 0, 0, fmt.Errorf("los datos no permiten calcular una regresión válida")
	}
	b := (n*sumLogXY - sumLogX*sumY) / denominator
	a := (sumY - b*sumLogX) / n

	return a, b, nil
}

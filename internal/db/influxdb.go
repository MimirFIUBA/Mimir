package db

import (
	"context"
	"fmt"
	"log/slog"
	"mimir/internal/consts"
	"sort"
	"strings"
	"sync"
	"time"
)

func processPoints(ctx context.Context, wg *sync.WaitGroup) {
	ticker := time.NewTicker(consts.DB_INFLUX_STORE_INTERVAL)
	for {
		select {
		case <-ticker.C:
			wg.Add(1)
			go func() {
				defer wg.Done()
				savePoints()
			}()
		case <-ctx.Done():
			ticker.Stop()
			slog.Info("canceled processing points", "error", ctx.Err())
			return
		}
	}
}

func savePoints() {
	influxDBClient := Database.getInfluxDBClient()
	dumpBuffer := ReadingsBuffer.Dump()
	if len(dumpBuffer) == 0 {
		return
	}
	if influxDBClient == nil {
		slog.Warn("InfluxDB client is not available.")
		return
	}

	slog.Info("Writing to InfluxDB")
	readingsByTopic := make(map[string][]float64)
	for _, reading := range dumpBuffer {
		if readingValue, ok := reading.Value.(float64); ok {
			readingsByTopic[reading.Topic] = append(readingsByTopic[reading.Topic], readingValue)
		}
	}

	if writeApi := Database.getInfluxWriteApi(); writeApi != nil {
		ctx := context.Background()
		for topic, values := range readingsByTopic {
			if len(values) > 0 {
				splittedTopic := strings.Split(topic, `/`)
				unit := splittedTopic[len(splittedTopic)-1]
				last, mean, _, _, max, min, _ := calculateStatistics(values)
				writeApi.WriteRecord(ctx, fmt.Sprintf("%s,unit=%s avg=%f,last=%f,max=%f,min=%f", topic, unit, mean, last, max, min))
			}
		}
		writeApi.Flush(ctx)
	}
}

func calculateStatistics(values []float64) (last float64, mean float64, median float64, mode float64, max float64, min float64, count int) {
	// Last value
	last = values[len(values)-1]
	count = len(values)

	// Mean
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	mean = sum / float64(len(values))

	// Median
	sorted := make([]float64, len(values))
	copy(sorted, values)
	sort.Float64s(sorted)
	min = sorted[0]
	max = sorted[len(sorted)-1]
	if len(sorted)%2 == 0 {
		median = (sorted[len(sorted)/2-1] + sorted[len(sorted)/2]) / 2
	} else {
		median = sorted[len(sorted)/2]
	}

	// Mode
	counts := make(map[float64]int)
	for _, v := range values {
		counts[v]++
	}
	maxCount := 0
	for k, v := range counts {
		if v > maxCount {
			maxCount = v
			mode = k
		}
	}
	return
}

package db

import (
	"context"
	"log/slog"
	"mimir/internal/consts"
	"strings"
	"sync"
	"time"

	"github.com/InfluxCommunity/influxdb3-go/influxdb3"
	"github.com/InfluxCommunity/influxdb3-go/influxdb3/batching"
)

func processPoints(ctx context.Context, wg *sync.WaitGroup) {
	ticker := time.NewTicker(consts.DB_INFLUX_STORE_INTERVAL)
	go func() {
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
	}()
}

func savePoints() {
	influxDBClient := Database.getInfluxDBClient()
	dumpBuffer := ReadingsBuffer.Dump()
	if len(dumpBuffer) > 0 && influxDBClient != nil {
		b := batching.NewBatcher(batching.WithSize(len(dumpBuffer)))
		for _, reading := range dumpBuffer {
			splittedTopic := strings.Split(reading.Topic, `/`)
			unit := splittedTopic[len(splittedTopic)-1]
			readingValue, ok := reading.Value.(float64)
			if ok {
				p := influxdb3.NewPoint(unit,
					map[string]string{"location": "fede"},
					map[string]any{
						"value": readingValue,
					},
					reading.Time)
				b.Add(p)
			}
		}

		if b.Ready() {
			err := influxDBClient.WritePoints(context.TODO(), b.Emit())
			if err != nil {
				slog.Error("Error writing points to influx db - ", "error", err)
			}
		}
	}
}

package db

import (
	"context"
	"log/slog"
	"mimir/internal/consts"
	"strings"
	"time"

	"github.com/InfluxCommunity/influxdb3-go/influxdb3"
	"github.com/InfluxCommunity/influxdb3-go/influxdb3/batching"
)

func processPoints(ctx context.Context) {
	ticker := time.NewTicker(consts.DB_INFLUX_STORE_INTERVAL)
	for {
		select {
		case <-ticker.C:
			savePoints()
		case <-ctx.Done():
			slog.Info("canceled processing points", "error", ctx.Err())
			return
		}
	}
}

func savePoints() {
	influxDBClient := Database.getInfluxDBClient()
	if len(ReadingsDBBuffer) > 0 && influxDBClient != nil {
		b := batching.NewBatcher(batching.WithSize(len(ReadingsDBBuffer)))
		for _, reading := range ReadingsDBBuffer {
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
			err := influxDBClient.WritePoints(context.Background(), b.Emit())
			if err != nil {
				slog.Error("Error writing points to influx db - ", "error", err)
			}
		}
		//TODO: see what to do if error while writing
		ReadingsDBBuffer = ReadingsDBBuffer[:0]
	}
}

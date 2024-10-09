package db

import (
	"context"
	"fmt"
	mimir "mimir/internal/mimir/models"
	"strings"
	"time"

	"github.com/InfluxCommunity/influxdb3-go/influxdb3"
	"github.com/InfluxCommunity/influxdb3-go/influxdb3/batching"
)

var (
	SensorsData = SensorsManager{
		idCounter: 0,
		sensors:   make([]mimir.Sensor, 0),
	}
	NodesData = NodesManager{
		idCounter: 0,
		nodes:     make([]mimir.Node, 0),
	}
	GroupsData = GroupsManager{
		idCounter: 0,
		groups:    make([]mimir.Group, 0),
	}

	ReadingsDBBuffer = make([]mimir.SensorReading, 0)

	DBClient *influxdb3.Client
)

func Run() {
	go func() {
		for {
			if len(ReadingsDBBuffer) > 0 && DBClient != nil {
				fmt.Println("writing, ", len(ReadingsDBBuffer))
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
					err := DBClient.WritePoints(context.Background(), b.Emit())
					if err != nil {
						panic(err)
					}
				}
				ReadingsDBBuffer = ReadingsDBBuffer[:0]
			}
			time.Sleep(5 * time.Second)
		}
	}()
}

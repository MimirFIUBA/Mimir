package db

import (
	"fmt"
	mimir "mimir/internal/mimir/models"
)

func StoreReading(reading mimir.SensorReading) error {
	fmt.Println("store reading: ", reading.SensorID)
	sensor, err := SensorsData.GetSensorById(reading.SensorID)
	if err != nil {
		return err
	}

	sensor.AddReading(reading)
	return nil
}

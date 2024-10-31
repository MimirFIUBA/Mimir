package db

import (
	"mimir/internal/models"
)

func StoreReading(reading models.SensorReading) error {
	if reading.SensorID != "" {
		sensor, err := SensorsData.GetSensorById(reading.SensorID)
		if err != nil {
			return err
		}
		sensor.AddReading(reading)
	} else {
		sensor, err := SensorsData.GetSensorByTopic(reading.Topic)
		if err != nil {
			return err
		}
		sensor.AddReading(reading)
	}

	ReadingsBuffer.AddReading(reading)

	return nil
}

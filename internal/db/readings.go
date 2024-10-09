package db

import (
	mimir "mimir/internal/mimir/models"
)

func StoreReading(reading mimir.SensorReading) error {
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

	ReadingsDBBuffer = append(ReadingsDBBuffer, reading)

	return nil
}

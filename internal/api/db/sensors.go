package db

import (
	"fmt"
	"mimir/internal/api/models"
	"strconv"
)

type SensorsManager struct {
	sensors   []models.Sensor
	idCounter int
}

func (s *SensorsManager) GetNewId() int {
	s.idCounter++
	return s.idCounter
}

func (s *SensorsManager) GetSensors() []models.Sensor {
	return s.sensors
}

func (s *SensorsManager) GetSensorById(id string) (*models.Sensor, error) {
	for index, sensor := range s.sensors {
		if sensor.ID == id {
			return &s.sensors[index], nil
		}
	}
	return nil, fmt.Errorf("sensor %s not found", id)
}

func (s *SensorsManager) IdExists(id string) bool {
	_, err := s.GetSensorById(id)
	if err != nil {
		return false
	}

	return true
}

func (s *SensorsManager) CreateSensor(sensor *models.Sensor) error {
	newId := s.GetNewId()
	sensor.ID = strconv.Itoa(newId)

	s.sensors = append(s.sensors, *sensor)
	err := NodesData.AddSensorToNodeById(sensor.NodeID, sensor)
	if err != nil {
		return err
	}

	return nil
}

func (s *SensorsManager) UpdateSensor(sensor *models.Sensor) (*models.Sensor, error) {
	oldSensor, err := s.GetSensorById(sensor.ID)
	if err != nil {
		return nil, err
	}

	oldSensor.Update(sensor)
	return sensor, nil
}

func (s *SensorsManager) DeleteSensor(id string) error {
	sensorIndex := -1
	for i := range s.sensors {
		sensor := &s.sensors[i]
		if sensor.ID == id {
			sensorIndex = i
			break
		}
	}

	if sensorIndex == -1 {
		return fmt.Errorf("sensor %s not found", id)
	}

	s.sensors[sensorIndex] = s.sensors[len(s.sensors)-1]
	s.sensors = s.sensors[:len(s.sensors)-1]
	return nil
}

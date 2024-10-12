package db

import (
	"fmt"
	"log"
	mimir "mimir/internal/mimir/models"

	"go.mongodb.org/mongo-driver/bson"
)

type SensorsManager struct {
	sensors   []mimir.Sensor
	idCounter int
}

func (s *SensorsManager) GetNewId() int {
	s.idCounter++
	return s.idCounter
}

func (s *SensorsManager) GetSensors() []mimir.Sensor {
	return s.sensors
}

func (s *SensorsManager) GetSensorById(id string) (*mimir.Sensor, error) {
	for index, sensor := range s.sensors {
		if sensor.ID == id {
			return &s.sensors[index], nil
		}
	}
	return nil, fmt.Errorf("sensor %s not found", id)
}

func (s *SensorsManager) GetSensorByTopic(topic string) (*mimir.Sensor, error) {
	for index, sensor := range s.sensors {
		if sensor.Topic == topic {
			return &s.sensors[index], nil
		}
	}
	return nil, fmt.Errorf("sensor with topic %s not found", topic)
}

func (s *SensorsManager) IdExists(id string) bool {
	_, err := s.GetSensorById(id)
	return err == nil
}

func (s *SensorsManager) CreateSensor(sensor *mimir.Sensor) error {
	// TODO(#20) - Add Body validation

	sensor, err := Database.insertTopic(sensor)
	if err != nil {
		log.Fatal(err)
		return err
	}

	s.sensors = append(s.sensors, *sensor)
	err = NodesData.AddSensorToNodeById(sensor.NodeID, sensor)
	if err != nil {
		return err
	}
	return nil
}

func (s *SensorsManager) UpdateSensor(sensor *mimir.Sensor) (*mimir.Sensor, error) {
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

func (s *SensorsManager) LoadSensors(sensors []*mimir.Sensor) {
	values := bson.A{}
	sensorsMap := make(map[string]*mimir.Sensor)
	for _, sensor := range sensors {
		values = append(values, bson.D{{Key: "name", Value: sensor.Name}})
		sensorsMap[sensor.Name] = sensor
	}

	filter := bson.D{{Key: "$or", Value: values}}

	results, err := Database.findTopics(filter)
	if err != nil {
		log.Fatal(err)
		return
	}

	existingSensorsMap := make(map[string]mimir.Sensor)
	for _, result := range results {
		existingSensorsMap[result.Name] = result
	}

	var sensorsToInsert []interface{}
	for _, sensor := range sensors {
		_, exists := existingSensorsMap[sensor.Name]
		if !exists {
			sensorsToInsert = append(sensorsToInsert, sensor)
		}
	}

	if len(sensorsToInsert) > 0 {
		Database.insertTopics(sensorsToInsert)
	}
}

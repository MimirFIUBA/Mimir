package db

import (
	"fmt"
	"log/slog"
	"mimir/internal/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SensorsManager struct {
	sensors   []*models.Sensor
	idCounter int
}

func (s *SensorsManager) GetNewId() int {
	s.idCounter++
	return s.idCounter
}

func (s *SensorsManager) GetSensors() []*models.Sensor {
	return s.sensors
}

func (s *SensorsManager) GetSensorsMap() map[string]*models.Sensor {
	sensorsMap := make(map[string]*models.Sensor)
	for _, sensor := range s.sensors {
		sensorsMap[sensor.Topic] = sensor
	}
	return sensorsMap
}

func (s *SensorsManager) GetSensorById(id string) (*models.Sensor, error) {
	for index, sensor := range s.sensors {
		if sensor.GetId() == id {
			return s.sensors[index], nil
		}
	}
	return nil, fmt.Errorf("sensor %s not found", id)
}

func (s *SensorsManager) GetSensorByTopic(topic string) (*models.Sensor, error) {
	//TODO: change error for bool
	for index, sensor := range s.sensors {
		if sensor.Topic == topic {
			return s.sensors[index], nil
		}
	}
	return nil, fmt.Errorf("sensor with topic %s not found", topic)
}

func (s *SensorsManager) IdExists(id string) bool {
	_, err := s.GetSensorById(id)
	return err == nil
}

func (s *SensorsManager) CreateSensor(sensor *models.Sensor) (*models.Sensor, error) {
	// TODO(#20) - Add Body validation

	sensor, err := Database.insertTopic(sensor)
	if err != nil {
		slog.Error("error inserting topic", "error", err, "topic", sensor)
		return nil, err
	}

	s.sensors = append(s.sensors, sensor)
	err = NodesData.AddSensorToNodeById(sensor.NodeID, sensor)
	if err != nil {
		slog.Error("error adding sensor to node", "error", err, "topic", sensor)
		return nil, err
	}
	return sensor, nil
}

func (s *SensorsManager) UpdateSensor(sensor *models.Sensor, id string) (*models.Sensor, error) {
	oldSensor, err := s.GetSensorById(id)
	if err != nil {
		return nil, err
	}

	oldSensor.Update(sensor)
	return sensor, nil
}

func (s *SensorsManager) SetSensorsToInactive() {
	for i := range s.sensors {
		s.sensors[i].IsActive = false
	}
}

func (s *SensorsManager) DeleteSensor(id string) error {
	sensorIndex := -1
	for i, sensor := range s.sensors {
		if sensor.GetId() == id {
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

func (s *SensorsManager) LoadSensors(sensors []*models.Sensor) {
	existingSensorsMap := make(map[string]*models.Sensor)
	if len(sensors) > 0 {
		filter := buildTopicFilter(sensors)
		results, err := Database.FindTopics(filter)
		if err != nil {
			slog.Error("fail to find topics", "sensors", sensors)
			return
		}

		for _, result := range results {
			existingSensorsMap[result.Topic] = &result
		}

		var sensorsToInsert []interface{}
		var sensorsToReactivate []*models.Sensor
		for _, sensor := range sensors {
			s.sensors = append(s.sensors, sensor)
			existingSensor, exists := existingSensorsMap[sensor.Topic]
			if !exists {
				sensorsToInsert = append(sensorsToInsert, sensor)
				NodesData.AddSensorToNodeById(sensor.NodeID, sensor)
			} else {
				sensor.ID = existingSensor.ID
				sensorsToReactivate = append(sensorsToReactivate, sensor)
				NodesData.AddSensorToNodeById(sensor.NodeID, sensor)
			}
		}
		insertAndUpdateTopicIds(sensorsToInsert)
		activateTopics(sensorsToReactivate)
	}
}

func insertAndUpdateTopicIds(sensorsToInsert []interface{}) {
	if len(sensorsToInsert) > 0 {
		insertedIds := Database.insertTopics(sensorsToInsert)
		for i, id := range insertedIds {
			objectId, ok := id.(primitive.ObjectID)
			if ok {
				insertedSensor, ok := sensorsToInsert[i].(*models.Sensor)
				if ok {
					insertedSensor.ID = objectId
				}
			}
		}
	}
}

func activateTopics(sensorsToActivate []*models.Sensor) {
	if len(sensorsToActivate) > 0 {
		Database.ActivateTopics(sensorsToActivate)
	}
}

func GetHandlerForTopic(topic string) {

}

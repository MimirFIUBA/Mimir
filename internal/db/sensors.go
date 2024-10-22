package db

import (
	"context"
	"fmt"
	"log"
	"mimir/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		if sensor.GetId() == id {
			return &s.sensors[index], nil
		}
	}
	return nil, fmt.Errorf("sensor %s not found", id)
}

func (s *SensorsManager) GetSensorByTopic(topic string) (*models.Sensor, error) {
	//TODO: change error for bool
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

func (s *SensorsManager) CreateSensor(sensor *models.Sensor) error {
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
	for i := range s.sensors {
		sensor := &s.sensors[i]
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

func buildTopicFilter(sensors []models.Sensor) bson.D {
	values := bson.A{}
	for _, sensor := range sensors {
		values = append(values, sensor.Topic)
	}

	return bson.D{{Key: "topic", Value: bson.D{{Key: "$in", Value: values}}}}
}

func (s *SensorsManager) LoadSensors(sensors []models.Sensor) {
	existingSensorsMap := make(map[string]*models.Sensor)
	if len(sensors) > 0 {
		filter := buildTopicFilter(sensors)
		results, err := Database.findTopics(filter)
		if err != nil {
			log.Fatal(err)
			return
		}

		for _, result := range results {
			existingSensorsMap[result.Topic] = &result
		}
	}

	var sensorsToInsert []interface{}
	for _, sensor := range sensors {
		s.sensors = append(s.sensors, sensor)
		_, exists := existingSensorsMap[sensor.Topic]
		if !exists {
			sensorsToInsert = append(sensorsToInsert, sensor)
		}
	}

	if len(sensorsToInsert) > 0 {
		Database.insertTopics(sensorsToInsert)
	}
}

func (d *DatabaseManager) insertTopic(topic *models.Sensor) (*models.Sensor, error) {
	mongoClient := d.getMongoClient()
	if mongoClient != nil {
		topicsCollection := mongoClient.Database(MONGO_DB_MIMIR).Collection(TOPICS_COLLECTION)
		result, err := topicsCollection.InsertOne(context.TODO(), topic)
		if err != nil {
			fmt.Println("error inserting group ", err)
			return nil, err
		}

		topicId, ok := result.InsertedID.(primitive.ObjectID)
		if !ok {
			return nil, fmt.Errorf("error converting id for group")
		}
		topic.ID = topicId
	}

	return topic, nil
}

func (d *DatabaseManager) insertTopics(topics []interface{}) {
	mongoClient := d.getMongoClient()
	if mongoClient != nil {
		topicsCollection := mongoClient.Database(MONGO_DB_MIMIR).Collection(TOPICS_COLLECTION)
		_, err := topicsCollection.InsertMany(context.TODO(), topics)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (d *DatabaseManager) DeactivateTopics(sensors []models.Sensor) {
	filter := buildTopicFilter(sensors)
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "is_active", Value: false}}}}
	mongoClient := d.getMongoClient()
	if mongoClient != nil {
		topicsCollection := mongoClient.Database(MONGO_DB_MIMIR).Collection(TOPICS_COLLECTION)
		_, err := topicsCollection.UpdateMany(context.TODO(), filter, update)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (d *DatabaseManager) findTopics(filter primitive.D) ([]models.Sensor, error) {
	var results []models.Sensor
	mongoClient := d.getMongoClient()
	if mongoClient != nil {
		topicsCollection := mongoClient.Database(MONGO_DB_MIMIR).Collection(TOPICS_COLLECTION)
		cursor, err := topicsCollection.Find(context.TODO(), filter)
		if err != nil {
			return nil, err
		} else {
			defer cursor.Close(context.TODO())
		}

		if err = cursor.All(context.TODO(), &results); err != nil {
			return nil, err
		}
	}
	return results, nil
}

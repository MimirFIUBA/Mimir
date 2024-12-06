package db

import (
	"mimir/internal/models"

	"go.mongodb.org/mongo-driver/bson"
)

func buildNameFilterForVariables(variables []*UserVariable) bson.D {
	values := bson.A{}
	for _, variable := range variables {
		values = append(values, variable.Name)
	}

	return bson.D{{Key: "name", Value: bson.D{{Key: "$in", Value: values}}}}
}

func buildTopicFilter(sensors []*models.Sensor) bson.D {
	values := bson.A{}
	for _, sensor := range sensors {
		values = append(values, sensor.Topic)
	}

	return bson.D{{Key: "topic", Value: bson.D{{Key: "$in", Value: values}}}}
}

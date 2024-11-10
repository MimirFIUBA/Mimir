package db

import "go.mongodb.org/mongo-driver/bson"

func buildNameFilterForVariables(variables []*UserVariable) bson.D {
	values := bson.A{}
	for _, variable := range variables {
		values = append(values, variable.Name)
	}

	return bson.D{{Key: "name", Value: bson.D{{Key: "$in", Value: values}}}}
}

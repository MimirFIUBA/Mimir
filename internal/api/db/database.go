package db

import "mimir/internal/api/models"

var (
	SensorsData = SensorsManager{
		sensors: []models.Sensor{
			{ID: "1", Name: "Example1", Description: "Example description"},
			{ID: "2", Name: "Example2", Description: "Example description"},
		},
		idCounter: 2,
	}
)

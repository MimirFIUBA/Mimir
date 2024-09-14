package db

import "mimir/internal/api/models"

var (
	SensorsData = SensorsManager{
		idCounter: 0,
		sensors:   make([]models.Sensor, 0),
	}
	NodesData = NodesManager{
		idCounter: 0,
		nodes:     make([]models.Node, 0),
	}
	GroupsData = GroupsManager{
		idCounter: 0,
		groups:    make([]models.Group, 0),
	}
)

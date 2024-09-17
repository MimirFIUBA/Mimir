package db

import mimir "mimir/internal/mimir/models"

var (
	SensorsData = SensorsManager{
		idCounter: 0,
		sensors:   make([]mimir.Sensor, 0),
	}
	NodesData = NodesManager{
		idCounter: 0,
		nodes:     make([]mimir.Node, 0),
	}
	GroupsData = GroupsManager{
		idCounter: 0,
		groups:    make([]mimir.Group, 0),
	}
)

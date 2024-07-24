package mimir

import (
	"fmt"

	"github.com/google/uuid"
)

type DataManager struct {
	groups  []Group
	nodes   []Node
	sensors []Sensor
}

func (d *DataManager) AddGroup(group *Group) *Group {
	if group.ID == uuid.Nil {
		group.ID = uuid.New()
	}
	d.groups = append(d.groups, *group)
	return group
}

func (d *DataManager) GetGroups() []Group {
	return d.groups
}

func (d *DataManager) GetGroup(ID uuid.UUID) *Group {
	for i := range d.groups {
		group := &d.groups[i]
		if group.ID == ID {
			return group
		}
	}
	return nil
}

func (d *DataManager) AddNode(node *Node) *Node {
	if node.ID == uuid.Nil {
		node.ID = uuid.New()
	}

	group := d.GetGroup(node.GroupID)
	if group != nil {
		group.Nodes = append(group.Nodes, *node)
	}

	d.nodes = append(d.nodes, *node)
	return node
}

func (d *DataManager) GetNodes() []Node {
	return d.nodes
}

func (d *DataManager) GetNode(ID uuid.UUID) *Node {
	for i := range d.nodes {
		node := &d.nodes[i]
		if node.ID == ID {
			return node
		}
	}
	return nil
}

func (d *DataManager) getNewSensorId() int {
	return len(d.sensors)
}

func (d *DataManager) storeReading(reading SensorReading) {
	for i := range d.sensors {
		sensor := &d.sensors[i]
		if sensor.ID == reading.SensorID {
			sensor.addReading(reading)
			break
		}
	}
}

func AddSensor(sensor *Sensor) *Sensor {
	sensor.ID = Data.getNewSensorId()
	Data.sensors = append(Data.sensors, *sensor)
	fmt.Printf("New sensor created: %+v\n", sensor)
	return sensor
}

func GetSensors() []Sensor {
	return Data.sensors
}

func GetSensor(id int) *Sensor {
	for _, sensor := range Data.sensors {
		if sensor.ID == id {
			return &sensor
		}
	}
	return nil
}

func StoreReading(reding SensorReading) {
	Data.storeReading(reding)
}

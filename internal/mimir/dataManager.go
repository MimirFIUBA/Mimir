package mimir

import (
	"fmt"

	"github.com/google/uuid"
)

type DataManager struct {
	groups       []Group
	nodes        []Node
	sensors      []Sensor
	topicChannel chan string
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

func (d *DataManager) StoreReading(reading SensorReading) {
	for i := range d.sensors {
		sensor := &d.sensors[i]
		if sensor.ID == reading.SensorID {
			sensor.addReading(reading)
			break
		}
	}
}

func (d *DataManager) AddSensor(sensor *Sensor) *Sensor {
	sensor.ID = Data.getNewSensorId()
	topicName := "topic/"

	node := d.GetNode(sensor.NodeID)
	if node != nil {
		node.Sensors = append(node.Sensors, *sensor)
		topicName += node.Name + "/"
	}
	d.sensors = append(d.sensors, *sensor)

	topicName += sensor.DataName

	fmt.Printf("New topic: %+v\n", topicName)
	// Topics.AddTopic(topicName)
	d.topicChannel <- topicName

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

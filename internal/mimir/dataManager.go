package mimir

import (
	"fmt"
	"slices"
	"strconv"

	"mimir/internal/consts"
	dh "mimir/internal/dataHandler"
	"mimir/internal/triggers"

	"github.com/google/uuid"
)

type DataManager struct {
	groups  []dh.Data
	nodes   []Node
	sensors []Sensor

	topicChannel chan string
}

func (d *DataManager) AddGroup(group *Group) *Group {
	if group.ID == "" {
		group.ID = uuid.New().String()
	}
	d.groups = append(d.groups, group)
	return group
}

func (d *DataManager) GetGroups() []dh.Data {
	return d.groups
}

func (d *DataManager) GetGroup(ID string) *Group {
	for i := range d.groups {
		data := d.groups[i]
		if data.GetId() == ID {
			group, ok := (data).(*Group)
			if ok {
				return group
			}
		}
	}
	return nil
}

func (d *DataManager) UpdateGroup(group *Group) *Group {
	existingGroup := d.GetGroup(group.ID)
	if existingGroup != nil {
		existingGroup.Update(group)
	}
	return existingGroup
}

func (d *DataManager) DeleteGroup(id string) {
	var groupIndex int
	for i := range d.groups {
		group := d.groups[i]
		if group.GetId() == id {
			groupIndex = i
			break
		}
	}

	d.groups[groupIndex] = d.groups[len(d.groups)-1]
	d.groups = d.groups[:len(d.groups)-1]
}

func (d *DataManager) AddNode(node *Node) *Node {
	if node.ID == "" {
		node.ID = uuid.New().String()
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

func (d *DataManager) GetNode(ID string) *Node {
	idx := slices.IndexFunc(d.nodes, func(n Node) bool {
		return n.ID == ID
	})
	if idx >= 0 {
		node := &d.nodes[idx]
		return node
	}
	return nil
}

func (d *DataManager) UpdateNode(node *Node) *Node {
	existingNode := d.GetNode(node.ID)
	if existingNode != nil {
		existingNode.Update(node)
	}
	return existingNode
}

func (d *DataManager) DeleteNode(id string) {
	var nodeIndex int
	for i := range d.nodes {
		node := &d.nodes[i]
		if node.ID == id {
			nodeIndex = i
			break
		}
	}

	d.nodes[nodeIndex] = d.nodes[len(d.nodes)-1]
	d.nodes = d.nodes[:len(d.nodes)-1]
}

func (d *DataManager) getNewSensorId() string {
	return strconv.Itoa(len(d.sensors))
}

func (d *DataManager) StoreReading(reading SensorReading) {
	fmt.Println("Store reading")
	fmt.Printf("reading: %v\n", reading)
	for i := range d.sensors {
		sensor := &d.sensors[i]
		if sensor.GetId() == reading.SensorID {
			fmt.Println("sensor add reading")
			sensor.addReading(reading)
			break
		}
	}
}

func (d *DataManager) GetSensors() []Sensor {
	return d.sensors
}

func (d *DataManager) GetSensor(id string) *Sensor {
	for i := range d.sensors {
		sensor := &d.sensors[i]
		if sensor.ID == id {
			return sensor
		}
	}
	return nil
}

func (d *DataManager) AddSensor(sensor *Sensor) *Sensor {
	fmt.Println("Add sensor")
	sensor.ID = d.getNewSensorId()
	sensor.Topic = consts.TopicPrefix + "/"
	nodeId := sensor.NodeID

	node := d.GetNode(sensor.NodeID)

	if node != nil {
		sensor.Topic += node.Name + "/" + sensor.DataName
		node.Sensors = append(node.Sensors, *sensor)
	} else {
		sensor.Topic += sensor.DataName
	}

	for _, data := range d.groups {
		group := data.(*Group)
		for i := range group.Nodes {
			node := &group.Nodes[i]
			if node.ID == nodeId {
				node.Sensors = append(node.Sensors, *sensor)
			}
		}
	}

	d.sensors = append(d.sensors, *sensor)

	fmt.Printf("New topic: %+v\n", sensor.Topic)

	d.topicChannel <- sensor.Topic

	fmt.Printf("New sensor created: %+v\n", sensor)
	return sensor
}

func (d *DataManager) UpdateSensor(sensor *Sensor) *Sensor {
	existingSensor := d.GetSensor(sensor.ID)
	if existingSensor != nil {
		existingSensor.Update(sensor)
	}
	return existingSensor
}

func (d *DataManager) DeleteSensor(id string) {
	var sensorIndex int
	for i := range d.sensors {
		sensor := &d.sensors[i]
		if sensor.ID == id {
			sensorIndex = i
			break
		}
	}

	d.sensors[sensorIndex] = d.sensors[len(d.sensors)-1]
	d.sensors = d.sensors[:len(d.sensors)-1]
}

func (d *DataManager) GetTriggersBySensorId() map[string][]triggers.TriggerObserver {
	var triggersBySensorId = make(map[string][]triggers.TriggerObserver)
	for _, sensor := range d.sensors {
		if len(sensor.triggerList) > 0 {
			var triggerList = sensor.triggerList
			triggersBySensorId[sensor.ID] = triggerList
		}
	}
	return triggersBySensorId
}

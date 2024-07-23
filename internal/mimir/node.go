package mimir

type Node struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	GroupID     int      `json:"groupId"`
	Sensors     []Sensor `json:"data"`
}

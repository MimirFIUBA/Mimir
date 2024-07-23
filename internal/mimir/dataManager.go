package mimir

type DataManager struct {
	groups []Group
	nodes  []Node
}

func (d *DataManager) AddGroup(group *Group) *Group {
	d.groups = append(d.groups, *group)
	return group
}

func (d *DataManager) GetGroups() []Group {
	return d.groups
}

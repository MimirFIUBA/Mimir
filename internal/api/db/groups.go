package db

import (
	"fmt"
	"mimir/internal/api/models"
	"strconv"
)

type GroupsManager struct {
	groups    []models.Group
	idCounter int
}

func (g *GroupsManager) GetNewId() int {
	g.idCounter++
	return g.idCounter
}

func (g *GroupsManager) GetGroups() []models.Group {
	return g.groups
}

func (g *GroupsManager) GetGroupById(id string) (*models.Group, error) {
	for index, group := range g.groups {
		if group.ID == id {
			return &g.groups[index], nil
		}
	}

	// TODO(#19) - Improve error handling
	return nil, fmt.Errorf("group %s not found", id)
}

func (g *GroupsManager) CreateGroup(group *models.Group) error {
	newId := g.GetNewId()
	group.ID = strconv.Itoa(newId)

	g.groups = append(g.groups, *group)
	// TODO - Add node relationship
	return nil
}

func (g *GroupsManager) UpdateGroup(group *models.Group) (*models.Group, error) {
	oldGroup, err := g.GetGroupById(group.ID)
	// TODO(#19) - Improve error handling
	if err != nil {
		return nil, err
	}

	oldGroup.Update(group)
	return group, nil
}

func (g *GroupsManager) DeleteGroup(id string) error {
	groupIndex := -1
	for i := range g.groups {
		if g.groups[i].ID == id {
			groupIndex = i
			break
		}
	}

	// TODO(#19) - Improve error handling
	if groupIndex == -1 {
		return fmt.Errorf("group %s not found", id)
	}

	g.groups[groupIndex] = g.groups[len(g.groups)-1]
	g.groups = g.groups[:len(g.groups)-1]
	return nil
}

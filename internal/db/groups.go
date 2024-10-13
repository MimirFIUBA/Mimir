package db

import (
	"context"
	"fmt"
	"mimir/internal/mimir/models"
	mimir "mimir/internal/mimir/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GroupsManager struct {
	groups    []mimir.Group
	idCounter int
}

func (g *GroupsManager) GetNewId() int {
	g.idCounter++
	return g.idCounter
}

func (g *GroupsManager) GetGroups() []mimir.Group {
	return g.groups
}

func (g *GroupsManager) GetGroupById(id string) (*mimir.Group, error) {
	for index, group := range g.groups {
		if group.GetId() == id {
			return &g.groups[index], nil
		}
	}

	return nil, fmt.Errorf("group %s not found", id)
}

func (g *GroupsManager) IdExists(id string) bool {
	_, err := g.GetGroupById(id)
	return err == nil
}

func (g *GroupsManager) CreateGroup(group *mimir.Group) error {
	group, err := Database.insertGroup(group)
	if err != nil {
		return err
	}

	g.AddGroup(group)

	return nil
}

func (g *GroupsManager) AddGroup(group *mimir.Group) error {
	g.groups = append(g.groups, *group)

	return nil
}

func (g *GroupsManager) UpdateGroup(group *mimir.Group) (*mimir.Group, error) {
	oldGroup, err := g.GetGroupById(group.GetId())
	if err != nil {
		return nil, err
	}

	oldGroup.Update(group)
	return group, nil
}

func (g *GroupsManager) DeleteGroup(id string) error {
	groupIndex := -1
	for i := range g.groups {
		if g.groups[i].GetId() == id {
			groupIndex = i
			break
		}
	}

	if groupIndex == -1 {
		return fmt.Errorf("group %s not found", id)
	}

	g.groups[groupIndex] = g.groups[len(g.groups)-1]
	g.groups = g.groups[:len(g.groups)-1]
	return nil
}

func (g *GroupsManager) AddNodeToGroupById(id string, node *mimir.Node) error {
	oldGroup, err := g.GetGroupById(id)
	if err != nil {
		return nil
	}

	return oldGroup.AddNode(node)
}

func (d *DatabaseManager) insertGroup(group *models.Group) (*models.Group, error) {
	mongoClient := d.getMongoClient()
	if mongoClient != nil {
		groupsCollection := mongoClient.Database(MONGO_DB_MIMIR).Collection(GROUPS_COLLECTION)
		result, err := groupsCollection.InsertOne(context.TODO(), group)
		if err != nil {
			fmt.Println("error inserting group ", err)
			return nil, err
		}

		groupId, ok := result.InsertedID.(primitive.ObjectID)
		if !ok {
			return nil, fmt.Errorf("error converting id for group")
		}
		group.ID = groupId
	}
	return group, nil
}

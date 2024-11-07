package db

import (
	"context"
	"fmt"
	"log/slog"
	"mimir/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (g *GroupsManager) CreateGroup(group *models.Group) error {
	group, err := Database.InsertGroup(group)
	if err != nil {
		return err
	}

	g.AddGroup(group)

	return nil
}

func (g *GroupsManager) AddGroup(group *models.Group) error {
	g.groups = append(g.groups, *group)

	return nil
}

func (g *GroupsManager) UpdateGroup(group *models.Group) (*models.Group, error) {
	oldGroup, err := g.GetGroupById(group.GetId())
	if err != nil {
		return nil, err
	}

	oldGroup.Update(group)

	_, err = Database.UpdateGroup(oldGroup)
	if err != nil {
		return nil, err
	}

	return oldGroup, nil
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

func (g *GroupsManager) AddNodeToGroupById(id string, node *models.Node) error {
	oldGroup, err := g.GetGroupById(id)
	if err != nil {
		return nil
	}

	err = oldGroup.AddNode(node)
	if err != nil {
		return err
	}

	_, err = Database.UpdateGroup(oldGroup)

	return err
}

func (d *DatabaseManager) InsertGroup(group *models.Group) (*models.Group, error) {
	mongoClient := d.getMongoClient()
	if mongoClient != nil {
		groupsCollection := mongoClient.Database(MONGO_DB_MIMIR).Collection(GROUPS_COLLECTION)
		result, err := groupsCollection.InsertOne(context.TODO(), group)
		if err != nil {
			slog.Error("error inserting group", "error", err)
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

func (d *DatabaseManager) UpdateGroup(group *models.Group) (*models.Group, error) {
	mongoClient := d.getMongoClient()
	if mongoClient != nil {
		groupsCollection := mongoClient.Database(MONGO_DB_MIMIR).Collection(GROUPS_COLLECTION)
		objectId, err := primitive.ObjectIDFromHex(group.GetId())
		if err != nil {
			return nil, err
		}
		update := bson.D{{Key: "$set", Value: group}}
		_, err = groupsCollection.UpdateByID(context.TODO(), objectId, update)
		if err != nil {
			fmt.Println("error updating group ", err)
			return nil, err
		}
	}
	return group, nil
}

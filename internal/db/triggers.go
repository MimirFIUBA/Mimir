package db

import (
	"context"
	"encoding/json"
	"fmt"
	"mimir/internal/consts"
	"mimir/triggers"
	"os"

	"github.com/gookit/ini/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Trigger struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name,omitempty"`
	Filename  string             `json:"filename" bson:"filename,omitempty"`
	IsActive  bool               `json:"active" bson:"active"`
	Topics    []string           `json:"topics" bson:"topics"`
	Condition Condition          `json:"condition" bson:"condition"`
	Actions   []Action           `json:"actions" bson:"actions,omitempty"`
}

type Condition string

type Action struct {
	Name    string `bson:"name"`
	Type    string `bson:"type"`
	Message string `bson:"message,omitempty"`
}

func (a Action) ToTriggerAction() triggers.Action {
	var action triggers.Action

	switch a.Type {
	case "print":
		action := triggers.NewPrintAction()
		action.Message = a.Message
		action.Name = a.Name
	case "alert":
	case "webSocket":
	default:
		fmt.Println("Action type not recognized")
	}

	return action
}

func (t *Trigger) BuildTriggerObserver() triggers.TriggerObserver {
	trigger := triggers.NewTrigger(t.Name)
	// trigger.Condition = config.BuildConditionFromString(string(t.Condition))
	//TODO: build Actions

	return trigger
}

func (t *Trigger) BuildFileName(suffix string) string {
	filename := t.Name
	if suffix != "" {
		filename = t.Name + "_" + suffix
	}
	return ini.String(consts.TRIGGERS_DIR_CONFIG_NAME) + "/" + filename + consts.TRIGGERS_FILE_SUFFIX
}

func AddNewTriggerFromMap(triggerMap map[string]interface{}) {

}

func AddNewTriggerObserver(trigger *triggers.TriggerObserver) (*Trigger, error) {
	return nil, nil
}

func AddNewTrigger(trigger *Trigger) (*Trigger, error) {
	trigger, err := Database.InsertTrigger(trigger)
	return trigger, err

}

func (d *DatabaseManager) GetTriggers() []triggers.TriggerObserver {
	var triggerList []triggers.TriggerObserver
	for _, sensor := range SensorsData.sensors {
		triggerList = append(triggerList, sensor.GetTriggers()...)
	}
	return triggerList
}

func (d *DatabaseManager) GetTrigger(id string) (*Trigger, error) {
	mongoClient := d.getMongoClient()
	if mongoClient != nil {
		triggersCollection := mongoClient.Database(MONGO_DB_MIMIR).Collection(TRIGGERS_COLLECTION)
		objectId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, err
		}
		filter := bson.D{{Key: "_id", Value: objectId}}
		var trigger Trigger
		err = triggersCollection.FindOne(context.TODO(), filter).Decode(&trigger)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return nil, nil
			}
			return nil, err
		}
		return &trigger, nil
	}
	return nil, fmt.Errorf("mongo client not running")
}

func (d *DatabaseManager) InsertTrigger(t *Trigger) (*Trigger, error) {
	mongoClient := d.getMongoClient()
	if mongoClient == nil {
		return nil, fmt.Errorf("mongo client not running")
	}
	triggersCollection := mongoClient.Database(MONGO_DB_MIMIR).Collection(TRIGGERS_COLLECTION)
	result, err := triggersCollection.InsertOne(context.TODO(), t)
	if err != nil {
		return nil, err
	}

	triggerId, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, fmt.Errorf("error converting id for trigger")
	}
	t.ID = triggerId

	if t.Filename == "" {
		filename := t.BuildFileName("")

		contador := 1
		for {
			if _, err := os.Stat(filename); os.IsNotExist(err) {
				break
			}
			filename = t.BuildFileName(fmt.Sprintf("%d", contador))
			contador++
		}
		fmt.Println("filename ", filename)

		jsonString, err := json.MarshalIndent(t, "", "    ")
		if err != nil {
			fmt.Println("Error ", err)
		}

		os.WriteFile(filename, jsonString, os.ModePerm)

	}

	return t, nil
}

func RegisterTrigger(trigger *Trigger) {
	triggerObserver := trigger.BuildTriggerObserver()

	ActiveTriggers = append(ActiveTriggers, triggerObserver)
	for _, topic := range trigger.Topics {
		sensor, err := SensorsData.GetSensorByTopic(topic)
		if err == nil {
			sensor.Register(triggerObserver)
		}
	}
}

func triggerToDBTrigger(trigger *triggers.Trigger) *Trigger {
	id, err := primitive.ObjectIDFromHex(trigger.ID)
	if err != nil {
		id = primitive.NilObjectID
	}

	return &Trigger{
		ID:        id,
		Name:      trigger.Name,
		Filename:  "",
		Condition: Condition(trigger.GetConditionAsString()),
		Actions:   []Action{}, //TODO
	}
}

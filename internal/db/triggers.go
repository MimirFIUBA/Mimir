package db

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"mimir/internal/consts"
	"mimir/triggers"
	"os"
	"strings"

	"github.com/gookit/ini/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (t *Trigger) BuildFileName(suffix string) string {
	filename := strings.ReplaceAll(t.Name, " ", "_")
	if suffix != "" {
		filename += "_" + suffix
	}
	return ini.String(consts.TRIGGERS_DIR_CONFIG_NAME) + "/" + filename + consts.TRIGGERS_FILE_SUFFIX
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
	if t.Filename == "" {
		filename, err := insertTriggerFile(t)
		if err != nil {
			return nil, err
		}
		t.Filename = filename
	}

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

	return t, nil
}

func (d *DatabaseManager) UpdateTrigger(id string, triggerUpdate *Trigger) (*Trigger, error) {
	mongoClient := d.getMongoClient()
	if mongoClient == nil {
		return nil, fmt.Errorf("mongo client not running")
	}
	triggersCollection := mongoClient.Database(MONGO_DB_MIMIR).Collection(TRIGGERS_COLLECTION)
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{Key: "_id", Value: objectId}}
	update := bson.D{{Key: "$set", Value: triggerUpdate}}
	_, err = triggersCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, err
	}

	for _, trigger := range ActiveTriggers {
		fmt.Println("trigger: ", trigger)
		if trigger.GetID() == id {
			fmt.Println("trigger id match")
		}
	}

	return triggerUpdate, nil
}

func RegisterTrigger(trigger triggers.TriggerObserver, topics []string) {
	ActiveTriggers = append(ActiveTriggers, trigger)
	for _, topic := range topics {
		sensor, err := SensorsData.GetSensorByTopic(topic)
		if err == nil {
			sensor.Register(trigger)
		}
	}
}

func insertTriggerFile(t *Trigger) (string, error) {
	filename := getFilenameForTrigger(t)

	jsonString, err := json.MarshalIndent(t, "", "    ")
	if err != nil {
		slog.Error("error marshalling trigger to json", "error", err)
		return "", err
	}

	err = os.WriteFile(filename, jsonString, os.ModePerm)
	if err != nil {
		slog.Error("error writing trigger file", "filename", filename, "error", err)
		return filename, err
	}
	return filename, nil
}

func getFilenameForTrigger(t *Trigger) string {
	filename := t.BuildFileName("")

	contador := 1
	for {
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			break
		}
		filename = t.BuildFileName(fmt.Sprintf("%d", contador))
		contador++
	}

	return filename
}

func (d *DatabaseManager) FindTriggers(filter bson.D) ([]Trigger, error) {
	var results []Trigger
	mongoClient := d.getMongoClient()
	if mongoClient != nil {
		triggersCollection := mongoClient.Database(MONGO_DB_MIMIR).Collection(TRIGGERS_COLLECTION)
		cursor, err := triggersCollection.Find(context.TODO(), filter)
		if err != nil {
			return nil, err
		} else {
			defer cursor.Close(context.TODO())
		}

		if err = cursor.All(context.TODO(), &results); err != nil {
			return nil, err
		}
	}
	return results, nil
}

func (d *DatabaseManager) UpsertTriggers(triggersToUpsert []Trigger) (*mongo.BulkWriteResult, error) {
	mongoClient := d.getMongoClient()
	if mongoClient != nil {
		triggersCollection := mongoClient.Database(MONGO_DB_MIMIR).Collection(TRIGGERS_COLLECTION)

		writeModels := make([]mongo.WriteModel, 0)
		for _, trigger := range triggersToUpsert {
			writeModels = append(writeModels, mongo.WriteModel(
				mongo.NewUpdateOneModel().SetFilter(bson.D{{Key: "filename", Value: trigger.Filename}}).SetUpsert(true).SetUpdate(bson.M{"$set": trigger}),
			))
		}

		opts := options.BulkWrite().SetOrdered(false)
		return triggersCollection.BulkWrite(context.TODO(), writeModels, opts)
	}
	return nil, fmt.Errorf("mongo client not available")
}

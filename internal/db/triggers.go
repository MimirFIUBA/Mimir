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
	"time"

	"github.com/gookit/ini/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Trigger struct {
	ID           primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name         string             `json:"name,omitempty" bson:"name,omitempty"`
	Filename     string             `json:"filename,omitempty" bson:"filename,omitempty"`
	IsActive     bool               `json:"active" bson:"active"`
	Topics       []string           `json:"topics" bson:"topics"`
	Condition    Condition          `json:"condition" bson:"condition"`
	Actions      []Action           `json:"actions" bson:"actions,omitempty"`
	TrueActions  []Action           `json:"trueActions" bson:"trueActions,omitempty"`
	FalseActions []Action           `json:"falseActions" bson:"falseActions,omitempty"`
	Type         string             `json:"type" bson:"type"`
	Timeout      int                `json:"timeout,omitempty" bson:"timeout,omitempty"`
	Frequency    int                `json:"frequency,omitempty" bson:"frequency,omitempty"`
	Scheduled    bool               `json:"scheduled" bson:"scheduled"`
	CronExpr     string             `json:"cron,omitempty" bson:"cron,omitempty"`
}

type Condition string

type Action struct {
	Name          string `json:"name" bson:"name"`
	Type          string `json:"type" bson:"type"`
	Message       string `json:"message,omitempty" bson:"message,omitempty"`
	Topic         string `json:"topic,omitempty" bson:"topic,omitempty"`
	Command       string `json:"command,omitempty" bson:"command,omitempty"`
	CommandArgs   string `json:"args,omitempty" bson:"args,omitempty"`
	TriggerName   string `json:"triggerName,omitempty" bson:"triggerName,omitempty"`
	TriggerStatus bool   `json:"triggerStatus,omitempty" bson:"triggerStatus,omitempty"`
}

func (t *Trigger) BuildFileName(suffix string) string {
	filename := strings.ReplaceAll(t.Name, " ", "_")
	if suffix != "" {
		filename += "_" + suffix
	}
	return ini.String(consts.TRIGGERS_DIR_CONFIG_NAME) + "/" + filename + consts.TRIGGERS_FILE_SUFFIX
}

func (d *DatabaseManager) GetTriggers() []triggers.Trigger {
	var triggerList []triggers.Trigger
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

func (d *DatabaseManager) UpdateTrigger(id string, triggerUpdate *Trigger, actions []triggers.Action) (*Trigger, error) {
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
		if trigger.GetID() == id {
			switch trigger.GetType() {
			case triggers.EVENT_TRIGGER:
				eventTrigger, ok := trigger.(*triggers.EventTrigger)
				if !ok {
					return nil, fmt.Errorf("trigger type error")
				}
				eventTrigger.Name = triggerUpdate.Name
			case triggers.TIMER_TRIGGER:
				timerTrigger, ok := trigger.(*triggers.TimerTrigger)
				if !ok {
					return nil, fmt.Errorf("trigger type error")
				}
				timerTrigger.Name = triggerUpdate.Name
				timerTrigger.UpdateTimeout(time.Duration(triggerUpdate.Timeout) * time.Second)
			case triggers.FREQUENCY_TRIGGER:
				frequencyTrigger, ok := trigger.(*triggers.FrequencyTrigger)
				if !ok {
					return nil, fmt.Errorf("trigger type error")
				}
				frequencyTrigger.Name = triggerUpdate.Name
				frequencyTrigger.Frequency = time.Duration(triggerUpdate.Frequency) * time.Second
			default:
				return nil, fmt.Errorf("trigger type not recognized")
			}

			trigger.UpdateCondition(string(triggerUpdate.Condition))
			trigger.UpdateActions(actions, triggers.TriggerOptions{})
			trigger.SetStatus(triggerUpdate.IsActive)

			//TODO add more stuff to update
		}
	}

	filename, exists := TriggerFilenamesById[id]
	if exists {
		triggerUpdate.Filename = filename
		saveTriggerFile(triggerUpdate)
	}

	return triggerUpdate, nil
}

func RegisterTrigger(trigger triggers.Trigger, topics []string) {
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

func saveTriggerFile(triggerData *Trigger) error {
	jsonString, err := json.MarshalIndent(triggerData, "", "    ")
	if err != nil {
		return err
	}

	fileName := triggerData.Filename

	return os.WriteFile(fileName, jsonString, os.ModePerm)
}

func (d *DatabaseManager) DeleteTrigger(id string) error {
	mongoClient := d.getMongoClient()
	if mongoClient != nil {
		triggersCollection := mongoClient.Database(MONGO_DB_MIMIR).Collection(TRIGGERS_COLLECTION)
		objectId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return err
		}
		filter := bson.D{{Key: "_id", Value: objectId}}
		_, err = triggersCollection.DeleteOne(context.TODO(), filter)
		if err != nil {
			return err
		}
	}

	filename, exists := TriggerFilenamesById[id]
	if exists {
		deleteTriggerFile(filename)
	}
	removeTriggerFromWokflow(id)

	return nil
}

func deleteTriggerFile(filename string) {
	newName := strings.Replace(filename, ".json", "_deleted.json", 1)
	err := os.Rename(filename, newName)
	if err != nil {
		slog.Error("error renaming file for deletion", "error", err)
	}
}

func removeTriggerFromWokflow(id string) {
	indexToRemove := -1
	for i, trigger := range ActiveTriggers {
		if trigger.GetID() == id {
			trigger.StopWatching()
			indexToRemove = i
		}
	}

	if indexToRemove >= 0 {
		ActiveTriggers = append(ActiveTriggers[:indexToRemove], ActiveTriggers[indexToRemove+1:]...)
	}
}

func (d *DatabaseManager) DeactivateTriggers(ctx context.Context) {
	mongoClient := d.getMongoClient()
	if mongoClient != nil {
		filter := buildTriggerIdFilter(ActiveTriggers)
		update := bson.D{{Key: "$set", Value: bson.D{{Key: "active", Value: false}}}}
		triggersCollection := mongoClient.Database(MONGO_DB_MIMIR).Collection(TRIGGERS_COLLECTION)
		triggersCollection.UpdateMany(ctx, filter, update)
	}
}

func buildTriggerIdFilter(triggers []triggers.Trigger) bson.D {
	values := bson.A{}
	for _, trigger := range triggers {
		id, err := primitive.ObjectIDFromHex(trigger.GetID())
		if err != nil {
			slog.Warn("could not set object id for trigger", "id", trigger.GetID())
			continue
		}
		values = append(values, id)
	}

	return bson.D{{Key: "_id", Value: bson.D{{Key: "$in", Value: values}}}}
}

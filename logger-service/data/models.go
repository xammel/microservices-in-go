package data

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"common/constants"
)

var client *mongo.Client

func New(mongo *mongo.Client) Models {
	client = mongo

	return Models{
		LogEntry: LogEntry{},
	}
}

type Models struct {
	LogEntry LogEntry
}

type LogEntry struct {
	ID        string    `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string    `bson:"name" json:"name"`
	Data      string    `bson:"data" json:"data"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

func (l *LogEntry) Insert(entry LogEntry) error {
	collection := client.Database(constants.MongoDBName).Collection(constants.MongoDBCollectionName)

	// why not just pass entry as second arg..?
	_, error := collection.InsertOne(context.TODO(), LogEntry{
		Name:      entry.Name,
		Data:      entry.Data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if error != nil {
		log.Println("Error inserting into MongoDB collection logs:", error)
		return error
	}

	return nil
}

func (l *LogEntry) All() ([]*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	collection := client.Database(constants.MongoDBName).Collection(constants.MongoDBCollectionName)

	opts := options.Find()
	opts.SetSort(bson.D{{"created_at", -1}})

	cursor, err := collection.Find(context.TODO(), bson.D{}, opts)

	if err != nil {
		log.Println("Find error:", err)
		return nil, err
	}

	defer cursor.Close(ctx)

	var logs []*LogEntry

	for cursor.Next(ctx) {
		var item LogEntry

		error := cursor.Decode(&item)
		if error != nil {
			log.Println("Error decoding log into slice", error)
		} else {
			logs = append(logs, &item)
		}
	}

	return logs, nil
}

func (l *LogEntry) GetOne(id string) (*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	collection := client.Database(constants.MongoDBName).Collection(constants.MongoDBCollectionName)

	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var entry LogEntry
	err = collection.FindOne(ctx, bson.M{"_id": docID}).Decode(&entry)
	if err != nil {
		return nil, err
	}

	return &entry, nil
}

func (l *LogEntry) DropCollection() error {
	context, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	collection := client.Database(constants.MongoDBName).Collection(constants.MongoDBCollectionName)

	if err := collection.Drop(context); err != nil {
		return err
	}

	return nil
}

func (l *LogEntry) Update() (*mongo.UpdateResult, error) {
	context, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	collection := client.Database(constants.MongoDBName).Collection(constants.MongoDBCollectionName)

	docId, err := primitive.ObjectIDFromHex(l.ID)
	if err != nil {
		return nil, err
	}

	result, err := collection.UpdateOne(
		context,
		bson.M{"_id": docId},
		bson.D{
			{"$set", bson.D{
				{"name", l.Name},
				{"data", l.Data},
				{"updated_at", time.Now()},
			}},
		},
	)

	if err != nil {
		return nil, err
	}

	return result, nil
}

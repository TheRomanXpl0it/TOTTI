package db

import (
	"context"
	"sub/log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Flag struct {
	Flag           string              `bson:"flag" json:"flag"`
	Username       string              `bson:"username" json:"username"`
	ExploitName    string              `bson:"exploit_name" json:"exploit_name"`
	TeamIP         string              `bson:"team_ip" json:"team_ip"`
	Time           primitive.Timestamp `bson:"time" json:"time"`
	Status         uint                `bson:"status" json:"status"`
	ServerResponse uint                `bson:"server_response" json:"server_response"`
}

type DB struct {
	client     *mongo.Client
	collection *mongo.Collection
}

const (
	DB_NSUB uint = iota
	DB_SUB
	DB_SUCC
	DB_ERR
	DB_EXPIRED
)

const DBName = "flags"
const flagCollection = "flags"

func ConnectMongo(uri string) *DB {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("error connecting to db: %v\n", err)
	}
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatalf("error pinging the db: %v\n", err)
	}

	collection := client.Database(DBName).Collection(flagCollection)

	return &DB{
		client:     client,
		collection: collection,
	}
}

func (db *DB) Disconnect() {
	db.client.Disconnect(context.TODO())
}

func (db *DB) CreateFlagsCollection() error {
	indexModel := mongo.IndexModel{
		Keys:    bson.M{"flag": 1},
		Options: options.Index().SetUnique(true),
	}

	_, err := db.collection.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) InsertFlags(flags []Flag) error {
	var documents []interface{}
	for _, flag := range flags {
		documents = append(documents, flag)
	}

	opts := options.InsertMany().SetOrdered(false)
	inserted, err := db.collection.InsertMany(context.TODO(), documents, opts)
	log.Noticef("Inserted %v flags of %v posted\n", len(inserted.InsertedIDs), len(flags))
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil
		}
		return err
	}

	return nil
}

func (db *DB) GetAllFlags() ([]Flag, error) {
	cursor, err := db.collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var flags []Flag
	for cursor.Next(context.TODO()) {
		var flag Flag
		err := cursor.Decode(&flag)
		if err != nil {
			return nil, err
		}
		flags = append(flags, flag)
	}

	return flags, nil
}

func (db *DB) GetFlags(expirationTime primitive.Timestamp) ([]Flag, error) {
	filter := bson.M{
		"time":   bson.M{"$gt": expirationTime},
		"status": DB_NSUB,
	}

	options := options.Find().SetSort(bson.D{{Key: "time", Value: -1}}).SetAllowDiskUse(true)

	cursor, err := db.collection.Find(context.TODO(), filter, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var flags []Flag
	for cursor.Next(context.TODO()) {
		var flag Flag
		err := cursor.Decode(&flag)
		if err != nil {
			return nil, err
		}
		flags = append(flags, flag)
	}

	return flags, nil
}

func (db *DB) UpdateFlag(flag Flag) error {
	filter := bson.M{"flag": flag.Flag}
	update := bson.M{
		"$set": bson.M{
			"status":          flag.Status,
			"server_response": flag.ServerResponse,
		},
	}

	_, err := db.collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) UpdateExpiredFlags(expirationTime primitive.Timestamp) (int64, error) {
	filter := bson.M{
		"status": DB_NSUB,
		"time":   bson.M{"$lte": expirationTime},
	}

	update := bson.M{
		"$set": bson.M{
			"server_response": DB_EXPIRED,
		},
	}

	updated, err := db.collection.UpdateMany(context.TODO(), filter, update)
	if err != nil {
		return 0, err
	}

	return updated.ModifiedCount, nil
}

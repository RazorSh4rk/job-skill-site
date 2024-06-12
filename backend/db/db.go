package db

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type InterpretedListing struct {
	ID         primitive.ObjectID `json:"id" bson:"_id"`
	Link       string             `json:"link" bson:"link"`
	Title      string             `json:"title" bson:"title"`
	Company    string             `json:"company" bson:"company"`
	CompanyURL string             `json:"companyUrl" bson:"companyUrl"`
	Location   string             `json:"location" bson:"location"`
	Latitude   string             `json:"latitude" bson:"latitude"`
	Longitude  string             `json:"longitude" bson:"longitude"`
	Techstack  []string           `json:"techstack" bson:"techstack"`
}

var (
	db *mongo.Client
)

func Connect() {
	mongoUrl := os.Getenv("MONGO_URL")
	if mongoUrl == "" {
		log.Fatal("I need an MONGO_URL either in .env or as an envvar")
	}
	mongoOptions := options.Client().ApplyURI(mongoUrl)
	_db, err := mongo.Connect(context.Background(), mongoOptions)
	if err != nil {
		log.Fatal("Could not connect to mongo: ", err.Error())
	}
	db = _db

	err = db.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal("Could not ping mongo: ", err.Error())
	}

	log.Println("Mongo connected")
}

func GetCollectionNames() []string {
	dbnames, err := db.Database("custom_jobs").ListCollectionNames(context.Background(), bson.D{})
	if err != nil {
		log.Println("Error listing db:", err)
		return []string{}
	}
	return dbnames
}

func GetCollection(collectionName string) []InterpretedListing {
	collection := db.Database("custom_jobs").Collection(collectionName)
	cursor, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		log.Println("Error getting collection: ", err.Error())
	}
	var records []InterpretedListing
	err = cursor.All(context.Background(), &records)
	if err != nil {
		log.Println("Error unmarshaling records: ", err.Error())
	}
	return records
}

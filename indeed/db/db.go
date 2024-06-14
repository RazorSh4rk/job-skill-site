package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"razorsh4rk.github.io/indeedscrape/ai"
	"razorsh4rk.github.io/indeedscrape/scraper"
)

var (
	db         *mongo.Client
	collection *mongo.Collection
)

func Connect(collectionName string) {
	mongoUrl := os.Getenv("MONGO_URL")
	if mongoUrl == "" {
		log.Fatal("I need an MONGO_URL either in .env or as an envvar")
	}
	mongoOptions := options.Client().ApplyURI(mongoUrl)
	db, err := mongo.Connect(context.Background(), mongoOptions)
	if err != nil {
		log.Fatal("Could not connect to mongo: ", err.Error())
	}

	err = db.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal("Could not ping mongo: ", err.Error())
	}

	collection = db.Database("custom_jobs").Collection(collectionName)
	log.Println("Mongo connected")
}

func CheckExists(record scraper.IListing) bool {
	err := collection.FindOne(context.Background(), bson.D{{Key: "link", Value: record.GetLink()}}, nil).Err()
	return err == nil
}

func InsertOne(record ai.InterpretedListing) {
	exists := CheckExists(&record)
	if exists {
		log.Println("Duplicate ", record.Title, record.Company)
		return
	}

	_, err := collection.InsertOne(context.Background(), record)

	if err != nil {
		fmt.Println("Error trying to insert to database: ", err.Error())
	}
}

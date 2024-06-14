package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

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
	Techstack  []string           `json:"techstack" bson:"techstack"`
}

type GeoData struct {
	Lat string `json:"lat"`
	Lon string `json:"lon"`
}

func main() {

	mongoUrl := "mongodb+srv://a:a@cluster0.prl3mgz.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"

	mongoOptions := options.Client().ApplyURI(mongoUrl)
	db, err := mongo.Connect(context.Background(), mongoOptions)
	if err != nil {
		log.Fatal("Could not connect to mongo: ", err.Error())
	}

	err = db.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal("Could not ping mongo: ", err.Error())
	}

	dBase := db.Database("custom_jobs")
	log.Println("Mongo connected")

	documents := []string{"javascript_developer", "frontend_developer", "react_developer", "java_developer", "python_developer", "golang_developer", "ios_developer", "rust_developer", "flutter_developer", "php_developer", "android_developer", "devops_engineer"}

	for _, cName := range documents {
		fmt.Println(cName)
		collection := dBase.Collection(cName)

		cursor, err := collection.Find(context.Background(), bson.D{})
		if err != nil {
			log.Println("Error getting collection: ", err.Error())
		}
		var records []InterpretedListing
		err = cursor.All(context.Background(), &records)
		if err != nil {
			log.Println("Error unmarshaling records: ", err.Error())
		}

		for _, r := range records {
			if r.Location != "" {
				fmt.Println("\t", r.ID, r.Title, r.Company)
				time.Sleep(time.Duration(1) * time.Second)

				url := "https://geocode.maps.co/search?q=%s&api_key=6660b6109539d568840582itg6230e9"
				loc := trimRepeat(r.Location)
				loc = strings.ReplaceAll(loc, " ", "+")
				url = fmt.Sprintf(url, loc)

				res, err := http.Get(url)
				if err != nil {
					fmt.Println(err.Error())
					continue
				}
				defer res.Body.Close()

				body, err := io.ReadAll(res.Body)
				if err != nil {
					fmt.Println(err.Error())
					continue
				}

				var gData []GeoData
				err = json.Unmarshal(body, &gData)
				if err != nil {
					fmt.Println(err.Error())
					continue
				}

				if len(gData) > 0 {
					fmt.Println("...")
					lData := gData[0]
					instruction := bson.M{
						"$set": bson.M{
							"latitude":  lData.Lat,
							"longitude": lData.Lon,
						},
					}

					collection.UpdateOne(context.Background(), bson.D{{Key: "_id", Value: r.ID}}, instruction)
				}
			}
		}
	}
}

func trimRepeat(s string) string {
	for i := 0; i < len(s)/2; i++ {
		if s[i] != s[i+len(s)/2] {
			return s
		}
	}
	return s[0 : len(s)/2]
}

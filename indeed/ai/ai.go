package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	openai "github.com/sashabaranov/go-openai"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"razorsh4rk.github.io/indeedscrape/scraper"
)

type InterpretedListing struct {
	ID         primitive.ObjectID `json:"id" bson:"_id"`
	Link       string             `json:"link" bson:"link"`
	Title      string             `json:"title" bson:"title"`
	Company    string             `json:"company" bson:"company"`
	CompanyURL string             `json:"companyUrl" bson:"companyUrl"`
	Location   string             `json:"location" bson:"location"`
	Lat        string             `json:"latitude" bson:"latitude"`
	Lon        string             `json:"longitude" bson:"longitude"`
	Techstack  []string           `json:"techstack" bson:"techstack"`
}

func (l *InterpretedListing) GetLink() string {
	return l.Link
}

func (il *InterpretedListing) FromListing(l scraper.Listing) {
	il.ID = primitive.NewObjectID()
	il.Link = l.Link
	il.Title = l.Title
	il.Company = l.Company
	il.CompanyURL = l.CompanyURL
	il.Location = l.Location
	il.Lat = l.Lat
	il.Lon = l.Lon
}

type Tech struct {
	Technologies []string `json:"technologies"`
}

var (
	client *openai.Client
	prompt = `Get all the expected technologies from this job listing.
	Summarize the findings in a json format, without returning any text other than the json.
	Follow this format: { "technologies": ["nodejs", "golang", "devops", ...]}
	Do not include the backticks and "json" as if it was markdown, only the json string itself.
	Here is the description: %s`
)

func Connect() {
	key := os.Getenv("OPENAI_KEY")
	if key == "" {
		log.Fatal("I need an OPENAI_KEY either in .env or as an envvar")
	}
	client = openai.NewClient(key)
}

func CreateInterpretedListing(l scraper.Listing) InterpretedListing {
	var il InterpretedListing
	il.FromListing(l)

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			ResponseFormat: &openai.ChatCompletionResponseFormat{
				Type: openai.ChatCompletionResponseFormatTypeJSONObject,
			},
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: fmt.Sprintf(prompt, l.Description),
				},
			},
		},
	)

	if err != nil {
		log.Println("Error getting completion: ", err)
		return InterpretedListing{}
	}

	message := resp.Choices[0].Message.Content
	var tech Tech
	err = json.Unmarshal([]byte(message), &tech)
	if err != nil {
		log.Println("Error unmarshaling: ", message, err)
		return InterpretedListing{}
	}
	il.Techstack = tech.Technologies

	return il
}

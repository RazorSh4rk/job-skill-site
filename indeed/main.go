package main

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"razorsh4rk.github.io/indeedscrape/ai"
	"razorsh4rk.github.io/indeedscrape/cmd"
	"razorsh4rk.github.io/indeedscrape/db"
	"razorsh4rk.github.io/indeedscrape/scraper"
)

func main() {
	godotenv.Load()

	term := os.Getenv("SEARCHTERM")
	location := os.Getenv("LOCATION")
	pages, _ := strconv.Atoi(os.Getenv("PAGES"))
	skipPages, _ := strconv.Atoi(os.Getenv("SKIP_PAGES"))
	browser := os.Getenv("BROWSER")

	log.Println(term, location, pages, browser)

	ai.Connect()
	db.Connect(strings.ReplaceAll(term, " ", "_"))

	cmd.Process(&cmd.Config{
		Term:      term,
		Location:  location,
		Pages:     pages,
		SkipPages: skipPages,
		Browser:   browser,
	}, func(l scraper.Listing) {
		if db.CheckExists(&l) {
			log.Println("Duplicate, skipping", l.Company, l.Title)
			return
		}
		if !l.IsEmpty() {
			go func() {
				parsed := ai.CreateInterpretedListing(l)
				db.InsertOne(parsed)
			}()
		}
	})
}

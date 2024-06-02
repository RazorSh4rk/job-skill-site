package main

import (
	"github.com/joho/godotenv"
	"razorsh4rk.github.io/jobsite/db"
	"razorsh4rk.github.io/jobsite/server"
)

func main() {
	godotenv.Load()
	db.Connect()
	server.Start()
}

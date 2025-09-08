package main

import (
	"log"
	"os"
)

const standardPort = ":8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = standardPort
	}
	log.Println("\"Οὐ κλέπτω τὴν νίκην.\"")
	log.Println("I will not steal my victory.")

	log.Println("starting up.....")
	log.Println("starting up and getting env variables")
}

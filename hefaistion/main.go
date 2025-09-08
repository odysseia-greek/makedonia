package main

import (
	"log"
	"os"
)

const standardPort = ":50060"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = standardPort
	}
	log.Println("\"Ἡφαιστίων Ἀλέξανδρος.\"")
	log.Println("Hefaistion is Alexander.")

	log.Println("starting up.....")
	log.Println("starting up and getting env variables")
}
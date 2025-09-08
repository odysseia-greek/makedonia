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
	log.Println("\"ἐγὼ μὲν ἂν δεξαίμην, εἰ Ἀλέξανδρός εἰμι.\"")
	log.Println("If I were Alexander, I would accept it.")

	log.Println("starting up.....")
	log.Println("starting up and getting env variables")
}
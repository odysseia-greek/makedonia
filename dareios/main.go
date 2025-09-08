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
	log.Println("\"οἱ βασιλέως λόγοι πρὸς τοὺς Ἕλληνας οὐκ ἀληθεῖς εἰσίν.\"")
	log.Println("The king's words to the Greeks are not to be trusted.")

	log.Println("starting up.....")
	log.Println("starting up and getting env variables")
}
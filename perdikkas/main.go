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
	log.Println("\"Μὴ καταφρόνει τῆς ἐρήμου, βασιλεῦ· οὐ παντὶ ἀνθρώπῳ ὁμοίως ὁ κίνδυνος.\"")
	log.Println("Do not underestimate the desert, O king; danger does not treat all men equally.")

	log.Println("starting up.....")
	log.Println("starting up and getting env variables")
}
package main

import (
	"log"
	"os"

	"github.com/odysseia-greek/agora/plato/logging"
)

const standardPort = ":50060"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = standardPort
	}
	logging.System("\"οἱ βασιλέως λόγοι πρὸς τοὺς Ἕλληνας οὐκ ἀληθεῖς εἰσίν.\"")
	logging.System("The king's words to the Greeks are not to be trusted.")

	logging.System("starting up.....")
	logging.System("starting up and getting env variables")
}

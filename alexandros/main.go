package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/makedonia/alexandros/gateway"
	"github.com/odysseia-greek/makedonia/alexandros/routing"
)

const standardPort = ":8080"
const legacyPort = ":5000"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = standardPort
	}
	//https://patorjk.com/software/taag/#p=display&f=Crawford2&t=ALEXANDROS
	logging.System(`
  ____  _        ___  __ __   ____  ____   ___    ____   ___   _____
 /    || |      /  _]|  |  | /    ||    \ |   \  |    \ /   \ / ___/
|  o  || |     /  [_ |  |  ||  o  ||  _  ||    \ |  D  )     (   \_ 
|     || |___ |    _]|_   _||     ||  |  ||  D  ||    /|  O  |\__  |
|  _  ||     ||   [_ |     ||  _  ||  |  ||     ||    \|     |/  \ |
|  |  ||     ||     ||  |  ||  |  ||  |  ||     ||  .  \     |\    |
|__|__||_____||_____||__|__||__|__||__|__||_____||__|\_|\___/  \___|
                                                                    
`)
	logging.System("\"Ου κλέπτω την νίκην’\"")
	logging.System("\"I will not steal my victory\"")
	logging.System("starting up.....")
	logging.System("starting up and getting env variables")

	handler, err := gateway.CreateNewConfig(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	graphqlServer := routing.InitRoutes(handler)

	// Start legacy server on port 5000 in a goroutine
	go func() {
		logging.System(fmt.Sprintf("Legacy endpoint running on port %s", legacyPort))
		legacyErr := http.ListenAndServe(legacyPort, graphqlServer)
		if legacyErr != nil {
			log.Fatal("Legacy server failed to start: ", legacyErr)
		}
	}()

	logging.System(fmt.Sprintf("Server running on port %s", port))
	err = http.ListenAndServe(port, graphqlServer)
	if err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}

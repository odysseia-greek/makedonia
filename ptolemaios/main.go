package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/attike/aristophanes/comedy"
	"github.com/odysseia-greek/makedonia/ptolemaios/aigyptos"
	v1 "github.com/odysseia-greek/makedonia/ptolemaios/gen/go/v1"
	"google.golang.org/grpc"
)

const standardPort = ":50060"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = standardPort
	}
	//https://patorjk.com/software/taag/#p=display&f=Crawford2&t=ptolemaios
	logging.System(`
 ____  ______   ___   _        ___  ___ ___   ____  ____  ___   _____
|    \|      | /   \ | |      /  _]|   |   | /    ||    |/   \ / ___/
|  o  )      ||     || |     /  [_ | _   _ ||  o  | |  ||     (   \_ 
|   _/|_|  |_||  O  || |___ |    _]|  \_/  ||     | |  ||  O  |\__  |
|  |    |  |  |     ||     ||   [_ |   |   ||  _  | |  ||     |/  \ |
|  |    |  |  |     ||     ||     ||   |   ||  |  | |  ||     |\    |
|__|    |__|   \___/ |_____||_____||___|___||__|__||____|\___/  \___|
`)
	logging.System("\"Πτολεμαῖος δ᾿ ὁ Σωτὴρ ὄναρ εἶδε τὸν ἐν Σινώπῃ τοῦ Πλούτωνος κολοσσόν.\"")
	logging.System("Ptolemy Soter saw in a dream the colossal statue of Pluto in Sinope.")

	logging.System("starting up.....")
	logging.System("starting up and getting env variables")

	ctx := context.Background()
	cfg, err := aigyptos.CreateNewConfig(ctx)
	if err != nil {
		logging.Error(err.Error())
		log.Fatal("death has found me")
	}

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var server *grpc.Server

	server = grpc.NewServer(
		grpc.UnaryInterceptor(
			comedy.UnaryServerInterceptor(
				cfg.Streamer,
				comedy.WithHeaderKey(config.HeaderKey),
				comedy.WithContextKeyName(config.DefaultTracingName),
				comedy.WithCloseHop(),
			),
		),
	)

	v1.RegisterPtolemaiosServiceServer(server, cfg)

	logging.Info(fmt.Sprintf("Server listening on %s", port))
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

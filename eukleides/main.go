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
	"github.com/odysseia-greek/makedonia/eukleides/geometrias"
	pb "github.com/odysseia-greek/makedonia/eukleides/proto"
	"google.golang.org/grpc"
)

const standardPort = ":50060"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = standardPort
	}

	//https://patorjk.com/software/taag/#p=display&f=Crawford2&t=EUKLEIDES
	logging.System(`
   ___  __ __  __  _  _        ___  ____  ___      ___  _____
  /  _]|  |  ||  |/ ]| |      /  _]|    ||   \    /  _]/ ___/
 /  [_ |  |  ||  ' / | |     /  [_  |  | |    \  /  [_(   \_ 
|    _]|  |  ||    \ | |___ |    _] |  | |  D  ||    _]\__  |
|   [_ |  :  ||     ||     ||   [_  |  | |     ||   [_ /  \ |
|     ||     ||  .  ||     ||     | |  | |     ||     |\    |
|_____| \__,_||__|\_||_____||_____||____||_____||_____| \___|
`)

	logging.System("\"στοιχεῖα τῆς γεωμετρίας\"")
	logging.System("The Elements of Geometry")

	logging.System("starting up.....")
	logging.System("starting up and getting env variables")

	ctx := context.Background()
	cfg, err := geometrias.CreateNewConfig(ctx)
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

	pb.RegisterEukleidesServer(server, cfg)

	logging.Info(fmt.Sprintf("Server listening on %s", port))
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

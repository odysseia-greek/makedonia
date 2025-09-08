package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/odysseia-greek/agora/plato/logging"
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

	log.Println("\"στοιχεῖα τῆς γεωμετρίας\"")
	log.Println("The Elements of Geometry")

	log.Println("starting up.....")
	log.Println("starting up and getting env variables")

	ctx := context.Background()
	config, err := geometrias.CreateNewConfig(ctx)
	if err != nil {
		logging.Error(err.Error())
		log.Fatal("death has found me")
	}

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var server *grpc.Server

	server = grpc.NewServer(grpc.UnaryInterceptor(geometrias.Interceptor))

	pb.RegisterEukleidesServer(server, config)

	logging.Info(fmt.Sprintf("Server listening on %s", port))
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

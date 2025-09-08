package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/odysseia-greek/agora/plato/logging"
	v1 "github.com/odysseia-greek/makedonia/antigonos/gen/go/v1"
	"github.com/odysseia-greek/makedonia/antigonos/monophthalmus"
	"github.com/odysseia-greek/makedonia/filippos/hetairoi"
	"google.golang.org/grpc"
)

const standardPort = ":50060"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = standardPort
	}
	//https://patorjk.com/software/taag/#p=display&f=Crawford2&t=Antigonos
	logging.System(`
  ____  ____   ______  ____   ____   ___   ____    ___   _____
 /    ||    \ |      ||    | /    | /   \ |    \  /   \ / ___/
|  o  ||  _  ||      | |  | |   __||     ||  _  ||     (   \_ 
|     ||  |  ||_|  |_| |  | |  |  ||  O  ||  |  ||  O  |\__  |
|  _  ||  |  |  |  |   |  | |  |_ ||     ||  |  ||     |/  \ |
|  |  ||  |  |  |  |   |  | |     ||     ||  |  ||     |\    |
|__|__||__|__|  |__|  |____||___,_| \___/ |__|__| \___/  \___|
`)
	log.Println("\"Μὴ εἴπῃς Ἀντίγονε μοι, ἀλλὰ πρόσταττε ὡς βασιλεύς.\"")
	log.Println("Do not call me Antigonos, command me as a king.")

	log.Println("starting up.....")
	log.Println("starting up and getting env variables")

	ctx := context.Background()
	config, err := monophthalmus.CreateNewConfig(ctx)
	if err != nil {
		logging.Error(err.Error())
		log.Fatal("death has found me")
	}

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var server *grpc.Server

	server = grpc.NewServer(grpc.UnaryInterceptor(hetairoi.Interceptor))

	v1.RegisterAntigonosServiceServer(server, config)

	logging.Info(fmt.Sprintf("Server listening on %s", port))
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

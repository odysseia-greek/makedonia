package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/makedonia/filippos/hetairoi"
	"github.com/odysseia-greek/makedonia/perdikkas/epimeleia"
	v1 "github.com/odysseia-greek/makedonia/perdikkas/gen/go/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const standardPort = ":50060"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = standardPort
	}
	//https://patorjk.com/software/taag/#p=display&f=Crawford2&t=PERDIKKAS&x=none&v=4&h=4&w=80&we=false
	logging.System(`
 ____   ___  ____   ___    ____  __  _  __  _   ____  _____
|    \ /  _]|    \ |   \  |    ||  |/ ]|  |/ ] /    |/ ___/
|  o  )  [_ |  D  )|    \  |  | |  ' / |  ' / |  o  (   \_ 
|   _/    _]|    / |  D  | |  | |    \ |    \ |     |\__  |
|  | |   [_ |    \ |     | |  | |     ||     ||  _  |/  \ |
|  | |     ||  .  \|     | |  | |  .  ||  .  ||  |  |\    |
|__| |_____||__|\_||_____||____||__|\_||__|\_||__|__| \___|
`)

	logging.System("\"Μὴ καταφρόνει τῆς ἐρήμου, βασιλεῦ· οὐ παντὶ ἀνθρώπῳ ὁμοίως ὁ κίνδυνος.\"")
	logging.System("Do not underestimate the desert, O king; danger does not treat all men equally.")

	logging.System("starting up.....")
	logging.System("starting up and getting env variables")

	ctx := context.Background()
	config, err := epimeleia.CreateNewConfig(ctx)
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
	reflection.Register(server)

	v1.RegisterPerdikkasServiceServer(server, config)

	logging.Info(fmt.Sprintf("Server listening on %s", port))
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

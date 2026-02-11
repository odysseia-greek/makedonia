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
	v1 "github.com/odysseia-greek/makedonia/parmenion/gen/go/v1"
	"github.com/odysseia-greek/makedonia/parmenion/strategos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const standardPort = ":50060"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = standardPort
	}
	//https://patorjk.com/software/taag/#p=display&f=Crawford2&t=PARMENION&x=none&v=4&h=4&w=80&we=false
	logging.System(`
 ____   ____  ____   ___ ___    ___  ____   ____  ___   ____  
|    \ /    ||    \ |   |   |  /  _]|    \ |    |/   \ |    \ 
|  o  )  o  ||  D  )| _   _ | /  [_ |  _  | |  ||     ||  _  |
|   _/|     ||    / |  \_/  ||    _]|  |  | |  ||  O  ||  |  |
|  |  |  _  ||    \ |   |   ||   [_ |  |  | |  ||     ||  |  |
|  |  |  |  ||  .  \|   |   ||     ||  |  | |  ||     ||  |  |
|__|  |__|__||__|\_||___|___||_____||__|__||____|\___/ |__|__|
`)
	logging.System("\"ἐγὼ μὲν ἂν δεξαίμην, εἰ Ἀλέξανδρός εἰμι.\"")
	logging.System("If I were Alexander, I would accept it.")

	logging.System("starting up.....")
	logging.System("starting up and getting env variables")

	ctx := context.Background()
	cfg, err := strategos.CreateNewConfig(ctx)
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

	reflection.Register(server)

	v1.RegisterParmenionServiceServer(server, cfg)

	logging.Info(fmt.Sprintf("Server listening on %s", port))
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

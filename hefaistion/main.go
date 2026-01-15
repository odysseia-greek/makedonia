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
	v1 "github.com/odysseia-greek/makedonia/hefaistion/gen/go/v1"
	"github.com/odysseia-greek/makedonia/hefaistion/philia"
	"google.golang.org/grpc"
)

const standardPort = ":50060"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = standardPort
	}
	//https://patorjk.com/software/taag/#p=display&f=Crawford2&t=hefaistion
	logging.System(`
 __ __    ___  _____   ____  ____ _____ ______  ____  ___   ____  
|  |  |  /  _]|     | /    ||    / ___/|      ||    |/   \ |    \ 
|  |  | /  [_ |   __||  o  | |  (   \_ |      | |  ||     ||  _  |
|  _  ||    _]|  |_  |     | |  |\__  ||_|  |_| |  ||  O  ||  |  |
|  |  ||   [_ |   _] |  _  | |  |/  \ |  |  |   |  ||     ||  |  |
|  |  ||     ||  |   |  |  | |  |\    |  |  |   |  ||     ||  |  |
|__|__||_____||__|   |__|__||____|\___|  |__|  |____|\___/ |__|__|
`)
	logging.System("\"Ἡφαιστίων Ἀλέξανδρος.\"")
	logging.System("Hefaistion is Alexander.")

	logging.System("starting up.....")
	logging.System("starting up and getting env variables")

	ctx := context.Background()
	cfg, err := philia.CreateNewConfig(ctx)
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

	v1.RegisterHefastionServiceServer(server, cfg)

	cfg.StartReporting(ctx)

	logging.Info(fmt.Sprintf("Server listening on %s", port))
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

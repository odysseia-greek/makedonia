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
	v1 "github.com/odysseia-greek/makedonia/antigonos/gen/go/v1"
	"github.com/odysseia-greek/makedonia/antigonos/monophthalmus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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

	logging.System("\"Μὴ εἴπῃς Ἀντίγονε μοι, ἀλλὰ πρόσταττε ὡς βασιλεύς.\"")
	logging.System("Do not call me Antigonos, command me as a king.")

	logging.System("starting up.....")
	logging.System("starting up and getting env variables")

	ctx := context.Background()
	cfg, err := monophthalmus.CreateNewConfig(ctx)
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

	v1.RegisterAntigonosServiceServer(server, cfg)

	logging.Info(fmt.Sprintf("Server listening on %s", port))
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

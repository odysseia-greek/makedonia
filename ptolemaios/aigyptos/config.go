package aigyptos

import (
	"context"
	"os"

	"github.com/odysseia-greek/agora/archytas"
	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/logging"
	aristophanes "github.com/odysseia-greek/attike/aristophanes/comedy"
)

func CreateNewConfig(ctx context.Context) (*ExtendedServiceImpl, error) {
	tracer, err := aristophanes.NewClientTracer(aristophanes.DefaultAddress)
	healthy := tracer.WaitForHealthyState()
	if !healthy {
		logging.Error("tracing service not ready - restarting seems the only option")
		os.Exit(1)
	}

	streamer, err := tracer.Chorus(ctx)

	client, err := config.CreateOdysseiaClient()
	if err != nil {
		return nil, err
	}

	cache, err := archytas.CreateBadgerClient()
	if err != nil {
		return nil, err
	}

	version := os.Getenv(config.EnvVersion)

	return &ExtendedServiceImpl{
		Client:   client,
		Archytas: cache,
		Version:  version,
		Streamer: streamer,
	}, nil
}

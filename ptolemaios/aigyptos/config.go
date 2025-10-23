package aigyptos

import (
	"context"
	"os"

	"github.com/odysseia-greek/agora/archytas"
	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/makedonia/filippos/hetairoi"
)

func CreateNewConfig(ctx context.Context) (*ExtendedServiceImpl, error) {
	hetairoi.SetStreamer(ctx)

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
	}, nil
}

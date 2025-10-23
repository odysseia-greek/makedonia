package gateway

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/logging"
	aristophanes "github.com/odysseia-greek/attike/aristophanes/comedy"
	pb "github.com/odysseia-greek/attike/aristophanes/proto"
	"github.com/odysseia-greek/makedonia/antigonos/monophthalmus"
	"github.com/odysseia-greek/makedonia/eukleides/geometrias"
	pbe "github.com/odysseia-greek/makedonia/eukleides/proto"
)

func CreateNewConfig(ctx context.Context) (*AlexandrosHandler, error) {
	randomizer, err := config.CreateNewRandomizer()
	if err != nil {
		return nil, err
	}

	var tracer *aristophanes.ClientTracer
	var streamer pb.TraceService_ChorusClient
	var eukleides *geometrias.CounterClient
	var eukleidesStreamer pbe.Eukleides_CreateNewEntryClient

	maxRetries := 3
	retryDelay := 10 * time.Second

	for i := 1; i <= maxRetries; i++ {
		tracer, err = aristophanes.NewClientTracer(aristophanes.DefaultAddress)
		if err == nil {
			break
		}

		logging.Error(fmt.Sprintf("failed to create tracer (attempt %d/%d): %s", i, maxRetries, err.Error()))

		if i < maxRetries {
			time.Sleep(retryDelay)
		}
	}

	for i := 1; i <= maxRetries; i++ {
		streamer, err = tracer.Chorus(ctx)
		if err == nil {
			break
		}

		logging.Error(fmt.Sprintf("failed to create chorus streamer (attempt %d/%d): %s", i, maxRetries, err.Error()))
		if i < maxRetries {
			time.Sleep(retryDelay)
		}
	}

	healthy := tracer.WaitForHealthyState()
	if !healthy {
		logging.Error("tracing service not ready - starting up without traces")
	}

	for i := 1; i <= maxRetries; i++ {
		eukleides, err = geometrias.NewEukleidesClient("eukleides:50060")
		if err == nil {
			break
		}

		logging.Error(fmt.Sprintf("failed to create counter (attempt %d/%d): %s", i, maxRetries, err.Error()))

		if i < maxRetries {
			time.Sleep(retryDelay)
		}
	}

	if err != nil {
		logging.Error("giving up after 3 retries to connect to counter")
		os.Exit(1)
	}

	for i := 1; i <= maxRetries; i++ {
		eukleidesStreamer, err = eukleides.CreateNewEntry(ctx)
		if err == nil {
			break
		}

		logging.Error(fmt.Sprintf("failed to create eukleides streamer (attempt %d/%d): %s", i, maxRetries, err.Error()))
		if i < maxRetries {
			time.Sleep(retryDelay)
		}
	}

	healthy = eukleides.WaitForHealthyState()
	if !healthy {
		logging.Error("eukleides service not ready - starting up without counter")
	}

	fuzzyClientAddress := config.StringFromEnv("ANTIGONOS_SERVICE", "antigonos:50060")
	fuzzyClient, err := NewGenericGrpcClient[*monophthalmus.FuzzyClient](
		fuzzyClientAddress,
		monophthalmus.NewAntigonosClient,
	)

	if err != nil {
		logging.Error(err.Error())
		return nil, err
	}

	fuzzyClientHealthy := fuzzyClient.client.WaitForHealthyState()
	if !fuzzyClientHealthy {
		logging.Debug("fuzzu client not ready - restarting seems the only option")
		os.Exit(1)
	}

	return &AlexandrosHandler{
		Streamer:        streamer,
		Randomizer:      randomizer,
		FuzzyClient:     fuzzyClient,
		CounterStreamer: eukleidesStreamer,
		Counter:         eukleides,
	}, nil
}

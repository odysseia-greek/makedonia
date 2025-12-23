package gateway

import (
	"context"
	"fmt"
	"time"

	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/logging"
	aristophanes "github.com/odysseia-greek/attike/aristophanes/comedy"
	pb "github.com/odysseia-greek/attike/aristophanes/proto"
	"github.com/odysseia-greek/makedonia/antigonos/monophthalmus"
	"github.com/odysseia-greek/makedonia/eukleides/geometrias"
	pbe "github.com/odysseia-greek/makedonia/eukleides/proto"
	"github.com/odysseia-greek/makedonia/hefaistion/philia"
	"github.com/odysseia-greek/makedonia/parmenion/strategos"
	"github.com/odysseia-greek/makedonia/ptolemaios/aigyptos"
)

func CreateNewConfig(ctx context.Context) (*AlexandrosHandler, error) {
	start := time.Now()
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

	healthyTracer := false
	if tracer != nil {
		healthyTracer = tracer.WaitForHealthyState()
	}

	counterClientAddress := config.StringFromEnv("EUKLEIDES_SERVICE", "eukleides:50060")
	for i := 1; i <= maxRetries; i++ {
		eukleides, err = geometrias.NewEukleidesClient(counterClientAddress)
		if err == nil {
			break
		}

		logging.Error(fmt.Sprintf("failed to create counter (attempt %d/%d): %s", i, maxRetries, err.Error()))

		if i < maxRetries {
			time.Sleep(retryDelay)
		}
	}

	for i := 1; i <= maxRetries; i++ {
		if eukleides == nil {
			break
		}
		eukleidesStreamer, err = eukleides.CreateNewEntry(ctx)
		if err == nil {
			break
		}

		logging.Error(fmt.Sprintf("failed to create eukleides streamer (attempt %d/%d): %s", i, maxRetries, err.Error()))
		if i < maxRetries {
			time.Sleep(retryDelay)
		}
	}

	healthyEukleides := false
	if eukleides != nil {
		healthyEukleides = eukleides.WaitForHealthyState()
	}

	fuzzyClientAddress := config.StringFromEnv("ANTIGONOS_SERVICE", "antigonos:50060")
	fuzzyClient, err := NewGenericGrpcClient[*monophthalmus.FuzzyClient](
		fuzzyClientAddress,
		monophthalmus.NewAntigonosClient,
	)

	if err != nil {
		logging.Error(err.Error())
	}

	fuzzyClientHealthy := false
	if fuzzyClient != nil {
		fuzzyClientHealthy = fuzzyClient.client.WaitForHealthyState()
	}

	phraseClientAddress := config.StringFromEnv("PARMENION_SERVICE", "parmenion:50060")
	phraseClient, err := NewGenericGrpcClient[*strategos.PhraseClient](
		phraseClientAddress,
		strategos.NewParmenionClient,
	)
	if err != nil {
		logging.Error(err.Error())
	}

	phraseClientHealthy := false
	if phraseClient != nil {
		phraseClientHealthy = phraseClient.client.WaitForHealthyState()
	}

	exactClientAddress := config.StringFromEnv("HEFAISTION_SERVICE", "hefaistion:50060")
	exactClient, err := NewGenericGrpcClient[*philia.ExactClient](
		exactClientAddress,
		philia.NewHefaistionClient,
	)
	if err != nil {
		logging.Error(err.Error())
	}

	exactClientHealthy := false
	if exactClient != nil {
		exactClientHealthy = exactClient.client.WaitForHealthyState()
	}

	extendedClientAddress := config.StringFromEnv("PTOLEMAIOS_SERVICE", "ptolemaios:50060")
	extendedClient, err := NewGenericGrpcClient[*aigyptos.ExtendedClient](
		extendedClientAddress,
		aigyptos.NewPtolemaiosClient,
	)
	if err != nil {
		logging.Error(err.Error())
	}

	extendedClientHealthy := false
	if extendedClient != nil {
		extendedClientHealthy = extendedClient.client.WaitForHealthyState()
	}

	elapsed := time.Since(start)

	logging.System(fmt.Sprintf(`Alexandros Configuration Overview:
- Initialization Time: %s
- Tracer Service:      %v (Address: %s)
- Eukleides Service:   %v (Address: %s)
- Antigonos Service:   %v (Address: %s)
- Parmenion Service:   %v (Address: %s)
- Hefaistion Service:  %v (Address: %s)
- Ptolemaios Service:  %v (Address: %s)
`,
		elapsed,
		healthyTracer, aristophanes.DefaultAddress,
		healthyEukleides, counterClientAddress,
		fuzzyClientHealthy, fuzzyClientAddress,
		phraseClientHealthy, phraseClientAddress,
		exactClientHealthy, exactClientAddress,
		extendedClientHealthy, extendedClientAddress,
	))

	return &AlexandrosHandler{
		Streamer:        streamer,
		Randomizer:      randomizer,
		FuzzyClient:     fuzzyClient,
		ExactClient:     exactClient,
		PhraseClient:    phraseClient,
		ExtendedClient:  extendedClient,
		CounterStreamer: eukleidesStreamer,
		Counter:         eukleides,
	}, nil
}

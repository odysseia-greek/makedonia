package gateway

import (
	"context"
	"time"

	"github.com/odysseia-greek/agora/hesiodos"
	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/randomizer"
	arv1 "github.com/odysseia-greek/attike/aristophanes/gen/go/v1"
	"github.com/odysseia-greek/makedonia/antigonos/monophthalmus"
	"github.com/odysseia-greek/makedonia/eukleides/geometrias"
	pbe "github.com/odysseia-greek/makedonia/eukleides/proto"
	"github.com/odysseia-greek/makedonia/hefaistion/philia"
	"github.com/odysseia-greek/makedonia/parmenion/strategos"
	"github.com/odysseia-greek/makedonia/perdikkas/epimeleia"
	"github.com/odysseia-greek/makedonia/ptolemaios/aigyptos"
	"google.golang.org/grpc/metadata"
)

type AlexandrosHandler struct {
	Streamer        arv1.TraceService_ChorusClient
	CounterStreamer pbe.Eukleides_CreateNewEntryClient
	Counter         *geometrias.CounterClient
	Randomizer      randomizer.Random
	FuzzyClient     *hesiodos.GenericGrpcClient[*monophthalmus.FuzzyClient]
	ExactClient     *hesiodos.GenericGrpcClient[*philia.ExactClient]
	ExtendedClient  *hesiodos.GenericGrpcClient[*aigyptos.ExtendedClient]
	PhraseClient    *hesiodos.GenericGrpcClient[*strategos.PhraseClient]
	PartialClient   *hesiodos.GenericGrpcClient[*epimeleia.PartialClient]
}

func (a *AlexandrosHandler) outgoingCtx(parent context.Context) (context.Context, context.CancelFunc, string) {
	ctx, cancel := context.WithTimeout(parent, 30*time.Second)

	reqID, _ := parent.Value(config.HeaderKey).(string)
	sessionID, _ := parent.Value(config.SessionIdKey).(string)

	kvs := make([]string, 0, 4)

	if reqID != "" {
		kvs = append(kvs, config.HeaderKey, reqID)
	}
	if sessionID != "" {
		kvs = append(kvs, config.SessionIdKey, sessionID)
	}

	if len(kvs) > 0 {
		ctx = metadata.AppendToOutgoingContext(ctx, kvs...)
	}

	return ctx, cancel, sessionID
}

package strategos

import (
	"context"
	"fmt"
	"time"

	"github.com/odysseia-greek/agora/archytas"
	"github.com/odysseia-greek/agora/aristoteles"
	"github.com/odysseia-greek/agora/plato/randomizer"
	"github.com/odysseia-greek/agora/plato/service"
	pbar "github.com/odysseia-greek/attike/aristophanes/proto"
	koinos "github.com/odysseia-greek/makedonia/filippos/gen/go/koinos/v1"
	v1 "github.com/odysseia-greek/makedonia/parmenion/gen/go/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

type PhraseService interface {
	WaitForHealthyState() bool
	Search(ctx context.Context, request *koinos.SearchQuery) (*v1.SearchResponse, error)
}

const (
	DEFAULTADDRESS string = "localhost:50060"
)

type PhraseServiceImpl struct {
	Elastic    aristoteles.Client
	Index      string
	Version    string
	Randomizer randomizer.Random
	Client     service.OdysseiaClient
	Streamer   pbar.TraceService_ChorusClient
	Archytas   archytas.Client
	v1.UnimplementedParmenionServiceServer
}

type PhraseServiceClient struct {
	Impl PhraseService
}
type PhraseClient struct {
	prhase v1.ParmenionServiceClient
}

func NewParmenionClient(address string) (*PhraseClient, error) {
	if address == "" {
		address = DEFAULTADDRESS
	}
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to tracing service: %w", err)
	}
	client := v1.NewParmenionServiceClient(conn)
	return &PhraseClient{prhase: client}, nil
}

func (p *PhraseClient) WaitForHealthyState() bool {
	timeout := 30 * time.Second
	checkInterval := 1 * time.Second
	endTime := time.Now().Add(timeout)

	for time.Now().Before(endTime) {
		response, err := p.Health(context.Background(), &emptypb.Empty{})
		if err == nil && response.Healthy {
			return true
		}

		time.Sleep(checkInterval)
	}

	return false
}

func (p *PhraseClient) Health(ctx context.Context, request *emptypb.Empty) (*koinos.HealthResponse, error) {
	return p.prhase.Health(ctx, request)
}

func (p *PhraseClient) Search(ctx context.Context, request *koinos.SearchQuery) (*v1.SearchResponse, error) {
	return p.prhase.Search(ctx, request)
}

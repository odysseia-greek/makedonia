package monophthalmus

import (
	"context"
	"fmt"
	"time"

	"github.com/odysseia-greek/agora/archytas"
	"github.com/odysseia-greek/agora/aristoteles"
	"github.com/odysseia-greek/agora/plato/randomizer"
	arv1 "github.com/odysseia-greek/attike/aristophanes/gen/go/v1"
	v1 "github.com/odysseia-greek/makedonia/antigonos/gen/go/v1"
	koinos "github.com/odysseia-greek/makedonia/filippos/gen/go/koinos/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

type FuzzyService interface {
	WaitForHealthyState() bool
	Search(ctx context.Context, request *koinos.SearchQuery) (*v1.SearchResponse, error)
}

const (
	DEFAULTADDRESS string = "localhost:50060"
)

type FuzzyServiceImpl struct {
	Elastic    aristoteles.Client
	Index      string
	Version    string
	Randomizer randomizer.Random
	Streamer   arv1.TraceService_ChorusClient
	Archytas   archytas.Client
	v1.UnimplementedAntigonosServiceServer
}

type FuzzyServiceClient struct {
	Impl FuzzyService
}
type FuzzyClient struct {
	fuzzy v1.AntigonosServiceClient
}

func NewAntigonosClient(address string) (*FuzzyClient, error) {
	if address == "" {
		address = DEFAULTADDRESS
	}
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to tracing service: %w", err)
	}
	client := v1.NewAntigonosServiceClient(conn)
	return &FuzzyClient{fuzzy: client}, nil
}

func (f *FuzzyClient) WaitForHealthyState() bool {
	timeout := 30 * time.Second
	checkInterval := 1 * time.Second
	endTime := time.Now().Add(timeout)

	for time.Now().Before(endTime) {
		response, err := f.Health(context.Background(), &emptypb.Empty{})
		if err == nil && response.Healthy {
			return true
		}

		time.Sleep(checkInterval)
	}

	return false
}

func (f *FuzzyClient) Health(ctx context.Context, request *emptypb.Empty) (*koinos.HealthResponse, error) {
	return f.fuzzy.Health(ctx, request)
}

func (f *FuzzyClient) Search(ctx context.Context, request *koinos.SearchQuery) (*v1.SearchResponse, error) {
	return f.fuzzy.Search(ctx, request)
}

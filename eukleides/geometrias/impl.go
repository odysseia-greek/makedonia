package geometrias

import (
	"context"
	"fmt"
	"github.com/odysseia-greek/agora/archytas"
	"github.com/odysseia-greek/agora/aristoteles"
	pbar "github.com/odysseia-greek/attike/aristophanes/proto"
	pb "github.com/odysseia-greek/makedonia/eukleides/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type CounterService interface {
	WaitForHealthyState() bool
	CreateNewEntry(ctx context.Context) (pb.Eukleides_CreateNewEntryClient, error)
}

const (
	DEFAULTADDRESS string = "localhost:50060"
)

type CounterServiceImpl struct {
	Elastic  aristoteles.Client
	Index    string
	Version  string
	Archytas archytas.Client
	Streamer pbar.TraceService_ChorusClient
	pb.UnimplementedEukleidesServer
}
type CounterServiceClient struct {
	Impl CounterService
}
type CounterClient struct {
	counter pb.EukleidesClient
}

func NewEukleidesClient(address string) (*CounterClient, error) {
	if address == "" {
		address = DEFAULTADDRESS
	}
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to tracing service: %w", err)
	}
	client := pb.NewEukleidesClient(conn)
	return &CounterClient{counter: client}, nil
}

func (m *CounterClient) WaitForHealthyState() bool {
	timeout := 30 * time.Second
	checkInterval := 1 * time.Second
	endTime := time.Now().Add(timeout)

	for time.Now().Before(endTime) {
		response, err := m.Health(context.Background(), &pb.HealthRequest{})
		if err == nil && response.Health {
			return true
		}

		time.Sleep(checkInterval)
	}

	return false
}

func (m *CounterClient) Health(ctx context.Context, request *pb.HealthRequest) (*pb.HealthResponse, error) {
	return m.counter.Health(ctx, request)
}

func (m *CounterClient) CreateNewEntry(ctx context.Context) (pb.Eukleides_CreateNewEntryClient, error) {
	return m.counter.CreateNewEntry(ctx)
}

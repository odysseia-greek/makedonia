package geometrias

import (
	"context"
	"fmt"
	"time"

	arv1 "github.com/odysseia-greek/attike/aristophanes/gen/go/v1"
	pb "github.com/odysseia-greek/makedonia/eukleides/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type CounterService interface {
	WaitForHealthyState() bool
	CreateNewEntry(ctx context.Context) (pb.Eukleides_CreateNewEntryClient, error)
	RetrieveTopFive(ctx context.Context, in *pb.TopFiveRequest) (*pb.TopFiveResponse, error)
	RetrieveTopFiveService(ctx context.Context, in *pb.TopFiveServiceRequest) (*pb.TopFive, error)
	RetrieveTopFiveForSession(ctx context.Context, in *pb.TopFiveSessionRequest) (*pb.TopFiveResponse, error)
}

const (
	DEFAULTADDRESS string = "localhost:50060"
)

type CounterServiceImpl struct {
	Version  string
	store    *Store
	Streamer arv1.TraceService_ChorusClient
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
		if err == nil && response.Healthy {
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

func (m *CounterClient) RetrieveTopFive(ctx context.Context, in *pb.TopFiveRequest) (*pb.TopFiveResponse, error) {
	return m.counter.RetrieveTopFive(ctx, in)
}

func (m *CounterClient) RetrieveTopFiveService(ctx context.Context, in *pb.TopFiveServiceRequest) (*pb.TopFive, error) {
	return m.counter.RetrieveTopFiveService(ctx, in)
}

func (m *CounterClient) RetrieveTopFiveForSession(ctx context.Context, in *pb.TopFiveSessionRequest) (*pb.TopFiveResponse, error) {
	return m.counter.RetrieveTopFiveForSession(ctx, in)
}

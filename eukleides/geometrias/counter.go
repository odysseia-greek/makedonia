package geometrias

import (
	"context"
	"io"
	"time"

	pb "github.com/odysseia-greek/makedonia/eukleides/proto"
)

func (c *CounterServiceImpl) Health(ctx context.Context, request *pb.HealthRequest) (*pb.HealthResponse, error) {
	return &pb.HealthResponse{
		Healthy: true,
		Time:    time.Now().String(),
		Version: c.Version,
	}, nil
}

// Server-side implementations (barebones) for Eukleides service
// CreateNewEntry receives a stream of CountCreationRequestSet and returns an ack.
func (c *CounterServiceImpl) CreateNewEntry(stream pb.Eukleides_CreateNewEntryServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.CountStreamResponse{Ack: "Received"})
		}
		if err != nil {
			return err
		}

		go c.createOrUpdate(in)
	}
}

func (c *CounterServiceImpl) createOrUpdate(request *pb.CountCreationRequestSet) {
	now := time.Now().UTC()
	for _, req := range request.Request {
		c.store.Inc(req.SessionId, req.ServiceName, req.Word, now)
	}
}

// RetrieveTopFive returns an empty TopFiveResponse for now (to be implemented later).
func (c *CounterServiceImpl) RetrieveTopFive(ctx context.Context, in *pb.TopFiveRequest) (*pb.TopFiveResponse, error) {
	return &pb.TopFiveResponse{TopFive: c.store.TopFiveGlobal()}, nil
}

// RetrieveTopFiveService returns an empty TopFive placeholder (to be implemented later).
func (c *CounterServiceImpl) RetrieveTopFiveService(ctx context.Context, in *pb.TopFiveServiceRequest) (*pb.TopFive, error) {
	top := c.store.TopFiveByService(in.Name)
	if len(top) == 0 {
		return &pb.TopFive{}, nil
	}
	return top[0], nil
}

// RetrieveTopFiveForSession returns an empty TopFiveResponse for the given session (to be implemented later).
func (c *CounterServiceImpl) RetrieveTopFiveForSession(ctx context.Context, in *pb.TopFiveSessionRequest) (*pb.TopFiveResponse, error) {
	return &pb.TopFiveResponse{TopFive: c.store.TopFiveForSession(in.SessionId)}, nil
}

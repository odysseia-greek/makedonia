package gateway

import (
	"context"
	"os"
	"time"

	"github.com/odysseia-greek/makedonia/alexandros/graph/model"
	"github.com/odysseia-greek/makedonia/antigonos/monophthalmus"
	"github.com/odysseia-greek/makedonia/hefaistion/philia"
	"github.com/odysseia-greek/makedonia/parmenion/strategos"
	"github.com/odysseia-greek/makedonia/perdikkas/epimeleia"
	"google.golang.org/protobuf/types/known/emptypb"

	koinos "github.com/odysseia-greek/makedonia/filippos/gen/go/koinos/v1"
)

func (a *AlexandrosHandler) Health(ctx context.Context) (*model.AggregatedHealthResponse, error) {
	var services []*model.ServiceHealth
	allHealthy := true

	type healthCheck struct {
		name   string
		client func(ctx context.Context) (bool, *model.DatabaseInfo, *string) // healthy, dbHealthy, version
	}

	checks := []healthCheck{
		{
			name: "fuzzy",
			client: func(ctx context.Context) (bool, *model.DatabaseInfo, *string) {
				var resp *koinos.HealthResponse
				err := a.FuzzyClient.CallWithReconnect(func(c *monophthalmus.FuzzyClient) error {
					var innerErr error
					resp, innerErr = c.Health(ctx, &emptypb.Empty{})
					return innerErr
				})
				if err != nil || resp == nil {
					return false, nil, nil
				}
				databaseHealth := &model.DatabaseInfo{
					Healthy:       resp.DatabaseHealth.Healthy,
					ClusterName:   &resp.DatabaseHealth.ClusterName,
					ServerName:    &resp.DatabaseHealth.ServerName,
					ServerVersion: &resp.DatabaseHealth.ServerVersion,
				}
				return resp.GetHealthy(), databaseHealth, ptr(resp.GetVersion())
			},
		},
		{
			name: "exact",
			client: func(ctx context.Context) (bool, *model.DatabaseInfo, *string) {
				var resp *koinos.HealthResponse
				err := a.ExactClient.CallWithReconnect(func(c *philia.ExactClient) error {
					var innerErr error
					resp, innerErr = c.Health(ctx, &emptypb.Empty{})
					return innerErr
				})
				if err != nil || resp == nil {
					return false, nil, nil
				}
				databaseHealth := &model.DatabaseInfo{
					Healthy:       resp.DatabaseHealth.Healthy,
					ClusterName:   &resp.DatabaseHealth.ClusterName,
					ServerName:    &resp.DatabaseHealth.ServerName,
					ServerVersion: &resp.DatabaseHealth.ServerVersion,
				}
				return resp.GetHealthy(), databaseHealth, ptr(resp.GetVersion())
			},
		},
		{
			name: "phrase",
			client: func(ctx context.Context) (bool, *model.DatabaseInfo, *string) {
				var resp *koinos.HealthResponse
				err := a.PhraseClient.CallWithReconnect(func(c *strategos.PhraseClient) error {
					var innerErr error
					resp, innerErr = c.Health(ctx, &emptypb.Empty{})
					return innerErr
				})
				if err != nil || resp == nil {
					return false, nil, nil
				}
				databaseHealth := &model.DatabaseInfo{
					Healthy:       resp.DatabaseHealth.Healthy,
					ClusterName:   &resp.DatabaseHealth.ClusterName,
					ServerName:    &resp.DatabaseHealth.ServerName,
					ServerVersion: &resp.DatabaseHealth.ServerVersion,
				}
				return resp.GetHealthy(), databaseHealth, ptr(resp.GetVersion())
			},
		},
		{
			name: "partial",
			client: func(ctx context.Context) (bool, *model.DatabaseInfo, *string) {
				var resp *koinos.HealthResponse
				err := a.PartialClient.CallWithReconnect(func(c *epimeleia.PartialClient) error {
					var innerErr error
					resp, innerErr = c.Health(ctx, &emptypb.Empty{})
					return innerErr
				})
				if err != nil || resp == nil {
					return false, nil, nil
				}
				databaseHealth := &model.DatabaseInfo{
					Healthy:       resp.DatabaseHealth.Healthy,
					ClusterName:   &resp.DatabaseHealth.ClusterName,
					ServerName:    &resp.DatabaseHealth.ServerName,
					ServerVersion: &resp.DatabaseHealth.ServerVersion,
				}
				return resp.GetHealthy(), databaseHealth, ptr(resp.GetVersion())
			},
		},
	}

	for _, check := range checks {
		outCtx, cancel, _ := a.outgoingCtx(ctx)
		defer cancel()
		healthy, dbHealthy, version := check.client(outCtx)
		cancel()

		serviceHealth := &model.ServiceHealth{
			Name:         check.name,
			Healthy:      healthy,
			DatabaseInfo: dbHealthy,
			Version:      version,
		}

		if !healthy {
			allHealthy = false
		}

		services = append(services, serviceHealth)
	}

	return &model.AggregatedHealthResponse{
		Healthy:  allHealthy,
		Time:     ptr(time.Now().Format(time.RFC3339)),
		Version:  ptr(os.Getenv("VERSION")),
		Services: services,
	}, nil
}

func ptr[T any](v T) *T {
	return &v
}

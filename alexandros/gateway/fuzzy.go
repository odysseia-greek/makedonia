package gateway

import (
	"context"

	"github.com/odysseia-greek/makedonia/alexandros/graph/model"
	antigonosv1 "github.com/odysseia-greek/makedonia/antigonos/gen/go/v1"
	"github.com/odysseia-greek/makedonia/antigonos/monophthalmus"
	pbe "github.com/odysseia-greek/makedonia/eukleides/proto"
	koinos "github.com/odysseia-greek/makedonia/filippos/gen/go/koinos/v1"
)

func (a *AlexandrosHandler) Fuzzy(ctx context.Context, request *koinos.SearchQuery) (*model.SearchResponse, error) {
	outCtx, cancel, sessionId := a.outgoingCtx(ctx)
	defer cancel()

	eukleidesUpdate := pbe.CountCreationRequest{
		Word:        request.Word,
		ServiceName: "fuzzy",
		SearchType:  "fuzzy",
		SessionId:   sessionId,
	}

	go a.pushToEukleides(&eukleidesUpdate)

	var grpcResponse *antigonosv1.SearchResponse

	err := a.FuzzyClient.CallWithReconnect(func(client *monophthalmus.FuzzyClient) error {
		var innerErr error
		grpcResponse, innerErr = client.Search(outCtx, request)
		return innerErr
	})
	if err != nil {
		return nil, err
	}

	lemmas := parseResults(grpcResponse.Results)

	resp := &model.SearchResponse{
		Results: lemmas, // if you ever build this from scratch, prefer [] over nil
		PageInfo: &model.PageInfo{
			Page:  grpcResponse.PageInfo.Page,
			Size:  grpcResponse.PageInfo.Size,
			Total: grpcResponse.PageInfo.Total,
		},
	}
	return resp, nil
}

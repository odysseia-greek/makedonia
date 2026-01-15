package gateway

import (
	"context"

	"github.com/odysseia-greek/makedonia/alexandros/graph/model"
	pbe "github.com/odysseia-greek/makedonia/eukleides/proto"
	koinos "github.com/odysseia-greek/makedonia/filippos/gen/go/koinos/v1"
	"github.com/odysseia-greek/makedonia/perdikkas/epimeleia"
	perdikkasv1 "github.com/odysseia-greek/makedonia/perdikkas/gen/go/v1"
)

func (a *AlexandrosHandler) Partial(ctx context.Context, request *koinos.SearchQuery) (*model.SearchResponse, error) {
	outCtx, cancel, sessionId := a.outgoingCtx(ctx)
	defer cancel()

	eukleidesUpdate := pbe.CountCreationRequest{
		Word:        request.Word,
		ServiceName: "partial",
		SearchType:  "partial",
		SessionId:   sessionId,
	}

	go a.pushToEukleides(&eukleidesUpdate)

	var grpcResponse *perdikkasv1.SearchResponse

	err := a.PartialClient.CallWithReconnect(func(client *epimeleia.PartialClient) error {
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

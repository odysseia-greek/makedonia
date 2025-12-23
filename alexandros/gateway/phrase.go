package gateway

import (
	"github.com/odysseia-greek/makedonia/alexandros/graph/model"
	pbe "github.com/odysseia-greek/makedonia/eukleides/proto"
	koinos "github.com/odysseia-greek/makedonia/filippos/gen/go/koinos/v1"
	parmenionv1 "github.com/odysseia-greek/makedonia/parmenion/gen/go/v1"
	"github.com/odysseia-greek/makedonia/parmenion/strategos"
)

func (a *AlexandrosHandler) Phrase(request *koinos.SearchQuery, requestID, sessionId string) (*model.SearchResponse, error) {
	phraseClientCtx, cancel := a.createRequestHeader(requestID, sessionId)
	defer cancel()

	eukleidesUpdate := pbe.CountCreationRequest{
		Word:        request.Word,
		ServiceName: "phrase",
		SearchType:  "phrase",
		SessionId:   sessionId,
	}

	go a.pushToEukleides(&eukleidesUpdate)

	var grpcResponse *parmenionv1.SearchResponse

	err := a.PhraseClient.CallWithReconnect(func(client *strategos.PhraseClient) error {
		var innerErr error
		grpcResponse, innerErr = client.Search(phraseClientCtx, request)
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

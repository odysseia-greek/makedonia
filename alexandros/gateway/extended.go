package gateway

import (
	"context"

	"github.com/odysseia-greek/makedonia/alexandros/graph/model"
	pbe "github.com/odysseia-greek/makedonia/eukleides/proto"
	"github.com/odysseia-greek/makedonia/ptolemaios/aigyptos"
	v1 "github.com/odysseia-greek/makedonia/ptolemaios/gen/go/v1"
)

func (a *AlexandrosHandler) Extended(ctx context.Context, request *v1.ExtendedSearch) (*model.AnalyzeTextResponse, error) {
	outCtx, cancel, sessionId := a.outgoingCtx(ctx)
	defer cancel()

	eukleidesUpdate := pbe.CountCreationRequest{
		Word:        request.Word,
		ServiceName: "extended",
		SearchType:  "textSearch",
		SessionId:   sessionId,
	}

	go a.pushToEukleides(&eukleidesUpdate)

	var grpcResponse *v1.ExtendedSearchResponse

	err := a.ExtendedClient.CallWithReconnect(func(client *aigyptos.ExtendedClient) error {
		var innerErr error
		grpcResponse, innerErr = client.ExtendedSearch(outCtx, request)
		return innerErr
	})
	if err != nil {
		return nil, err
	}

	resp := &model.AnalyzeTextResponse{
		Conjugations: nil,
		Texts:        nil,
		Rootword:     &grpcResponse.FoundInText.Rootword,
	}

	for _, conj := range grpcResponse.FoundInText.Conjugations {
		resp.Conjugations = append(resp.Conjugations, &model.ConjugationResponse{
			Rule: &conj.Rule,
			Word: &conj.Word,
		})
	}

	for _, text := range grpcResponse.FoundInText.Texts {
		textModel := &model.AnalyzeResult{
			Author:        &text.Author,
			Book:          &text.Book,
			Reference:     &text.Reference,
			ReferenceLink: &text.ReferenceLink,
			Text: &model.Rhema{
				Greek:        &text.Text.Greek,
				Section:      &text.Text.Section,
				Translations: nil,
			},
		}

		for _, translation := range text.Text.Translations {
			textModel.Text.Translations = append(textModel.Text.Translations, &translation)
		}

		resp.Texts = append(resp.Texts, textModel)
	}

	return resp, nil
}

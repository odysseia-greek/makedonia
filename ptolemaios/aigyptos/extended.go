package aigyptos

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/agora/plato/models"
	"github.com/odysseia-greek/agora/plato/service"
	pb "github.com/odysseia-greek/attike/aristophanes/proto"
	koinos "github.com/odysseia-greek/makedonia/filippos/gen/go/koinos/v1"
	"github.com/odysseia-greek/makedonia/filippos/hetairoi"
	v1 "github.com/odysseia-greek/makedonia/ptolemaios/gen/go/v1"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (e *ExtendedServiceImpl) Health(ctx context.Context, request *emptypb.Empty) (*koinos.HealthResponse, error) {
	return &koinos.HealthResponse{
		Healthy:        true,
		Time:           time.Now().String(),
		DatabaseHealth: nil,
		Version:        e.Version,
	}, nil
}
func (e *ExtendedServiceImpl) Search(ctx context.Context, request *v1.ExtendedSearch) (*v1.ExtendedSearchResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	var requestId string
	if ok {
		headerValue := md.Get(service.HeaderKey)
		if len(headerValue) > 0 {
			requestId = headerValue[0]
		}
	}

	analyseResult := &v1.ExtendedSearchResponse{}

	cacheItem, _ := e.Archytas.Read(request.Word)
	if cacheItem != nil {
		err := json.Unmarshal(cacheItem, &analyseResult)
		if err != nil {
			return nil, err
		}

		go hetairoi.CacheSpan(string(cacheItem), request.Word, ctx)
		return analyseResult, nil
	}

	herodotosSpan := &pb.ParabasisRequest{
		RequestType: &pb.ParabasisRequest_Span{Span: &pb.SpanRequest{
			Action: "analyseText",
			Status: fmt.Sprintf("querying Herodotos for word: %s", request.Word),
		}},
	}

	hetairoi.ServiceToServiceSpan(herodotosSpan, ctx)

	startTime := time.Now()
	r := models.AnalyzeTextRequest{Rootword: request.Word}
	jsonBody, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	foundInText, err := e.Client.Herodotos().Analyze(jsonBody, requestId)
	endTime := time.Since(startTime)

	if foundInText != nil {
		var source models.AnalyzeTextResponse
		defer foundInText.Body.Close()
		err = json.NewDecoder(foundInText.Body).Decode(&source)
		if err != nil {
			logging.Error(fmt.Sprintf("error while decoding: %s", err.Error()))
		}

		herodotosSpan = &pb.ParabasisRequest{
			RequestType: &pb.ParabasisRequest_Span{Span: &pb.SpanRequest{
				Action: "analyseText",
				Took:   fmt.Sprintf("%v", endTime),
				Status: fmt.Sprintf("querying Herodotos returned: %d", foundInText.StatusCode),
			}},
		}
		hetairoi.ServiceToServiceSpan(herodotosSpan, ctx)

		analyseResult.FoundInText = &v1.AnalyzeTextResponse{
			Rootword:     source.Rootword,
			PartOfSpeech: source.PartOfSpeech,
			Conjugations: []*v1.Conjugations{},
			Texts:        []*v1.AnalyzeResult{},
		}

		for _, text := range source.Results {
			parsedText := &v1.AnalyzeResult{
				ReferenceLink: text.ReferenceLink,
				Author:        text.Author,
				Book:          text.Book,
				Reference:     text.Reference,
				Text: &v1.Rhema{
					Greek:        text.Text.Greek,
					Translations: text.Text.Translations,
					Section:      text.Text.Section,
				},
			}
			analyseResult.FoundInText.Texts = append(analyseResult.FoundInText.Texts, parsedText)
		}

		for _, conjugation := range source.Conjugations {
			parsedConjugation := &v1.Conjugations{
				Word: conjugation.Word,
				Rule: conjugation.Rule,
			}

			analyseResult.FoundInText.Conjugations = append(analyseResult.FoundInText.Conjugations, parsedConjugation)
		}
	}

	analyseResultJson, _ := json.Marshal(analyseResult)

	standardDuration := time.Minute * 10
	err = e.Archytas.SetWithTTL(request.Word, string(analyseResultJson), standardDuration)
	if err != nil {
		logging.Error(err.Error())
	}

	return analyseResult, nil
}

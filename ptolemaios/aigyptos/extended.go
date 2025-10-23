package aigyptos

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/agora/plato/models"
	"github.com/odysseia-greek/agora/plato/service"
	pb "github.com/odysseia-greek/attike/aristophanes/proto"
	koinos "github.com/odysseia-greek/makedonia/filippos/gen/go/koinos/v1"
	v1 "github.com/odysseia-greek/makedonia/ptolemaios/gen/go/v1"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (e *ExtendedServiceImpl) Health(ctx context.Context, request *emptypb.Empty) (*koinos.HealthResponse, error) {
	elasticHealth := e.Elastic.Health().Info()
	dbHealth := &koinos.DatabaseHealth{
		Healthy:       elasticHealth.Healthy,
		ClusterName:   elasticHealth.ClusterName,
		ServerName:    elasticHealth.ServerName,
		ServerVersion: elasticHealth.ServerVersion,
	}

	return &koinos.HealthResponse{
		Healthy:        true,
		Time:           time.Now().String(),
		DatabaseHealth: dbHealth,
		Version:        e.Version,
	}, nil
}
func (e *ExtendedServiceImpl) Search(ctx context.Context, request *v1.ExtendedSearch) (*v1.AnalyzeTextResponse, error) {

	md, ok := metadata.FromIncomingContext(ctx)
	var requestId string
	if ok {
		headerValue := md.Get(service.HeaderKey)
		if len(headerValue) > 0 {
			requestId = headerValue[0]
		}
	}

	splitID := strings.Split(requestId, "+")

	traceCall := false
	var traceID, spanID string

	if len(splitID) >= 3 {
		traceCall = splitID[2] == "1"
	}

	if len(splitID) >= 1 {
		traceID = splitID[0]
	}
	if len(splitID) >= 2 {
		spanID = splitID[1]
	}

	if traceCall {
		herodotosSpan := &pb.ParabasisRequest{
			TraceId:      traceID,
			ParentSpanId: spanID,
			SpanId:       spanID,
			RequestType: &pb.ParabasisRequest_Span{Span: &pb.SpanRequest{
				Action: "analyseText",
				Status: fmt.Sprintf("querying Herodotos for word: %s", request.Word),
			}},
		}

		err := e.Streamer.Send(herodotosSpan)
		if err != nil {
			logging.Error(fmt.Sprintf("error returned from tracer: %s", err.Error()))
		}
	}

	startTime := time.Now()
	r := models.AnalyzeTextRequest{Rootword: request.Word}
	jsonBody, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	foundInText, err := e.Client.Herodotos().Analyze(jsonBody, requestId)
	endTime := time.Since(startTime)
	var analyseResult *v1.AnalyzeTextResponse

	if foundInText != nil {
		var source models.AnalyzeTextResponse
		defer foundInText.Body.Close()
		err = json.NewDecoder(foundInText.Body).Decode(&source)
		if err != nil {
			logging.Error(fmt.Sprintf("error while decoding: %s", err.Error()))
		}

		if traceCall {
			herodotosSpan := &pb.ParabasisRequest{
				TraceId:      traceID,
				ParentSpanId: spanID,
				SpanId:       spanID,
				RequestType: &pb.ParabasisRequest_Span{Span: &pb.SpanRequest{
					Action: "analyseText",
					Took:   fmt.Sprintf("%v", endTime),
					Status: fmt.Sprintf("querying Herodotos returned: %d", foundInText.StatusCode),
				}},
			}

			err := e.Streamer.Send(herodotosSpan)
			if err != nil {
				logging.Error(fmt.Sprintf("error returned from tracer: %s", err.Error()))
			}
		}

		analyseResult = &v1.AnalyzeTextResponse{
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
			analyseResult.Texts = append(analyseResult.Texts, parsedText)
		}

		for _, conjugation := range source.Conjugations {
			parsedConjugation := &v1.Conjugations{
				Word: conjugation.Word,
				Rule: conjugation.Rule,
			}

			analyseResult.Conjugations = append(analyseResult.Conjugations, parsedConjugation)
		}
	}

	return analyseResult, nil
}

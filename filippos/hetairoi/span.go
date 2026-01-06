package hetairoi

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/agora/plato/service"
	"github.com/odysseia-greek/attike/aristophanes/comedy"
	arv1 "github.com/odysseia-greek/attike/aristophanes/gen/go/v1"
	"google.golang.org/grpc/metadata"
)

func extractRequestIds(ctx context.Context) (string, string, bool) {
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

	return traceID, spanID, traceCall
}
func DatabaseSpan(query map[string]interface{}, hits, timeTook int64, ctx context.Context) {
	traceID, spanID, traceCall := extractRequestIds(ctx)

	if !traceCall {
		return
	}

	parsedQuery, _ := json.Marshal(query)

	dataBaseSpan := &arv1.ObserveRequest{
		TraceId:      traceID,
		ParentSpanId: spanID,
		SpanId:       spanID,
		Kind: &arv1.ObserveRequest_DbSpan{DbSpan: &arv1.ObserveDbSpan{
			Action: "search",
			Query:  string(parsedQuery),
			Hits:   hits,
			TookMs: timeTook,
		}},
	}

	err := streamer.Send(dataBaseSpan)
	if err != nil {
		logging.Error(fmt.Sprintf("error returned from tracer: %s", err.Error()))
	}
}

func CacheSpan(response string, sessionId string, ctx context.Context) {
	traceID, spanID, traceCall := extractRequestIds(ctx)

	if !traceCall {
		return
	}

	span := &arv1.ObserveRequest{
		TraceId:      traceID,
		ParentSpanId: spanID,
		SpanId:       spanID,
		Kind: &arv1.ObserveRequest_Action{Action: &arv1.ObserveAction{
			Action: fmt.Sprintf("taken from cache with key: %s", sessionId),
			Status: response,
		}},
	}

	err := streamer.Send(span)
	if err != nil {
		logging.Error(fmt.Sprintf("error returned from tracer: %s", err.Error()))
	}
}

func ServiceToServiceSpan(span *arv1.ObserveRequest, ctx context.Context) {
	traceID, spanID, traceCall := extractRequestIds(ctx)

	if !traceCall {
		return
	}

	span.TraceId = traceID
	span.SpanId = comedy.GenerateSpanID()
	span.ParentSpanId = spanID

	err := streamer.Send(span)
	if err != nil {
		logging.Error(fmt.Sprintf("error returned from tracer: %s", err.Error()))
	}

}

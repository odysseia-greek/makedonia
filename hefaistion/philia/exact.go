package philia

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/odysseia-greek/agora/aristoteles/models"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/agora/plato/transform"
	koinos "github.com/odysseia-greek/makedonia/filippos/gen/go/koinos/v1"
	"github.com/odysseia-greek/makedonia/filippos/hetairoi"
	v1 "github.com/odysseia-greek/makedonia/hefaistion/gen/go/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (e *ExactServiceImpl) Health(ctx context.Context, request *emptypb.Empty) (*koinos.HealthResponse, error) {
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

func (e *ExactServiceImpl) Search(ctx context.Context, request *koinos.SearchQuery) (*v1.SearchResponse, error) {
	baseWord, strippedWord := extractBaseWord(request.Word)

	var language string
	switch request.Language {
	case koinos.Language_LANG_GREEK:
		language = "greek"
	case koinos.Language_LANG_ENGLISH:
		language = "english"
	case koinos.Language_LANG_DUTCH:
		language = "dutch"
	default:
		return nil, fmt.Errorf("unsupported language: %v", request.Language)
	}

	if request.NumberOfResults == 0 {
		request.NumberOfResults = 5
	}

	elasticResponse, hitsTotal, err := e.queryElastic(ctx, baseWord, language, false, request.NumberOfResults)
	if err != nil {
		return nil, err
	}

	if len(elasticResponse.Hits.Hits) == 0 {
		logging.Debug("no hits found trying with a word without diacretics")
		elasticResponse, hitsTotal, err = e.queryElastic(ctx, strippedWord, language, true, request.NumberOfResults)
		if err != nil {
			return nil, err
		}
	}

	results := make([]*koinos.Lemma, 0, len(elasticResponse.Hits.Hits))
	for _, hit := range elasticResponse.Hits.Hits {
		source, _ := json.Marshal(hit.Source)
		var src hetairoi.LemmaSource
		if err := json.Unmarshal(source, &src); err != nil {
			b, _ := json.Marshal(hit.Source)
			if err2 := json.Unmarshal(b, &src); err2 != nil {
				return nil, fmt.Errorf("decode _source: %w", err2)
			}
		}

		results = append(results, hetairoi.LemmaFromSource(src))
	}

	resp := &v1.SearchResponse{
		Results:  results,
		PageInfo: &koinos.PageInfo{Page: 1, Size: int32(len(results)), Total: int32(hitsTotal)},
	}
	return resp, nil
}

func (e *ExactServiceImpl) queryElastic(ctx context.Context, word, language string, normalized bool, results int32) (*models.Response, int64, error) {
	var query map[string]interface{}

	if normalized {
		query = map[string]interface{}{
			"query": map[string]interface{}{
				"match": map[string]interface{}{
					"normalized": word,
				},
			},
		}
	} else {
		query = map[string]interface{}{
			"query": map[string]interface{}{
				"bool": map[string]interface{}{
					"should": []interface{}{
						map[string]interface{}{
							"prefix": map[string]interface{}{
								fmt.Sprintf("%s.keyword", language): fmt.Sprintf("%s,", word),
							},
						},
						map[string]interface{}{
							"term": map[string]interface{}{
								fmt.Sprintf("%s.keyword", language): word,
							},
						},
					},
				},
			},
			"size": results,
		}
	}

	logging.Debug(fmt.Sprintf("%v", query))

	elasticResponse, err := e.Elastic.Query().Match(e.Index, query)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying elastic: %w", err)
	}

	hitsTotal := int64(0)
	if elasticResponse.Hits.Hits != nil {
		hitsTotal = elasticResponse.Hits.Total.Value
	}
	go hetairoi.DatabaseSpan(query, hitsTotal, elasticResponse.Took, ctx)

	return elasticResponse, hitsTotal, nil
}

func extractBaseWord(queryWord string) (string, string) {
	splitWord := strings.Split(queryWord, " ")

	// Known Greek pronouns
	greekPronouns := map[string]bool{"η": true, "ο": true, "το": true}

	// Function to clean punctuation from a word
	cleanWord := func(word string) string {
		return strings.Trim(word, ",.!?-") // Add any other punctuation as needed
	}

	// Iterate through the words
	for _, word := range splitWord {
		cleanedWord := cleanWord(word)

		if strings.HasPrefix(cleanedWord, "-") {
			// Skip words starting with "-"
			continue
		}

		if _, isPronoun := greekPronouns[cleanedWord]; !isPronoun {
			// If the word is not a pronoun, it's likely the correct word
			// Normalize and split the input
			strippedWord := transform.RemoveAccents(strings.ToLower(cleanedWord))
			return cleanedWord, strippedWord
		}
	}

	return queryWord, transform.RemoveAccents(strings.ToLower(queryWord))
}

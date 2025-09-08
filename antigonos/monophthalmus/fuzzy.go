package monophthalmus

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/odysseia-greek/agora/plato/transform"
	v1 "github.com/odysseia-greek/makedonia/antigonos/gen/go/v1"
	koinos "github.com/odysseia-greek/makedonia/filippos/gen/go/koinos/v1"
	"github.com/odysseia-greek/makedonia/filippos/hetairoi"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (f *FuzzyServiceImpl) Health(ctx context.Context, request *emptypb.Empty) (*koinos.HealthResponse, error) {
	elasticHealth := f.Elastic.Health().Info()
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
		Version:        f.Version,
	}, nil
}

func (f *FuzzyServiceImpl) Search(ctx context.Context, request *koinos.SearchQuery) (*v1.SearchResponse, error) {
	baseWord := extractBaseWord(request.Word)

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"fuzzy": map[string]interface{}{
				request.Language.String(): map[string]interface{}{
					"value":     baseWord,
					"fuzziness": 2,
				},
			},
		},
		"size": 5,
	}

	elasticResponse, err := f.Elastic.Query().Match(f.Index, query)
	if err != nil {
		return nil, fmt.Errorf("error querying elastic: %w", err)
	}

	hitsTotal := int64(0)
	if elasticResponse.Hits.Hits != nil {
		hitsTotal = elasticResponse.Hits.Total.Value
	}
	go hetairoi.DatabaseSpan(query, hitsTotal, elasticResponse.Took, ctx)

	hits := elasticResponse.Hits.Hits

	if len(hits) == 0 {
		noResponseQuery := map[string]interface{}{
			"query": map[string]interface{}{
				"fuzzy": map[string]interface{}{
					request.Language.String(): map[string]interface{}{
						"value":     request.Word,
						"fuzziness": 2,
					},
				},
			},
			"size": 5,
		}

		noResponseResponse, err := f.Elastic.Query().Match(f.Index, noResponseQuery)
		if err != nil {
			return nil, err
		}

		hitsTotal = int64(0)
		if elasticResponse.Hits.Hits != nil {
			hitsTotal = elasticResponse.Hits.Total.Value
		}

		go hetairoi.DatabaseSpan(query, hitsTotal, noResponseResponse.Took, ctx)
		hits = noResponseResponse.Hits.Hits
	}

	results := make([]*koinos.Lemma, 0, len(hits))
	for _, hit := range hits {
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

func extractBaseWord(queryWord string) string {
	// Normalize and split the input
	strippedWord := transform.RemoveAccents(strings.ToLower(queryWord))
	splitWord := strings.Split(strippedWord, " ")

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
			return cleanedWord
		}
	}

	return queryWord
}

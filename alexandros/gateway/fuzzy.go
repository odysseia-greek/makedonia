package gateway

import (
	"github.com/odysseia-greek/makedonia/alexandros/graph/model"
	antigonosv1 "github.com/odysseia-greek/makedonia/antigonos/gen/go/v1"
	"github.com/odysseia-greek/makedonia/antigonos/monophthalmus"
	koinos "github.com/odysseia-greek/makedonia/filippos/gen/go/koinos/v1"
)

func (a *AlexandrosHandler) Fuzzy(request *koinos.SearchQuery, requestID, sessionId string) (*model.SearchResponse, error) {
	fuzzyClientCtx, cancel := a.createRequestHeader(requestID, sessionId)
	defer cancel()

	var grpcResponse *antigonosv1.SearchResponse

	err := a.FuzzyClient.CallWithReconnect(func(client *monophthalmus.FuzzyClient) error {
		var innerErr error
		grpcResponse, innerErr = client.Search(fuzzyClientCtx, request)
		return innerErr
	})
	if err != nil {
		return nil, err
	}

	var lemmas []*model.Lemma

	for _, result := range grpcResponse.Results {
		glosses := make([]*model.LocalizedGloss, len(result.QuickGlosses))

		for _, gloss := range result.QuickGlosses {
			glosses = append(glosses, &model.LocalizedGloss{
				Language: gloss.Language,
				Gloss:    gloss.Gloss,
			})
		}

		definitions := make([]*model.Definition, len(result.Definitions))

		for _, definition := range result.Definitions {
			def := &model.Definition{
				Grade:    definition.Grade,
				Meanings: nil,
			}
			for _, meaning := range definition.Meanings {
				def.Meanings = append(def.Meanings, &model.Meaning{
					Language:   meaning.Language,
					Definition: meaning.Definition,
					Notes:      meaning.Notes,
					Example:    &meaning.Example,
				})
			}

			definitions = append(definitions, def)
		}

		modernConnections := make([]*model.ModernConnection, len(result.ModernConnections))
		for _, mc := range result.ModernConnections {
			modernConnections = append(modernConnections, &model.ModernConnection{
				Term: mc.Term,
				Note: &mc.Note,
			})
		}

		lemma := &model.Lemma{
			ID:                &result.Id,
			Headword:          result.Headword,
			Normalized:        &result.Normalized,
			LinkedWord:        &result.LinkedWord,
			PartOfSpeech:      &result.PartOfSpeech,
			Article:           &result.Article,
			Gender:            &result.Gender,
			Noun:              nil,
			Verb:              nil,
			QuickGlosses:      glosses,
			Definitions:       definitions,
			ModernConnections: modernConnections,
		}

		if result.Noun != nil {
			lemma.Noun = &model.NounInfo{
				Declension: &result.Noun.Declension,
				Genitive:   &result.Noun.Genitive,
			}
		} else if result.Verb != nil {
			lemma.Verb = &model.VerbInfo{
				PrincipalParts: result.Verb.PrincipalParts,
			}
		}

		lemmas = append(lemmas, lemma)
	}
	resp := &model.SearchResponse{
		Results: lemmas,
		PageInfo: &model.PageInfo{
			Page:  grpcResponse.PageInfo.Page,
			Size:  grpcResponse.PageInfo.Size,
			Total: grpcResponse.PageInfo.Total,
		},
	}

	return resp, nil
}

package gateway

import (
	"github.com/odysseia-greek/makedonia/alexandros/graph/model"
	pbe "github.com/odysseia-greek/makedonia/eukleides/proto"
	koinos "github.com/odysseia-greek/makedonia/filippos/gen/go/koinos/v1"
	hefaistionv1 "github.com/odysseia-greek/makedonia/hefaistion/gen/go/v1"
	"github.com/odysseia-greek/makedonia/hefaistion/philia"
)

func (a *AlexandrosHandler) Exact(request *koinos.SearchQuery, requestID, sessionId string) (*model.SearchResponse, error) {
	fuzzyClientCtx, cancel := a.createRequestHeader(requestID, sessionId)
	defer cancel()

	eukleidesUpdate := pbe.CountCreationRequest{
		Word:        request.Word,
		ServiceName: "exact",
		SearchType:  "exact",
		SessionId:   sessionId,
	}

	go a.pushToEukleides(&eukleidesUpdate)

	var grpcResponse *hefaistionv1.SearchResponse

	err := a.ExactClient.CallWithReconnect(func(client *philia.ExactClient) error {
		var innerErr error
		grpcResponse, innerErr = client.Search(fuzzyClientCtx, request)
		return innerErr
	})
	if err != nil {
		return nil, err
	}

	var lemmas []*model.Lemma

	for _, result := range grpcResponse.Results {
		// QuickGlosses
		glosses := make([]*model.LocalizedGloss, 0, len(result.QuickGlosses))
		for _, gloss := range result.QuickGlosses {
			glosses = append(glosses, &model.LocalizedGloss{
				Language: gloss.Language,
				Gloss:    gloss.Gloss,
			})
		}

		// Definitions
		definitions := make([]*model.Definition, 0, len(result.Definitions))
		for _, definition := range result.Definitions {
			def := &model.Definition{
				Grade:    definition.Grade,
				Meanings: make([]*model.Meaning, 0, len(definition.Meanings)),
			}
			for _, meaning := range definition.Meanings {
				// Ensure notes is a non-nil slice for [String!]!
				notes := meaning.Notes
				if notes == nil {
					notes = []string{}
				}
				// Example is optional (String), nil is fine.
				def.Meanings = append(def.Meanings, &model.Meaning{
					Language:   meaning.Language,
					Definition: meaning.Definition,
					Notes:      notes,
					Example:    &meaning.Example,
				})
			}
			definitions = append(definitions, def)
		}

		// ModernConnections
		modernConnections := make([]*model.ModernConnection, 0, len(result.ModernConnections))
		for _, mc := range result.ModernConnections {
			note := mc.Note // optional
			modernConnections = append(modernConnections, &model.ModernConnection{
				Term: mc.Term,
				Note: &note,
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
			QuickGlosses:      glosses,     // non-nil, possibly empty
			Definitions:       definitions, // non-nil, possibly empty
			ModernConnections: modernConnections,
		}

		if result.Noun != nil {
			lemma.Noun = &model.NounInfo{
				Declension: &result.Noun.Declension,
				Genitive:   &result.Noun.Genitive,
			}
		}
		if result.Verb != nil {
			// Ensure principal parts is non-nil for [String!]! in your SDL
			parts := result.Verb.PrincipalParts
			if parts == nil {
				parts = []string{}
			}
			lemma.Verb = &model.VerbInfo{
				PrincipalParts: parts,
			}
		}

		lemmas = append(lemmas, lemma)
	}

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

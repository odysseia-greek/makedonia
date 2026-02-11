package hetairoi

import (
	koinos "github.com/odysseia-greek/makedonia/filippos/gen/go/koinos/v1"
)

func LemmaFromSource(s LemmaSource) *koinos.Lemma {
	quick := make([]*koinos.LocalizedGloss, 0, 2)
	if s.English != "" {
		quick = append(quick, &koinos.LocalizedGloss{Language: "en", Gloss: s.English})
	}
	if s.Dutch != "" {
		quick = append(quick, &koinos.LocalizedGloss{Language: "nl", Gloss: s.Dutch})
	}

	var noun *koinos.NounInfo
	if s.Noun != nil {
		noun = &koinos.NounInfo{
			Declension: s.Noun.Declension,
			Genitive:   s.Noun.Genitive,
		}
	}

	var verb *koinos.VerbInfo
	if s.Verb != nil {
		verb = &koinos.VerbInfo{
			PrincipalParts: s.Verb.PrincipalParts,
		}
	}

	defs := make([]*koinos.Definition, 0, len(s.Definitions))
	for _, d := range s.Definitions {
		ms := make([]*koinos.Meaning, 0, len(d.Meanings))
		for _, m := range d.Meanings {
			ms = append(ms, &koinos.Meaning{
				Language:   m.Language,
				Definition: m.Definition,
				Notes:      m.Notes,
				Example:    m.Example,
			})
		}
		defs = append(defs, &koinos.Definition{
			Grade:    int32(d.Grade),
			Meanings: ms,
		})
	}

	mconns := make([]*koinos.ModernConnection, 0, len(s.ModernConns))
	for _, mc := range s.ModernConns {
		mconns = append(mconns, &koinos.ModernConnection{Term: mc.Term, Note: mc.Note})
	}

	return &koinos.Lemma{
		Id:                s.ID,
		Headword:          s.Greek,
		Normalized:        s.Normalized,
		LinkedWord:        s.LinkedWord,
		PartOfSpeech:      s.PartOfSpeech,
		Article:           s.Article,
		Gender:            s.Gender,
		Noun:              noun,
		Verb:              verb,
		QuickGlosses:      quick,
		Definitions:       defs,
		ModernConnections: mconns,
	}
}

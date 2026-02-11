package graph

import (
	"github.com/odysseia-greek/makedonia/alexandros/graph/model"
	koinos "github.com/odysseia-greek/makedonia/filippos/gen/go/koinos/v1"
)

func parseLanguage(inputLanguage *model.Language) koinos.Language {
	var language koinos.Language
	switch *inputLanguage {
	case model.LanguageLangGreek:
		language = koinos.Language_LANG_GREEK
	case model.LanguageLangEnglish:
		language = koinos.Language_LANG_ENGLISH
	case model.LanguageLangDutch:
		language = koinos.Language_LANG_DUTCH
	default:
		language = koinos.Language_LANG_GREEK
	}

	return language
}

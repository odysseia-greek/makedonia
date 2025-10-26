package gateway

import (
	"fmt"
	"net/http"

	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/agora/plato/middleware"
	"github.com/odysseia-greek/agora/plato/models"
	koinos "github.com/odysseia-greek/makedonia/filippos/gen/go/koinos/v1"
)

func (a *AlexandrosHandler) SearchWord(w http.ResponseWriter, req *http.Request) {
	var requestId string
	fromContext := req.Context().Value(config.DefaultTracingName)
	if fromContext == nil {
		requestId = req.Header.Get(config.HeaderKey)
	} else {
		requestId = fromContext.(string)
	}

	queryWord := req.URL.Query().Get("word")

	request := &koinos.SearchQuery{
		Word:     queryWord,
		Language: koinos.Language_LANG_GREEK,
	}
	response, err := a.Exact(request, requestId, "internal")
	if err != nil {
		middleware.ResponseWithCustomCode(w, http.StatusInternalServerError, err)
		return
	}

	var results models.ExtendedResponse

	logging.Debug(fmt.Sprintf("responseLengtht: %d", len(response.Results)))
	for _, resp := range response.Results {
		hit := models.Hit{
			Hit: models.Meros{
				Greek: resp.Headword,
			},
			FoundInText: nil,
		}

		for _, gloss := range resp.QuickGlosses {
			if gloss.Language == "en" {
				hit.Hit.English = gloss.Gloss
			}

			if gloss.Language == "nl" {
				hit.Hit.Dutch = gloss.Gloss

			}
		}

		results.Hits = append(results.Hits, hit)
	}
	middleware.ResponseWithCustomCode(w, http.StatusOK, results)
}

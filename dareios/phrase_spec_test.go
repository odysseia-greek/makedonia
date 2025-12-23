package main

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	gq "github.com/odysseia-greek/makedonia/dareios/internal/graphql"
)

type phraseResponse struct {
	Phrase struct {
		Results []struct {
			Headword     string `json:"headword"`
			PartOfSpeech string `json:"partOfSpeech"`
			Normalized   string `json:"normalized"`
			QuickGlosses []struct {
				Language string `json:"language"`
				Gloss    string `json:"gloss"`
			} `json:"quickGlosses"`
			LinkedWord string `json:"linkedWord"`
		} `json:"results"`
		PageInfo struct {
			Page  int `json:"page"`
			Total int `json:"total"`
		} `json:"pageInfo"`
	} `json:"phrase"`
}

var _ = Describe("phrase query", func() {
	It("returns paged phrase results for a simple word", func(ctx context.Context) {
		c, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		const q = `query($input: SearchQueryInput!) { phrase(input: $input) {
		results {
			headword
			partOfSpeech
			normalized
			quickGlosses{
				language
				gloss
			}
			verb{
				principalParts
			}
			modernConnections{
				term
				note
			}
			definitions{
				grade
				meanings{
					definition
					language
				}
			}
			linkedWord
		}
		pageInfo{
			page
			total
		}
	} 
}`
		vars := map[string]any{
			"input": map[string]any{
				"word": "λόγος",
				"size": 2,
			},
		}
		var resp phraseResponse
		err := gq.Execute(c, baseURL, q, vars, &resp)
		Expect(err).NotTo(HaveOccurred())

		f := resp.Phrase
		Expect(f.PageInfo.Page).To(BeNumerically(">=", 1))
		Expect(f.PageInfo.Total).To(BeNumerically(">=", 0))
		// We allow empty results if dataset is empty, but structure must be valid
		if len(f.Results) > 0 {
			r := f.Results[0]
			Expect(r.Headword).NotTo(BeEmpty())
		}
	}, SpecTimeout(20*time.Second))
})

package main

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	gq "github.com/odysseia-greek/makedonia/dareios/internal/graphql"
)

type exactResponse struct {
	Exact struct {
		Results []struct {
			Headword     string `json:"headword"`
			PartOfSpeech string `json:"partOfSpeech"`
			Normalized   string `json:"normalized"`
			QuickGlosses []struct {
				Language string `json:"language"`
				Gloss    string `json:"gloss"`
			} `json:"quickGlosses"`
			Noun struct {
				Declension string `json:"declension"`
				Genitive   string `json:"genitive"`
			} `json:"noun"`
			Verb struct {
				PrincipalParts []string `json:"principalParts"`
			} `json:"verb"`
			ModernConnections []struct {
				Term string `json:"term"`
				Note string `json:"note"`
			} `json:"modernConnections"`
			Definitions []struct {
				Grade    int `json:"grade"`
				Meanings []struct {
					Definition string `json:"definition"`
					Language   string `json:"language"`
				} `json:"meanings"`
			} `json:"definitions"`
			LinkedWord string `json:"linkedWord"`
		} `json:"results"`
		PageInfo struct {
			Page  int `json:"page"`
			Total int `json:"total"`
		} `json:"pageInfo"`
		SimilarWords []struct {
			English  string `json:"english"`
			Greek    string `json:"greek"`
			Original string `json:"original"`
		} `json:"similarWords"`
		FoundInText struct {
			Rootword string `json:"rootword"`
			Texts    []struct {
				Author        string `json:"author"`
				Book          string `json:"book"`
				Reference     string `json:"reference"`
				ReferenceLink string `json:"referenceLink"`
				Text          struct {
					Greek        string   `json:"greek"`
					Section      string   `json:"section"`
					Translations []string `json:"translations"`
				} `json:"text"`
			} `json:"texts"`
			Conjugations []struct {
				Rule string `json:"rule"`
				Word string `json:"word"`
			} `json:"conjugations"`
		} `json:"foundInText"`
	} `json:"exact"`
}

const q = `query($input: ExpandableSearchQueryInput!) { exact(input: $input) {
		results {
			headword
			partOfSpeech
			normalized
			quickGlosses{
				language
				gloss
			}
			noun{
				declension
				genitive
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
		similarWords{
			english
			greek
			original
		}
		foundInText{
			rootword
			texts{
				author
				book
				reference
				referenceLink
				text{
					greek
					section
					translations
				}
			}
			conjugations{
				rule
				word
			}
		}
	}
}`

var _ = Describe("exact query", func() {
	It("returns paged exact results for a simple word without expand", func(ctx context.Context) {
		c, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		vars := map[string]any{
			"input": map[string]any{
				"word":   "λόγος",
				"expand": false,
				"size":   1,
			},
		}
		var resp exactResponse
		err := gq.Execute(c, baseURL, q, vars, &resp)
		Expect(err).NotTo(HaveOccurred())

		f := resp.Exact
		Expect(f.PageInfo.Page).To(BeNumerically(">=", 1))
		Expect(f.PageInfo.Total).To(BeNumerically(">=", 0))
		// We allow empty results if dataset is empty, but structure must be valid
		if len(f.Results) > 0 {
			r := f.Results[0]
			Expect(r.Headword).NotTo(BeEmpty())
		}
	}, SpecTimeout(20*time.Second))

	It("an expanded result has texts included", func(ctx context.Context) {
		c, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		vars := map[string]any{
			"input": map[string]any{
				"word":   "λόγος",
				"expand": true,
				"size":   1,
			},
		}
		var resp exactResponse
		err := gq.Execute(c, baseURL, q, vars, &resp)
		Expect(err).NotTo(HaveOccurred())

		f := resp.Exact
		Expect(len(f.FoundInText.Texts)).To(BeNumerically(">=", 1))
		Expect(f.FoundInText.Texts[0].Text.Greek).NotTo(BeEmpty())
		Expect(len(f.SimilarWords)).To(BeNumerically(">=", 1))
	}, SpecTimeout(20*time.Second))

	It("the exact result for the word λόγος also includes extra fields such as connections and a noun part", func(ctx context.Context) {
		c, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		vars := map[string]any{
			"input": map[string]any{
				"word":   "λόγος",
				"expand": false,
				"size":   1,
			},
		}
		var resp exactResponse
		err := gq.Execute(c, baseURL, q, vars, &resp)
		Expect(err).NotTo(HaveOccurred())

		f := resp.Exact
		Expect(len(f.Results[0].ModernConnections)).To(BeNumerically(">=", 1))
		Expect(f.Results[0].Noun.Declension).To(Equal("second"))
		Expect(len(f.Results[0].Definitions)).To(BeNumerically(">=", 1))
		Expect(f.Results[0].Definitions[0].Meanings[0].Definition).NotTo(BeEmpty())
		Expect(f.Results[0].Verb.PrincipalParts).To(BeEmpty())

	}, SpecTimeout(20*time.Second))

	It("the exact result for the word γίγνομαι also includes extra fields such as connections and a verb part", func(ctx context.Context) {
		c, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		vars := map[string]any{
			"input": map[string]any{
				"word":   "γίγνομαι",
				"expand": false,
				"size":   1,
			},
		}
		var resp exactResponse
		err := gq.Execute(c, baseURL, q, vars, &resp)
		Expect(err).NotTo(HaveOccurred())

		f := resp.Exact
		Expect(len(f.Results[0].ModernConnections)).To(BeNumerically(">=", 1))
		Expect(f.Results[0].Noun.Genitive).To(BeEmpty())
		Expect(len(f.Results[0].Definitions)).To(BeNumerically(">=", 1))
		Expect(f.Results[0].Definitions[0].Meanings[0].Definition).NotTo(BeEmpty())
		Expect(len(f.Results[0].Verb.PrincipalParts)).To(BeNumerically(">=", 1))

	}, SpecTimeout(20*time.Second))
})

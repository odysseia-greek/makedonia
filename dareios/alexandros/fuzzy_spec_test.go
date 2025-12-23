package alexandros

import (
    "context"
    "time"

    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"

    gq "github.com/odysseia-greek/makedonia/dareios/internal/graphql"
)

type fuzzyResponse struct {
    Fuzzy struct {
        Results []struct {
            Headword       string `json:"headword"`
            PartOfSpeech   string `json:"partOfSpeech"`
            Normalized     string `json:"normalized"`
            QuickGlosses   []struct {
                Language string `json:"language"`
                Gloss    string `json:"gloss"`
            } `json:"quickGlosses"`
            LinkedWord string `json:"linkedWord"`
        } `json:"results"`
        PageInfo struct {
            Page  int `json:"page"`
            Total int `json:"total"`
        } `json:"pageInfo"`
    } `json:"fuzzy"`
}

var _ = Describe("fuzzy query", func() {
    It("returns paged fuzzy results for a simple word", func() {
        c, cancel := context.WithTimeout(context.Background(), 15*time.Second)
        defer cancel()

        const q = `query($input: FuzzyInput!) { fuzzy(input: $input) { results { headword partOfSpeech normalized quickGlosses { language gloss } linkedWord } pageInfo { page total } } }`
        vars := map[string]any{
            "input": map[string]any{
                "word": "λογος",
                "size": 5,
            },
        }
        var resp fuzzyResponse
        err := gq.Execute(c, baseURL, q, vars, &resp)
        Expect(err).NotTo(HaveOccurred())

        f := resp.Fuzzy
        Expect(f.PageInfo.Page).To(BeNumerically(">=", 1))
        Expect(f.PageInfo.Total).To(BeNumerically(">=", 0))
        // We allow empty results if dataset is empty, but structure must be valid
        if len(f.Results) > 0 {
            r := f.Results[0]
            Expect(r.Headword).NotTo(BeEmpty())
        }
    }, SpecTimeout(20*time.Second))
})

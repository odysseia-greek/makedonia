package alexandros

import (
    "context"
    "time"

    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"

    gq "github.com/odysseia-greek/makedonia/dareios/internal/graphql"
)

type counterTopFiveResponse struct {
    CounterTopFive struct {
        TopFive []struct {
            LastUsed   string `json:"lastUsed"`
            ServiceName string `json:"serviceName"`
            Word       string `json:"word"`
            Count      int    `json:"count"`
        } `json:"topFive"`
    } `json:"counterTopFive"`
}

var _ = Describe("counterTopFive query", func() {
    It("returns top five counters with non-negative counts", func() {
        c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        defer cancel()

        const q = `query { counterTopFive { topFive { lastUsed serviceName word count } } }`
        var resp counterTopFiveResponse
        err := gq.Execute(c, baseURL, q, nil, &resp)
        Expect(err).NotTo(HaveOccurred())

        tf := resp.CounterTopFive.TopFive
        Expect(tf).NotTo(BeNil())
        for _, e := range tf {
            Expect(e.Word).NotTo(BeEmpty())
            Expect(e.ServiceName).NotTo(BeEmpty())
            Expect(e.Count).To(BeNumerically(">=", 0))
        }
    }, SpecTimeout(15*time.Second))
})

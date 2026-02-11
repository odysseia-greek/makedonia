package main

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	gq "github.com/odysseia-greek/makedonia/dareios/internal/graphql"
)

type healthResponse struct {
	Health struct {
		Time     string `json:"time"`
		Healthy  bool   `json:"healthy"`
		Version  string `json:"version"`
		Services []struct {
			Name         string `json:"name"`
			Version      string `json:"version"`
			Healthy      bool   `json:"healthy"`
			DatabaseInfo struct {
				Healthy       bool   `json:"healthy"`
				ServerName    string `json:"serverName"`
				ServerVersion string `json:"serverVersion"`
				ClusterName   string `json:"clusterName"`
			} `json:"databaseInfo"`
		} `json:"services"`
	} `json:"health"`
}

var _ = Describe("health query", func() {
	It("returns overall health and service statuses", func(ctx context.Context) {
		c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		const q = `query { health { time healthy version services { name version healthy databaseInfo { healthy serverName serverVersion clusterName } } } }`
		var resp healthResponse
		err := gq.Execute(c, baseURL, q, nil, &resp)
		Expect(err).NotTo(HaveOccurred())

		h := resp.Health
		Expect(h.Time).NotTo(BeEmpty())
		Expect(h.Healthy).To(BeTrue())
		Expect(h.Services).NotTo(BeNil())
		if len(h.Services) > 0 {
			s := h.Services[0]
			Expect(s.Name).NotTo(BeEmpty())
			Expect(s.DatabaseInfo.ServerVersion).NotTo(BeEmpty())
		}
	}, SpecTimeout(15*time.Second))
})

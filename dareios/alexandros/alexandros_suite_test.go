package alexandros

import (
    "os"
    "testing"

    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
)

// ALEXANDROS_URL can point to the GraphQL endpoint, e.g. http://localhost:8080/query
var baseURL string

func TestAlexandros(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "Alexandros GraphQL Suite")
}

var _ = BeforeSuite(func() {
    baseURL = os.Getenv("ALEXANDROS_URL")
    if baseURL == "" {
        baseURL = "http://localhost:8080/query"
    }
})

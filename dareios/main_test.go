package main

import (
	"fmt"
	"testing"

	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/logging"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	baseURL string
)

func TestDareios(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Dareios Integration Suite")
}

var _ = BeforeSuite(func() {
	logging.System(`
 ___     ____  ____     ___  ____  ___   _____
|   \   /    ||    \   /  _]|    |/   \ / ___/
|    \ |  o  ||  D  ) /  [_  |  ||     (   \_ 
|  D  ||     ||    / |    _] |  ||  O  |\__  |
|     ||  _  ||    \ |   [_  |  ||     |/  \ |
|     ||  |  ||  .  \|     | |  ||     |\    |
|_____||__|__||__|\_||_____||____|\___/  \___|
    `)
	logging.System("\"οἱ βασιλέως λόγοι πρὸς τοὺς Ἕλληνας οὐκ ἀληθεῖς εἰσίν.\"")
	logging.System("The king's words to the Greeks are not to be trusted.")

	// Get the URL from environment, defaulting to local if not set
	baseURL = config.StringFromEnv("ALEXANDROS_URL", "http://byzantium.odysseia-greek:8080/alexandros/graphql")

	logging.System(fmt.Sprintf("Integration tests targeting: %s", baseURL))
	logging.System("Starting up integration tests...")

})

//go:generate go run github.com/99designs/gqlgen generate
package graph

import (
	"github.com/odysseia-greek/makedonia/alexandros/gateway"
)

// Resolver struct for dependency injection
type Resolver struct {
	Handler *gateway.AlexandrosHandler
}

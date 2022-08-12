package clients

import (
	"testing"

	"atayalan.com/go-service-sdk/servicediscovery"
	"github.com/gofiber/fiber/v2"
)

var app *fiber.App = fiber.New(fiber.Config{
	DisableStartupMessage: true,
})

func setup() {
	servicediscovery.OnInit()
	err := OnInit()
	if err != nil {
		panic(err)
	}

	// TODO JWT validation via SDK security filter
	// not required for for single-tenant deployements, but
	// need it for multi-tenant deployments
	// FIXME later

	// User APIs
	app.Get("/usage/clients", GetClientsUsage)
	app.Get("/usage/clients/:clientId", GetClientsUsageByClientId)
	app.Get("/usage/dnns", GetDnnUsage)

}

func TestCreatePolicy(t *testing.T) {
	// this  function has to be called with setup
	setup()
	/// Todo
}

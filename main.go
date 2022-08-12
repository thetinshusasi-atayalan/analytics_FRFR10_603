package main

import (
	"context"

	"atayalan.com/go-service-sdk/logging"
	"atayalan.com/go-service-sdk/metrics"
	"atayalan.com/go-service-sdk/rest"
	"atayalan.com/go-service-sdk/server"
	"atayalan.com/go-service-sdk/servicediscovery"

	// "github.com/gofiber/fiber/v2/middleware/cors"
	"atayalan.com/analytics/clients"

	"atayalan.com/analytics/influxDB"
)

var counter metrics.Counter
var log = logging.CreateLogger("main")

type Analytics struct{}

func (Analytics) OnInit() error {
	var ctx = context.Background()

	servicediscovery.OnInit()

	config, err := servicediscovery.GetServiceConfig("influxDB")
	if err != nil {
		log.Errorf(ctx, "Failed to get analytics service config  %v", err)
		return err
	}
	log.Infof(ctx, "service config %v", config)

	client, err := influxDB.GetClient(ctx)
	if err != nil {
		log.Errorf(ctx, "Failed to connect to InfluxDB %v", err)
		return err
	}
	defer client.Close()

	isServerActive := influxDB.CheckPing(ctx, client)
	if !isServerActive {
		log.Errorf(ctx, "InfluxDB health check ping failed", err)
		return err
	}

	err = clients.OnInit()
	if err != nil {
		log.Fatalf(ctx, "Failed to create the clients bucket %v", err)
		return err
	}

	// rest.RegisterPostApi("/usage/clients", clients.CreateClientMeasurement)
	rest.RegisterGetApi("/usage/clients", clients.GetClientsUsage)
	rest.RegisterGetApi("/usage/clients/:clientId", clients.GetClientsUsageByClientId)
	rest.RegisterGetApi("/usage/dnns", clients.GetDnnUsage)

	// Service Config APIs

	return nil
}

//Stop do any clean up there
func (Analytics) Stop() {
}

func (Analytics) Start() error {
	return nil
}

//Status this is the health check api
func (Analytics) Status() bool {
	return true
}

func main() {
	rest.EnableCors = true
	service := Analytics{}
	server.RegisterService(service)
	server.Run(true)
}

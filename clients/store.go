package clients

import (
	"context"
	"errors"

	"atayalan.com/analytics/influxDB"
	"atayalan.com/go-service-sdk/logging"
)

const (

	// bucket name
	clientsDatabase  = "clients"
	measurementUsage = "usage"
)

var log = logging.CreateLogger("clients")

func OnInit() error {
	var err error
	var ctx = context.Background()
	dbClient, err := influxDB.GetClient(ctx)
	if err != nil {
		return err
	}
	defer dbClient.Close()

	isServerActive := influxDB.CheckPing(ctx, dbClient)
	if !isServerActive {
		log.Errorf(ctx, "InfluxDB health check failed", err)
		return errors.New("InfluxDB health check failed")
	}

	_, err = influxDB.CreateDatabaseWithName(ctx, dbClient, clientsDatabase)
	if err != nil {
		log.Errorf(ctx, "Failed to create the bucket", err)
		return err
	}

	return nil
}

package influxDB

import (
	"context"
	"fmt"
	"time"

	"atayalan.com/go-service-sdk/logging"
	"atayalan.com/go-service-sdk/security"

	"atayalan.com/go-service-sdk/servicediscovery"

	// "github.com/gofiber/fiber/v2/middleware/cors"
	"atayalan.com/analytics/utils"

	influxdb2 "github.com/influxdata/influxdb1-client/v2"
)

var log = logging.CreateLogger("influxDB")

const OrgNameAtayalan string = "Atayalan"

var pingTimeInSec = int64(15)
var SpaceString = " "

func GetClient(ctx context.Context) (influxdb2.Client, error) {

	config, err := servicediscovery.GetServiceConfig("influxDB")
	if err != nil {
		log.Errorf(ctx, "Failed to get analytics service config  %v", err)
		return nil, err
	}
	log.Infof(ctx, "service config %v", config.Port)

	// dbToken := os.Getenv("INFLUX_DB_TOKEN")
	// if err != nil {
	// 	log.Errorf(ctx, "Failed to get analytics INFLUX_DB_TOKEN from environment variables  %v", err)
	// 	return nil, err
	// }

	// log.Infof(ctx, "Token values =======", dbToken)

	// client := influxdb2.NewClient("http://"+config.Host+":"+utils.Int64ToStr(int64(config.Port)), dbToken)
	client, err := influxdb2.NewHTTPClient(influxdb2.HTTPConfig{
		Addr: "http://" + config.Host + ":" + utils.Int64ToStr(int64(config.Port)),
	})
	if err != nil {
		log.Errorf(ctx, "Failed to create influxdb client %v", err)
		return nil, err
	}

	return client, err
}

func CheckPing(ctx context.Context, client influxdb2.Client) bool {
	isServerActive := true
	_, _, err := client.Ping(time.Duration(pingTimeInSec) * time.Second)
	if err != nil {
		log.Errorf(ctx, "InfluxDB health check failed", err)
		isServerActive = false
		return isServerActive
	}
	return isServerActive

}

func CreateDatabaseWithName(ctx context.Context, client influxdb2.Client, databaseName string) (*influxdb2.Response, error) {
	// Todo: Add database retention policy code here
	query := influxdb2.NewQuery(fmt.Sprintf("CREATE DATABASE %s", databaseName), "", "")
	res, err := client.Query(query)
	if err != nil || res.Error() != nil {
		return nil, err
	}

	return res, nil

}

func WritePointToMeasurement(ctx context.Context, client influxdb2.Client, databaseName, measurementName string, tags map[string]string, fields map[string]interface{}) error {

	bp, err := influxdb2.NewBatchPoints(influxdb2.BatchPointsConfig{
		Database: databaseName,
	})
	if err != nil {
		log.Errorf(ctx, "Failed to create a batch points %v", err)
		return nil
	}
	// Adding TentantId in the measurement
	tags["tentantId"] = security.FromContextGetUserContext(ctx).GetTenantId()

	pt, err := influxdb2.NewPoint(measurementName, tags, fields, time.Now())
	if err != nil {
		log.Errorf(ctx, "Failed to create a point in influxdb %v", err)
		return nil
	}
	bp.AddPoint(pt)
	err = client.Write(bp)
	if err != nil {
		log.Errorf(ctx, "Failed to write the point %v", err)
		return nil
	}

	return nil
}

func Read(ctx context.Context, client influxdb2.Client, databaseName, cmd string) ([]influxdb2.Result, error) {
	q := influxdb2.Query{
		Command:  cmd,
		Database: databaseName,
	}
	res, err := client.Query(q)

	if err != nil || res.Error() != nil {
		log.Errorf(ctx, "Failed to query the database %s  res Err: %s", err, res.Error())
		return nil, err
	}

	return res.Results, nil
}

package clients

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"atayalan.com/analytics/influxDB"
	"atayalan.com/analytics/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/influxdata/influxdb1-client/models"
	influxdb2 "github.com/influxdata/influxdb1-client/v2"
)

func CreateClientMeasurement(c *fiber.Ctx) error {
	ctx := c.UserContext()
	var data Usage
	err := c.BodyParser(&data)
	if err != nil {
		log.Errorf(ctx, "Failed to create client measurement. Err: %s", err)
		return c.Status(fiber.StatusBadRequest).JSON(utils.GetErrorDetails("Failed to create client measurement. Err: %s", err))
	}

	errors := data.ValidateStruct()
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)

	}

	err = createClientMeasurementAsync(ctx, data)
	if err != nil {
		log.Errorf(ctx, "Failed to create client measurement. Err: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(errors)

	}
	return c.SendStatus(fiber.StatusOK)
}

func createClientMeasurementAsync(ctx context.Context, data Usage) error {

	clientCon, err := influxDB.GetClient(ctx)
	if err != nil {
		log.Errorf(ctx, "Failed to connect to InfluxDB %v", err)
		return err
	}
	defer clientCon.Close()
	tags := map[string]string{
		"clientId":   data.ClientId,
		"clientType": string(data.ClientType),
		"dnn":        data.Dnn,
	}

	fields := map[string]interface{}{
		"bytesToClient":   data.BytesToClient,
		"bytesFromClient": data.BytesFromClient,
		"ip":              data.IP,
		"upfIp":           data.UPFIP,
	}

	err = influxDB.WritePointToMeasurement(ctx, clientCon, clientsDatabase, measurementUsage, tags, fields)
	if err != nil {
		log.Errorf(ctx, "Failed to  write the measurement to influxdb %v", err)
		return err
	}

	return nil

}

func GetClientsUsage(c *fiber.Ctx) error {
	ctx := c.UserContext()
	timeRange := c.Query("timeRange")
	if timeRange == "" {
		log.Errorf(ctx, "Invalid Time range query paramater %v", timeRange)
		return c.Status(fiber.StatusBadRequest).JSON(fmt.Sprintf("Invalid Time Range. Err: %s", errors.New("Invalid Time Range value")))
	}
	_, err := strconv.ParseFloat(timeRange, 64)
	if err != nil {
		log.Errorf(ctx, "Invalid Time range query paramater %v", timeRange)
		return c.Status(fiber.StatusBadRequest).JSON(fmt.Sprintf("Invalid Time Range. Err: %s", errors.New("Invalid Time Range value")))
	}

	resp, err := getClientsUsage(ctx, timeRange, "", "", "")
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fmt.Sprintf("Failed to read from the influxDB. Err: %s", err))

	}

	tags := []string{
		"clientId",
		"clientType",
	}
	fields := []string{
		"bytesIn",
		"bytesOut",
	}

	data, _ := extractClientsDbResults(resp, tags, fields)

	return c.Status(fiber.StatusOK).JSON(data)

}

func GetClientsUsageByClientId(c *fiber.Ctx) error {
	ctx := c.UserContext()
	timeRange := c.Query("timeRange")
	if timeRange == "" {
		log.Errorf(ctx, "Invalid Time range query paramater %v", timeRange)
		return c.Status(fiber.StatusBadRequest).JSON(fmt.Sprintf("Invalid Time Range. Err: %s", errors.New("Invalid Time Range value")))
	}
	_, err := strconv.ParseFloat(timeRange, 64)
	if err != nil {
		log.Errorf(ctx, "Invalid Time range query paramater %v", timeRange)
		return c.Status(fiber.StatusBadRequest).JSON(fmt.Sprintf("Invalid Time Range. Err: %s", errors.New("Invalid Time Range value")))
	}

	clientId := c.Params("clientId")
	if clientId == "" {
		log.Errorf(ctx, "Invalid Client id %v", clientId)
		return c.Status(fiber.StatusBadRequest).JSON(fmt.Sprintf("Invalid Client id. Err: %s", errors.New("Invalid Client id value")))

	}

	resp, err := getClientsUsage(ctx, timeRange, "", "clientId", clientId)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fmt.Sprintf("Failed to read from the influxDB. Err: %s", err))

	}

	tags := []string{
		"clientId",
		"clientType",
	}
	fields := []string{
		"bytesIn",
		"bytesOut",
	}

	data := extractClientsByClientIdFromDbResults(resp, tags, fields)

	return c.Status(fiber.StatusOK).JSON(data)

}

func GetDnnUsage(c *fiber.Ctx) error {
	ctx := c.UserContext()
	timeRange := c.Query("timeRange")
	if timeRange == "" {
		log.Errorf(ctx, "Invalid Time range query paramater %v", timeRange)
		return c.Status(fiber.StatusBadRequest).JSON(fmt.Sprintf("Invalid Time Range. Err: %s", errors.New("Invalid Time Range value")))
	}
	_, err := strconv.ParseFloat(timeRange, 64)
	if err != nil {
		log.Errorf(ctx, "Invalid Time range query paramater %v", timeRange)
		return c.Status(fiber.StatusBadRequest).JSON(fmt.Sprintf("Invalid Time Range. Err: %s", errors.New("Invalid Time Range value")))
	}

	resp, err := getClientsUsage(ctx, timeRange, "dnn", "", "")
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fmt.Sprintf("Failed to read from the influxDB. Err: %s", err))

	}

	tags := []string{
		"dnn",
	}
	fields := []string{
		"bytesIn",
		"bytesOut",
	}

	data, _ := extractClientsDbResults(resp, tags, fields)

	return c.Status(fiber.StatusOK).JSON(data)

}

func getClientsUsage(ctx context.Context, timeRange string, groupByKey, whereKey, whereValue string) (*influxdb2.Result, error) {
	clientCon, err := influxDB.GetClient(ctx)
	if err != nil {
		log.Errorf(ctx, "Failed to connect to InfluxDB %v", err)
		return nil, err
	}

	cmd := fmt.Sprintf(`SELECT SUM("bytesToClient") AS bytesIn, SUM("bytesFromClient") AS bytesOut FROM "%s" WHERE time > now()-%sh`,
		measurementUsage,
		timeRange) +
		influxDB.SpaceString

	if whereKey != "" && whereValue != "" {
		cmd += fmt.Sprintf(`AND "%s" = '%s'`, whereKey, whereValue) + influxDB.SpaceString

	}

	groupByCmd := fmt.Sprintf(`GROUP BY "%s" , "%s"`, "clientId", "clientType") + influxDB.SpaceString
	if groupByKey != "" {
		groupByCmd = fmt.Sprintf(`GROUP BY "%s"`, groupByKey) + influxDB.SpaceString

	}
	cmd += groupByCmd

	log.Infof(ctx, "Influx cmd : ", cmd)

	resp, err := influxDB.Read(ctx, clientCon, clientsDatabase, cmd)

	if err != nil || len(resp) == 0 {
		log.Errorf(ctx, "Failed to read clients measurement %v", err)
		return nil, err
	}

	return &resp[0], nil

}

func extractClientsDbResults(res *influxdb2.Result, tags []string, fields []string) ([]map[string]interface{}, int) {
	list := make([]map[string]interface{}, 0)
	series := res.Series
	if len(series) == 0 {
		return list, len(list)
	}
	for _, val := range series {
		data := extractClientsRowData(val, tags, fields)

		list = append(list, data)

	}

	return list, len(list)

}

func extractClientsRowData(val models.Row, tags []string, fields []string) map[string]interface{} {
	data := map[string]interface{}{}
	influxDbResTags := val.Tags
	for _, tag := range tags {
		if tagVal, ok := influxDbResTags[tag]; ok {
			data[tag] = tagVal
		}

	}

	if len(val.Columns) != 0 && len(val.Values) != 0 {
		influxDbResCols := val.Columns

		influxDbResValues := val.Values[0]

		for _, field := range fields {
			idx := utils.IndexOfForStringArray(field, influxDbResCols)
			if idx != -1 {
				data[field] = influxDbResValues[idx]
			}

		}
	}
	return data
}
func extractClientsByClientIdFromDbResults(res *influxdb2.Result, tags []string, fields []string) map[string]interface{} {
	data := map[string]interface{}{}
	series := res.Series
	if len(series) == 0 {
		return data
	}
	return extractClientsRowData(series[0], tags, fields)
}

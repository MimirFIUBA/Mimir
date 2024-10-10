package inlfuxdb

import (
	"errors"
	"os"

	"github.com/InfluxCommunity/influxdb3-go/influxdb3"
)

// Connect to an Influx Database reading the credentials from
// environement variables INFLUXDB_TOKEN, INFLUXDB_URL
// return influxdb Client or errors
func ConnectToInfluxDB() (*influxdb3.Client, error) {

	dbToken := os.Getenv("INFLUXDB_TOKEN")
	if dbToken == "" {
		return nil, errors.New("INFLUXDB_TOKEN must be set")
	}

	dbURL := os.Getenv("INFLUXDB_URL")
	if dbURL == "" {
		return nil, errors.New("INFLUXDB_URL must be set")
	}

	bucketName := os.Getenv("INFLUXDB_BUCKET")
	if bucketName == "" {
		return nil, errors.New("INFLUXDB_BUCKET must be set")
	}

	dbOrg := os.Getenv("INFLUXDB_ORG")
	if dbOrg == "" {
		return nil, errors.New("INFLUXDB_ORG must be set")
	}

	client, err := influxdb3.New(influxdb3.ClientConfig{
		Host:         dbURL,
		Token:        dbToken,
		Database:     bucketName,
		Organization: dbOrg,
	})

	// client := influxdb3.NewClientWithOptions(dbURL, dbToken, influxdb3.DefaultOptions().SetBatchSize(20))

	// validate client connection health
	// _, err := client.Health(context.Background())

	return client, err
}

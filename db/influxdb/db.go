package inlfuxdb

import (
	"errors"
	"os"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

// Connect to an Influx Database reading the credentials from
// environement variables INFLUXDB_TOKEN, INFLUXDB_URL
// return influxdb Client or errors
func ConnectToInfluxDB() (influxdb2.Client, error) {
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

	client := influxdb2.NewClient(dbURL, dbToken)

	return client, nil
}

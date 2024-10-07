package db_test

import (
	"mimir/db"
	"testing"

	"github.com/joho/godotenv"
)

func Test_connectToInfluxDB(t *testing.T) {

	//load environment variable from a file for test purposes
	godotenv.Load("../db/test_influxdb.env")

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "Successful connection to InfluxDB",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := db.ConnectToInfluxDB()
			if (err != nil) != tt.wantErr {
				t.Errorf("ConnectToInfluxDB() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// health, err := got.Health(context.Background())
			// if (err != nil) && health.Status == domain.HealthCheckStatusPass {
			// 	t.Errorf("connectToInfluxDB() error. database not healthy")
			// 	return
			// }
			got.Close()
		})
	}
}

# Mimir
## Monitoring system for hydroponics 

### Server required configuration

For the mongo database, create the file db/mongodb/test_mongodb.env with the content:
```
MONGODB_PROTOCOL=mongodb
MONGODB_USERNAME=tp
MONGODB_PASSWORD=tp
MONGODB_HOSTNAME=localhost:27017
```

For the influx database, create the file db/influxdb/test_influxdb.env with the content:
```
INFLUXDB_TOKEN=tp
INFLUXDB_URL=localhost:8086
INFLUXDB_BUCKET=Mimir
INFLUXDB_ORG=tp
```

### Server startup

To start all the docker services required by the project at once:
```
sudo INFLUXDB_USERNAME=tp INFLUXDB_PASSWORD=tp INFLUXDB_ORG=tp INFLUXDB_BUCKET=Mimir INFLUXDB_TOKEN=tp MONGODB_USERNAME=tp MONGODB_PASSWORD=tp docker compose up -d
```

To start the server run:
```
go run cmd/mimir-server/main.go
```
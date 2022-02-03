package influxdb

import (
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	log "github.com/sirupsen/logrus"
)

func New(hostname string, token string) influxdb2.Client {
	client := influxdb2.NewClient(hostname, token)
	log.Infof("Successfully connected to InfluxDB at %s", hostname)
	return client
}

package influxdb

import (
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func New(hostname string, token string) influxdb2.Client {
	client := influxdb2.NewClient(hostname, token)
	return client
}

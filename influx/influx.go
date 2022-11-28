package influx

import (
	"context"
	"fmt"

	"github.com/influxdata/influxdb-client-go/v2"
)

func NewInfluxClient(server, token string) (*influxdb2.Client, error) {
	newClient := influxdb2.NewClient(server, token)
	_, err := newClient.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("influx错误:%s", err.Error())
	}
	return &newClient, nil
}

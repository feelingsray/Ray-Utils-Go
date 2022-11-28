package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

func NewRedisClient(host string, port int, passwd string, db int) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: passwd,
		DB:       db,
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return rdb, nil
}

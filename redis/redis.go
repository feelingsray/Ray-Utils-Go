package redis

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

type RedisStore struct {
	Host     string
	Port     int
	Password string
	Db       int
}

func NewRedisStore(host string, password string, db int) *RedisStore {
	return &RedisStore{
		Host:     host,
		Password: password,
		Port:     6379,
		Db:       db,
	}
}

func (r *RedisStore) GetRStore() (redis.Conn, error) {
	c, err := redis.Dial(
		"tcp", fmt.Sprintf("%s:%d", r.Host, r.Port),
		redis.DialPassword(r.Password),
		redis.DialDatabase(r.Db),
		redis.DialConnectTimeout(5*time.Second),
		redis.DialReadTimeout(1*time.Second),
		redis.DialWriteTimeout(1*time.Second),
		redis.DialKeepAlive(1*time.Second),
	)
	if err != nil {
		return nil, err
	} else {
		return c, nil
	}
}

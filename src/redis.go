package main

import (
	"github.com/garyburd/redigo/redis"
	"strconv"
	"time"
)

var redis_pool *redis.Pool = nil

func init() {
	redis_pool = newPool()
}

func newPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:   80,
		MaxActive: 1024,
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", config.Redis.Hostname+":"+strconv.Itoa(config.Redis.Port.(int)))
			if err != nil {
				return nil, err
			}

			if config.Redis.Password != "" {
				if _, err := conn.Do("AUTH", config.Redis.Password); err != nil {
					conn.Close()
					return nil, err
				}
			}

			return conn, err
		},
	}

}

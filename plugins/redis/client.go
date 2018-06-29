package redis

import (
	"github.com/go-redis/redis"
	"time"
)

var Client *redis.Client

var Nil = redis.Nil

func Set(key string, value string, ttl time.Duration) error {
	return Client.Set(key, value, ttl).Err()
}

func Get(key string) (string, error) {
	return Client.Get(key).Result()
}

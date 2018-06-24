package redis

import (
	"github.com/mholt/caddy"
	"github.com/go-redis/redis"
	"strconv"
	"errors"
)

func init() {
	caddy.RegisterPlugin("redis", caddy.Plugin{
		ServerType: "http",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	redisCfg, err := redisConfigParse(c)
	if err != nil {
		return err
	}
	Client = redis.NewClient(redisCfg)
	return nil
}

func redisConfigParse(c *caddy.Controller) (*redis.Options, error) {
	redisOption := &redis.Options{}
	for c.Next() {
		if len(redisOption.Addr) > 0 {
			return redisOption, errors.New("duplication redis config")
		}
		for c.NextBlock() {
			switch c.Val() {
			case "addr":
				if !c.NextArg() {
					return redisOption, c.ArgErr()
				}
				redisOption.Addr = c.Val()
			case "password":
				if !c.NextArg() {
					return redisOption, c.ArgErr()
				}
				redisOption.Password = c.Val()
			case "db":
				if !c.NextArg() {
					return redisOption, c.ArgErr()
				}
				db, err := strconv.Atoi(c.Val())
				if err != nil {
					return redisOption, err
				}
				redisOption.DB = db
			}
		}
	}
	return redisOption, nil
}

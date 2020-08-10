package redisClient

import (
	"time"

	"github.com/astaxie/beego"

	"github.com/garyburd/redigo/redis"
)

func PoolConnect() redis.Conn {
	pool := &redis.Pool{
		MaxIdle:     1,  // 最大的空闲连接数
		MaxActive:   10, // 最大连接数
		IdleTimeout: 180 * time.Second,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", beego.AppConfig.String("redisDB"))
			if err != nil {
				return nil, err
			}
			return c, nil
		},
	}
	return pool.Get()
}

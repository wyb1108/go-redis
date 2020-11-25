package redis

import (
	"errors"
	"fmt"
	"github.com/FZambia/sentinel"
	"github.com/gomodule/redigo/redis"
	"github.com/wyb1108/go-utils/util"
	"strings"
	"time"
)

var (
	pool *redis.Pool
)

func Init(file string) {
	err := initConfig(file)
	util.CheckErrorAndExit(err)
	if redisConfig.RedisServers == "" {
		fmt.Println("redis服务地址未配置, 请在redisConfig.json中配置")
	}
	if redisConfig.RedisSentinel {
		sntnl := &sentinel.Sentinel{
			Addrs:      strings.Split(redisConfig.RedisServers, ","),
			MasterName: redisConfig.RedisMasterName,
			Dial: func(addr string) (redis.Conn, error) {
				timeout := 500 * time.Millisecond
				c, err := redis.Dial("tcp", addr, redis.DialConnectTimeout(timeout), redis.DialReadTimeout(timeout), redis.DialWriteTimeout(timeout))
				if err != nil {
					return nil, err
				}
				return c, nil
			},
		}
		pool = &redis.Pool{
			MaxActive:   redisConfig.PoolMaxActive,
			MaxIdle:     redisConfig.PoolMinIdle,
			Wait:        true,
			IdleTimeout: 240 * time.Second,
			Dial: func() (redis.Conn, error) {
				masterAddr, err := sntnl.MasterAddr()
				if err != nil {
					return nil, err
				}
				c, err := redis.Dial("tcp", masterAddr)
				if err != nil {
					return nil, err
				}
				return c, nil
			},
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				if !sentinel.TestRole(c, "master") {
					return errors.New("role check failed")
				} else {
					return nil
				}
			},
		}
	} else {
		pool = &redis.Pool{
			MaxActive: redisConfig.PoolMaxActive,
			MaxIdle:   redisConfig.PoolMinIdle,

			IdleTimeout: 300 * time.Second,

			Dial: func() (conn redis.Conn, e error) {
				c, err := redis.Dial("tcp", redisConfig.RedisServers,
					redis.DialPassword(redisConfig.RedisPassword),
					redis.DialConnectTimeout(time.Duration(redisConfig.RedisTimeout)*time.Millisecond),
					redis.DialReadTimeout(time.Duration(redisConfig.RedisTimeout)*time.Millisecond),
					redis.DialWriteTimeout(time.Duration(redisConfig.RedisTimeout)*time.Millisecond))
				util.CheckError(err)
				return c, err
			},
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				if time.Since(t) < time.Minute {
					return nil
				}
				_, err := c.Do("PING")
				return err
			},
			MaxConnLifetime: 300 * time.Millisecond,
			Wait:            true,
		}
	}
}

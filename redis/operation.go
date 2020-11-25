package redis

import (
	"github.com/gomodule/redigo/redis"
)

func Get(key string) string {
	conn := pool.Get()
	defer conn.Close()
	value, _ := redis.String(conn.Do("GET", key))
	return value
}

func Set(key, value string) error {
	conn := pool.Get()
	defer conn.Close()

	_, err := redis.String(conn.Do("SET", key, value))
	return err
}

func SetWithExpire(key, value string, expireTime int) error {
	conn := pool.Get()
	defer conn.Close()

	_, err := redis.String(conn.Do("SETEX", key, expireTime, value))
	return err
}

func SetNX(key, value string) (string, error) {
	conn := pool.Get()
	defer conn.Close()

	return redis.String(conn.Do("SETNX", key, value))
}

func Expire(key string, expireTime int) error {
	conn := pool.Get()
	defer conn.Close()

	_, err := redis.String(conn.Do("EXPIRE", key, expireTime))
	return err
}

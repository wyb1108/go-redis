package redis

import "github.com/wyb1108/go-utils/util"

type redisConf struct {
	RedisServers    string `json:"redis.servers"`
	RedisPassword   string `json:"redis.password"`
	RedisSentinel   bool   `json:"redis.sentinel"`
	RedisMasterName string `json:"redis.masterName"`
	RedisTimeout    int    `json:"redis.timeout"`

	PoolMaxActive     int  `json:"pool.maxActive"`
	PoolMinIdle       int  `json:"pool.minIdle"`
	PollMaxWaitMillis int  `json:"poll.maxWaitMillis"`
	PoolTestOnBorrow  bool `json:"pool.testOnBorrow"`
	PoolTestOnReturn  bool `json:"pool.testOnReturn"`
	PoolTestWhileIdle bool `json:"pool.testWhileIdle"`
}

var (
	redisConfig redisConf
)

func initConfig(configFile string) error {
	return util.ParseJsonFile(configFile, &redisConfig)
}

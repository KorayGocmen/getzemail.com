package main

import (
	redis "github.com/go-redis/redis/v7"
)

var (
	redisdb *redis.Client
)

func initRedis() {
	logger.Println("Creating redis client")

	redisdb = redis.NewClient(&redis.Options{
		Addr:         config.Redis.Addr,
		Password:     config.Redis.Pass,
		DB:           config.Redis.DB,
		PoolSize:     config.Redis.PoolSize,
		MinIdleConns: config.Redis.MinIdleConns,
		IdleTimeout:  timeDuration(config.Redis.IdleTimeout),
	})

	if _, err := redisdb.Ping().Result(); err != nil {
		logger.Fatalln("Failed to connect to redis", err)
	}
}

package config

import (
    "github.com/go-redis/redis/v8"
	"fmt"
)

func SetupRedis() *redis.Client {
    return redis.NewClient(&redis.Options{
        Addr: fmt.Sprintf("%s:%d", 
            "redis",
            6379,
        ),
    })
}
package database

import (
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
)

func GetClient() *redis.Client {
	addr := os.Getenv("REDIS_SERVER")
	port := os.Getenv("REDIS_PORT")
	pwd := os.Getenv("REDIS_PWD")

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", addr, port),
		Password: pwd,
	})
	return rdb
}

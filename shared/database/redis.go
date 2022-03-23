package database

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

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

func Get[T SavebleItem](keyPattern string, id string) (T, error) {
	ctx := context.Background()
	rdb := GetClient()

	key := fmt.Sprintf("%s:%s", keyPattern, id)
	result, err := rdb.Get(ctx, key).Result()
	if err != nil {
		fmt.Printf("[Redis][Get][Error] %v\n", err.Error())
	}
	var item T
	json.Unmarshal([]byte(result), &item)
	return item, err
}

func Save[T SavebleItem](keyPattern string, item T) error {
	ctx := context.Background()
	rdb := GetClient()

	key := fmt.Sprintf("%s:%s", keyPattern, item.GetId())
	value, _ := json.Marshal(item)
	rdb.LPush(ctx, keyPattern, item.GetId()).Err()
	return rdb.Set(ctx, key, value, 0).Err()
}

func List[T SavebleItem](keyPattern string, offset int64, limit int64) ([]T, error) {
	ctx := context.Background()
	rdb := GetClient()

	if limit == 0 {
		limit = 100
	}
	keys, err := rdb.Sort(
		ctx,
		"urls",
		&redis.Sort{
			By:     "id",
			Offset: offset,
			Count:  limit,
			Order:  "ASC",
		},
	).Result()

	if err != nil {
		fmt.Printf("[Redis][List][Error] %v\n", err.Error())
		return []T{}, err
	}

	var itens []T
	for _, key := range keys {
		eachItem, _ := Get[T]("urls", key)
		itens = append(itens, eachItem)
	}

	return itens, nil
}

func GetCounter(key string) int64 {
	ctx := context.Background()
	rdb := GetClient()
	v, _ := rdb.Get(ctx, key).Result()
	counter, _ := strconv.Atoi(v)
	return int64(counter)
}

func Increment(key string) {
	ctx := context.Background()
	rdb := GetClient()
	_, err := rdb.Incr(ctx, key).Result()
	if err != nil {
		fmt.Printf("[Redis][List][Increment] %v\n", err.Error())
	}
}

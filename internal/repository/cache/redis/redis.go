package rediscache

import (
	"context"
	"encoding/json"
	"hezzltask5/internal/models"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

const DataKey = "datakey"

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(client *redis.Client) *RedisCache {

	return &RedisCache{
		client: client,
	}
}

func NewRedisClient(addr string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (rc *RedisCache) Set(ctx context.Context, items []models.Item) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()
	data, err := json.Marshal(items)
	if err != nil {
		logrus.Printf("cachelevel: set method:error: %v", err)

	}

	if err := rc.client.Set(ctx, DataKey, data, 60*time.Second).Err(); err != nil {
		logrus.Printf("cachelevel: set method:error: %v", err)

	}
}

func (rc *RedisCache) Get(ctx context.Context) ([]models.Item, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()
	data, err := rc.client.Get(ctx, DataKey).Bytes()
	if err != nil {
		return nil, err
	}
	var items []models.Item
	err = json.Unmarshal(data, &items)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (rc *RedisCache) Remove(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()
	_, err := rc.client.Del(ctx, DataKey).Result()
	if err != nil {
		logrus.Printf("cachelevel: del method:error: %v", err)

	}
}

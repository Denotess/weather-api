package models

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache interface {
	Get(ctx context.Context, key string) (CachedWeather, bool, error)
	Set(ctx context.Context, key string, value CachedWeather, ttl time.Duration) error
}

type CachedWeather struct {
	Temp       float64 `json:"temp"`
	Conditions string  `json:"conditions"`
}

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{client: client}
}

func (c *RedisCache) Get(ctx context.Context, key string) (CachedWeather, bool, error) {
	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return CachedWeather{}, false, nil
	}
	if err != nil {
		return CachedWeather{}, false, err
	}

	var data CachedWeather
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		return CachedWeather{}, false, err
	}
	return data, true, nil
}

func (c *RedisCache) Set(ctx context.Context, key string, value CachedWeather, ttl time.Duration) error {
	payload, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, payload, ttl).Err()
}

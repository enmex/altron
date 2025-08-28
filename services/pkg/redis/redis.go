package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

var defaultExpirationTime = time.Duration(24 * time.Hour)

type Client[T any] struct {
	cfg            *Config
	client         *redis.Client
	expirationTime time.Duration
}

func NewRedisCache[T any](cfg *Config) *Client[T] {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       1,
	})

	return &Client[T]{
		cfg:            cfg,
		client:         client,
		expirationTime: defaultExpirationTime,
	}
}

func (w *Client[T]) WithExpirationTime(expirationTime time.Duration) *Client[T] {
	w.expirationTime = expirationTime
	return w
}

func (w *Client[T]) Set(ctx context.Context, key string, value T) error {
	json, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return w.client.Set(ctx, w.cfg.Prefix+key, string(json), w.expirationTime).Err()
}

func (w *Client[T]) Get(ctx context.Context, key string) (*T, error) {
	var jsonStr string
	err := w.client.Get(ctx, w.cfg.Prefix+key).Scan(&jsonStr)
	if len(jsonStr) == 0 {
		return nil, nil
	}

	var value T
	if err := json.Unmarshal([]byte(jsonStr), &value); err != nil {
		return nil, err
	}

	return &value, err
}

func (w *Client[T]) Delete(ctx context.Context, key string) error {
	return w.client.Del(ctx, w.cfg.Prefix+key).Err()
}

func (w *Client[T]) Exist(ctx context.Context, key string) (bool, error) {
	exists, err := w.client.Exists(ctx, w.cfg.Prefix+key).Result()
	if err != nil {
		return false, err
	}
	return exists == 1, nil
}

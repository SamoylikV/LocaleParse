package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/SamoylikV/LocaleParse/internal/config"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

type Client struct {
	*redis.Client
	ctx context.Context
}

func NewClient(cfg *config.Config) (*Client, error) {
	redisDB, _ := strconv.Atoi(cfg.RedisDB)
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       redisDB,
	})

	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}
	return &Client{client, ctx}, nil
}

func (c *Client) GetLocaleData(key string) (map[string]map[string]string, error) {
	data, err := c.Client.Get(c.ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("failed to get data from redis: %w", err)
	}

	var localeData map[string]map[string]string
	err = json.Unmarshal([]byte(data), &localeData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal redis data: %w", err)
	}
	return localeData, nil
}

func (c *Client) SetLocaleData(key string, data map[string]map[string]string, expiration time.Duration) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data to JSON: %w", err)
	}
	err = c.Client.Set(c.ctx, key, jsonData, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set data in redis: %w", err)
	}

	return nil
}

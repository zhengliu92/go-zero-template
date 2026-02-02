package svc

import (
	"context"
	"fmt"
	"go-zero-template/internal/config"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func MustInitRedis(redisConfig config.RedisConfig) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     redisConfig.Addr,
		Password: redisConfig.Password,
		DB:       redisConfig.DB,
	})

	// 检查 Ping，看看 Redis 能否打通
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatalf("failed to connect Redis: %v", err)
	}

	log.Println("Redis connected successfully")
	return client
}

// PingRedis 检查 Redis 连接是否正常
func PingRedis(client *redis.Client) error {
	if client == nil {
		return fmt.Errorf("redis client is not initialized")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return client.Ping(ctx).Err()
}

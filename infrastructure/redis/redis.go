package redis

import (
	"context"
	"sync"
	"time"

	"orchid-starter/config"

	"github.com/go-redis/redismock/v9"
	redisV9 "github.com/redis/go-redis/v9"
	"github.com/yudhiana/logos"
)

var redisClient *redisV9.Client
var redisOnce sync.Once

// GetRedisClient returns a singleton redis client configured according to the
// application configuration.  The call is safe for concurrent use.
func GetRedisClient(cfg *config.LocalConfig) *redisV9.Client {
	redisOnce.Do(func() {

		addr := cfg.RedisConfig.Addr()
		logos.NewLogger().Info("Initialize Redis connection", "addr", addr)

		// pooled connection configuration (tune PoolSize/MinIdleConns as needed)
		redisClient = redisV9.NewClient(&redisV9.Options{
			Addr:     addr,
			Username: cfg.RedisConfig.RedisUsername,
			Password: cfg.RedisConfig.RedisPassword, // "" if none
			DB:       cfg.RedisConfig.RedisDB,       // usually 0

			PoolSize:        cfg.RedisConfig.RedisPoolSize,
			MinIdleConns:    cfg.RedisConfig.RedisMinIdleConn,
			ConnMaxIdleTime: time.Duration(cfg.RedisConfig.RedisConnMaxIdleTime) * time.Second,
		})

		if cfg.RedisDebug {
			// Add hook before testing connection
			redisClient.AddHook(NewRedisHook())
		}

		// Test connection
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		// simple ping to verify the connection
		if err := redisClient.Ping(ctx).Err(); err != nil {
			panic(err)
		}
	})

	return redisClient
}

func GetMockRedisConnection() (*redisV9.Client, redismock.ClientMock) {
	return redismock.NewClientMock()
}

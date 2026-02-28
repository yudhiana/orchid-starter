package redis

import (
	"context"
	"log"
	"net"
	"time"

	redisV9 "github.com/redis/go-redis/v9"
)

type RedisHook struct {
}

func NewRedisHook() *RedisHook {
	return &RedisHook{}
}

func (h *RedisHook) ProcessHook(next redisV9.ProcessHook) redisV9.ProcessHook {
	return func(ctx context.Context, cmd redisV9.Cmder) error {
		// Ignore internal go-redis maintenance command
		if cmd.Name() == "client" {
			return next(ctx, cmd)
		}

		log.Println("REDIS COMMAND > ", cmd.String())
		err := next(ctx, cmd)
		if err != nil {
			log.Println("REDIS ERROR > ", err)
		}
		return err
	}
}

func (h *RedisHook) ProcessPipelineHook(next redisV9.ProcessPipelineHook) redisV9.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redisV9.Cmder) error {
		log.Println("REDIS PIPELINE > ", len(cmds), "commands")
		return next(ctx, cmds)
	}
}

func (h *RedisHook) DialHook(next redisV9.DialHook) redisV9.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		log.Println("REDIS DIAL > ", "network:", network, "addr:", addr)

		start := time.Now()
		conn, err := next(ctx, network, addr)
		duration := time.Since(start)

		if err != nil {
			log.Printf("REDIS DIAL ERROR > %v (duration: %v)\n", err, duration)
			return nil, err
		}

		log.Printf("REDIS DIAL SUCCESS > %s (duration: %v)\n", addr, duration)
		return conn, nil
	}
}

package redis

import (
	"context"
	"fibera/internal/config"
	"sync"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var (
	redisCli *redis.Client
	once     sync.Once
	redisErr error
)

// GetRedis 获取Redis客户端
func GetRedis() (*redis.Client, error) {
	once.Do(func() {
		conf, _ := config.LoadConfig()

		redisCli = redis.NewClient(&redis.Options{
			Addr:     conf.Redis.Addr,
			Password: conf.Redis.Password,
			DB:       conf.Redis.DB,
		})

		// 测试连接
		_, redisErr = redisCli.Ping(context.Background()).Result()
		if redisErr != nil {
			zap.L().Error("failed to connect to redis", zap.Error(redisErr))
			return
		}
	})
	return redisCli, nil
}

package redis

import (
	"ForumWeb/setting"

	"github.com/go-redis/redis"
)

func Init(config *setting.RedisConfig) (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       config.Db,
		PoolSize: config.PoolSize,
	})
	_, err = rdb.Ping().Result()
	return
}

func Close() error {
	return rdb.Close()
}

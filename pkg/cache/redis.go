package cache

import (
	"gin-api/pkg/config"
	"gin-api/pkg/redis"
	"time"
)

// RedisDriver 实现 cache.CacheInterface
type RedisDriver struct {
	RedisClient *redis.RedisClient
	KeyPrefix   string
}

//NewRedis 返回实现了 cache.CacheInterface 的 RedisDriver
func NewRedis(address string, username string, password string, db int) *RedisDriver {
	rs             := &RedisDriver{}
	rs.RedisClient = redis.NewClient(address, username, password, db)
	rs.KeyPrefix   = config.GetString("app.name") + ":cache:"
	return rs
}

func (s *RedisDriver) Set(key string, value string, expireTime time.Duration) {
	s.RedisClient.Set(s.KeyPrefix+key, value, expireTime)
}

func (s *RedisDriver) Get(key string) string {
	return s.RedisClient.Get(s.KeyPrefix + key)
}

func (s *RedisDriver) Has(key string) bool {
	return s.RedisClient.Has(s.KeyPrefix + key)
}

func (s *RedisDriver) Forget(key string) {
	s.RedisClient.Del((s.KeyPrefix + key))
}

func (s *RedisDriver) Forever(key string, value string) {
	s.RedisClient.Set(s.KeyPrefix+key, value, 0)
}

func (s *RedisDriver) Flush() {
	s.RedisClient.FlushDB()
}

func (s *RedisDriver) Increment(parameters ...interface{}) {
	s.RedisClient.Increment(parameters...)
}

func (s *RedisDriver) Decrement(parameters ...interface{}) {
	s.RedisClient.Decrement(parameters...)
}

func (s *RedisDriver) IsAlive() error {
	return s.RedisClient.Ping()
}

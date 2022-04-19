// Package redis 工具包
package redis

import (
	"context"
	"fmt"
	redis "github.com/go-redis/redis/v8"
	"gin-api/pkg/config"
	"gin-api/pkg/logger"

	"sync"
	"time"
)

// RedisClient redisClient 服务
type RedisClient struct {
	Client *redis.Client
	Ctx    context.Context
}

// once 确保默认的 redisClient 对象只实例一次
var once sync.Once

var defaultClient *RedisClient

//DefaultClient 基于默认配置初始化 RedisClient
func DefaultClient() *RedisClient {
	once.Do(func() {
		address := fmt.Sprintf("%v:%v", config.GetString("redis.host"), config.GetString("redis.port"))
		username := config.GetString("redis.username")
		password := config.GetString("redis.password")
		dbIndex := config.GetInt("redis.database")
		defaultClient = NewClient(address, username, password, dbIndex)
	})
	return defaultClient
}

//NewClient 基于自定义配置初始化 RedisClient
func NewClient(address string, password string, username string, dbIndex int) *RedisClient {
	// 使用 redis 库里的 connect 初始化连接
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Username: username,
		Password: password,
		DB:       dbIndex,
	})

	//context
	ctx := context.Background()

	//初始化 RedisClient 对象
	rds := &RedisClient{Client: client, Ctx: ctx}

	// 测试连接
	err := rds.Ping()
	logger.LogIf("redis-Ping", err)

	return rds
}

//Lock 获取锁
func (rds *RedisClient) Lock(key string, expire time.Duration) bool {
	lockValue := 1
	err := rds.Client.SetNX(rds.Ctx, key, lockValue, expire).Err()
	if err != nil {
		logger.Log("Lock", err.Error())
		//获取锁失败后，检查是不是死锁
		if rds.Client.TTL(rds.Ctx, key).Val() == time.Duration(-1) {
			rds.Client.Expire(rds.Ctx, key, expire)
		}
		return false
	}
	return true
}

//ReleaseLock 释放锁
func (rds *RedisClient) ReleaseLock(key string) bool {
	err := rds.Client.Del(rds.Ctx, key).Err()
	if err != nil {
		logger.Log("ReleaseLock", err.Error())
		return false
	}
	return true
}

//Ping 用来测试连接
func (rds *RedisClient) Ping() error {
	_, err := rds.Client.Ping(rds.Ctx).Result()
	if err != nil {
		logger.Log("Ping", err.Error())
	}
	return err
}

//Redis 返回原生客户端 redis.Client
func (rds *RedisClient) Redis() *redis.Client {
	return rds.Client
}

// Set 存储 key 对应的 value，且设置 expiration 过期时间
func (rds *RedisClient) Set(key string, value interface{}, expiration time.Duration) bool {
	if err := rds.Client.Set(rds.Ctx, key, value, expiration).Err(); err != nil {
		logger.Log("Set", err.Error())
		return false
	}
	return true
}

// Get 获取 key 对应的 value
func (rds *RedisClient) Get(key string) string {
	result, err := rds.Client.Get(rds.Ctx, key).Result()
	if err != nil {
		if err != redis.Nil {
			logger.Log("Get", err.Error())
		}
		return ""
	}
	return result
}

// HGet 获取 hash 中指定字段的值
func (rds *RedisClient) HGet(key, field string) string {
	result, err := rds.Client.HGet(rds.Ctx, key, field).Result()
	if err != nil {
		if err != redis.Nil {
			logger.Log("Get", err.Error())
		}
		return ""
	}
	return result
}

//HGetAll 获取 hash 中指定 key 的所有字段和值
func (rds *RedisClient) HGetAll(key string) map[string]string {
	result, err := rds.Client.HGetAll(rds.Ctx, key).Result()
	if err != nil {
		if err != redis.Nil {
			logger.Log("HGetAll", err.Error())
		}
		return map[string]string{}
	}
	return result
}

func (rds *RedisClient) HSet(key string, values ...interface{}) bool {
	err := rds.Client.HSet(rds.Ctx, key, values...).Err()
	if err != nil {
		logger.Log("HSet", err.Error())
		return false

	}
	return true
}

func (rds *RedisClient) Expire(key string, expire time.Duration) bool {
	if err := rds.Client.Expire(rds.Ctx, key, expire).Err(); err != nil {
		logger.Log("Expire", err.Error())
		return false
	}
	return true
}

// Has 判断一个 key 是否存在，内部错误和 redis.Nil 都返回 false
func (rds *RedisClient) Has(key string) bool {
	_, err := rds.Client.Get(rds.Ctx, key).Result()
	if err != nil {
		if err != redis.Nil {
			logger.Log("Has", err.Error())
		}
		return false
	}
	return true
}

// Del 删除存储在 redis 里的数据，支持多个 key 传参
func (rds *RedisClient) Del(keys ...string) bool {
	if err := rds.Client.Del(rds.Ctx, keys...).Err(); err != nil {
		logger.Log("Del", err.Error())
		return false
	}
	return true
}

// FlushDB 清空当前 redis db 里的所有数据
func (rds *RedisClient) FlushDB(keys ...string) bool {
	if err := rds.Client.FlushDB(rds.Ctx).Err(); err != nil {
		logger.Log("FlushDB", err.Error())
		return false
	}
	return true
}

// Increment 当参数只有 1 个时，为 key，其值增加 1。
// 当参数有 2 个时，第一个参数为 key ，第二个参数为要增加的值 int64 类型。
func (rds RedisClient) Increment(parameters ...interface{}) bool {
	switch len(parameters) {
	case 1:
		key := parameters[0].(string)
		if err := rds.Client.Incr(rds.Ctx, key).Err(); err != nil {
			logger.Log("Increment", err.Error())
			return false
		}
	case 2:
		key := parameters[0].(string)
		value := parameters[0].(int64)
		if err := rds.Client.IncrBy(rds.Ctx, key, value).Err(); err != nil {
			logger.Log("Increment", err.Error())
			return false
		}
	default:
		logger.Log("Increment", "参数过多")
		return false
	}
	return true
}

// Decrement 当参数只有 1 个时，为 key，其值减去 1。
// 当参数有 2 个时，第一个参数为 key ，第二个参数为要减去的值 int64 类型。
func (rds RedisClient) Decrement(parameters ...interface{}) bool {
	switch len(parameters) {
	case 1:
		key := parameters[0].(string)
		if err := rds.Client.Decr(rds.Ctx, key).Err(); err != nil {
			logger.Log("Decrement", err.Error())
			return false
		}
	case 2:
		key := parameters[0].(string)
		value := parameters[0].(int64)
		if err := rds.Client.DecrBy(rds.Ctx, key, value).Err(); err != nil {
			logger.Log("Decrement", err.Error())
			return false
		}
	default:
		logger.Log("Decrement", "参数过多")
		return false
	}
	return true
}

// Package cache 缓存工具类，可以缓存各种类型包括 struct 对象
package cache

import (
	"encoding/json"
	"gin-api/pkg/helpers"
	"gin-api/pkg/logger"
	"sync"
	"time"
)

type Cache struct {
	Driver CacheInterface
}

var once sync.Once

var cache *Cache

//Init 接收实现了 CacheInterface 的驱动来实现 Cache 对象
func Init(driver CacheInterface) {
	once.Do(func() {
		cache = &Cache{
			Driver: driver,
		}
	})
}

func Set(key string, obj interface{}, expireTime time.Duration) {
	b, err := json.Marshal(&obj)
	logger.LogIf(helpers.CurrentFuncName(), err)
	cache.Driver.Set(key, string(b), expireTime)
}

func Get(key string) interface{} {
	stringValue := cache.Driver.Get(key)
	var wanted interface{}
	err := json.Unmarshal([]byte(stringValue), &wanted)
	logger.LogIf(helpers.CurrentFuncName(), err)
	return wanted
}

func Has(key string) bool {
	return cache.Driver.Has(key)
}

// GetObject 应该传地址，用法如下:
// 	model := user.User{}
// 	cache.GetObject("key", &model)
func GetObject(key string, wanted interface{}) {
	val := cache.Driver.Get(key)
	if len(val) > 0 {
		err := json.Unmarshal([]byte(val), &wanted)
		logger.LogIf(helpers.CurrentFuncName(), err)
	}
}

func Forget(key string) {
	cache.Driver.Forget(key)
}

func Forever(key string, value string) {
	cache.Driver.Set(key, value, 0)
}

func Flush() {
	cache.Driver.Flush()
}

func Increment(parameters ...interface{}) {
	cache.Driver.Increment(parameters...)
}

func Decrement(parameters ...interface{}) {
	cache.Driver.Decrement(parameters...)
}

func IsAlive() error {
	return cache.Driver.IsAlive()
}

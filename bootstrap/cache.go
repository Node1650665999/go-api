package bootstrap

import "fmt"
import "gin-api/pkg/config"
import "gin-api/pkg/cache"

func setupCache()  {
	address  := fmt.Sprintf("%v:%v", config.GetString("redis.host"), config.GetString("redis.port"))
	username := config.GetString("redis.username")
	password := config.GetString("redis.password")
	dbIndex  := config.GetInt("redis.database")
	driver   := cache.NewRedis(address, username, password, dbIndex)
	cache.Init(driver)
}
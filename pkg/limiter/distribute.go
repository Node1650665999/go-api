package limiter

import (
	"errors"
	"github.com/spf13/cast"
	"math"
	"gin-api/pkg/config"
	"gin-api/pkg/redis"
	"sync"
	"time"
)

//DistributeBucket 分布式限流
type DistributeBucket struct {
	redis      *redis.RedisClient
	keyPrefix  string
	lockPrefix string
	rate       int64 //固定的token放入速率, r/s
	capacity   int64 //桶的容量
	lock       sync.Mutex
}

//NewDistributeBucket 实例化一个分布式限流器
func NewDistributeBucket() *DistributeBucket {
	return &DistributeBucket{
		redis:      redis.DefaultClient(),
		keyPrefix:  config.GetString("app.name") + ":limiter:",
		lockPrefix: config.GetString("app.name") + ":lock:",
	}
}

//Check 检测请求是否超额
// key : 用来标识限流的对象, 可为单个 Ip, 单个Route, 或者 Ip+Route 的组合
// format : 即访问频率的控制, 格式如下：
//	5 reqs/second: "5-S"
//	10 reqs/minute: "10-M"
//	1000 reqs/hour: "1000-H"
//	2000 reqs/day: "2000-D"
func (d *DistributeBucket) Check(key string, format string) error {
	limit, everyDuration, err := ParseFormat(format)
	if err != nil {
		return err
	}

	//换算成每秒能处理几个,即令牌桶中的速率 rate(r/s)
	d.rate     = cast.ToInt64(math.Ceil(cast.ToFloat64(limit) / everyDuration.Seconds()))
	//容量设置为rate的5倍,即可以应对的最高突发性流量
	d.capacity = 10 * d.rate
	//多久能填满桶
	fillTime := math.Floor(float64(d.capacity / d.rate))
	//key 的过期时间设置为填满时间的2倍
	ttl      := time.Duration(2 * fillTime) * time.Second

	//获取分布式锁
	lockKey := d.lockPrefix + key
	d.redis.Lock(lockKey, ttl)
	defer d.redis.ReleaseLock(lockKey)

	//获取令牌
	fullKey := d.keyPrefix + key
	data := d.getTokenBucket(fullKey)

	var tokens int64
	var lastRefresh int64
	now := time.Now().Unix()

	if len(data) == 0 { // 没有值说明是第一次进入
		tokens = d.capacity
		lastRefresh = now
	} else { //根据速率和时间差,判断应放入多少令牌
		tokens = cast.ToInt64(tokens)
		lastRefresh = cast.ToInt64(data["last_refresh"])
		//重置令牌数量
		tokens = tokens + (now-lastRefresh)*d.rate
		if tokens > d.capacity {
			tokens = d.capacity
		}
		//重置令牌更新时间
		lastRefresh = now
	}

	//本次请求token数量是否足够
	if tokens > 0 {
		//允许本次请求,并计算token余量
		tokens = tokens - 1
		//更新令牌桶
		d.setTokenBucket(fullKey, tokens, lastRefresh, ttl)
		return nil
	} else {
		//没有令牌,则拒绝
		return errors.New("访问太频繁")
	}
}

//getTokenBucket 获取令牌桶
func (d *DistributeBucket) getTokenBucket(fullKey string)map[string]string{
	return d.redis.HGetAll(fullKey)
}

//setTokenBucket 更新令牌桶
func (d *DistributeBucket) setTokenBucket(key string, tokens, lastRefresh int64, expire time.Duration) {
	d.redis.HSet(key, "tokens", tokens, "last_refresh", lastRefresh)
	d.redis.Expire(key, expire)
}

//HealthCheck 检查分布式限流器,用来在单机和分布式之间切换
func (d *DistributeBucket) HealthCheck() bool{
	err := d.redis.Ping()
	if err != nil {
		return false
	}
	return true
}

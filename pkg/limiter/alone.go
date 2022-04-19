package limiter

import (
	"errors"
	"golang.org/x/time/rate"
	"sync"
	"time"
)

// 创建一个自定义visitor结构体，包含每个访问者的限流器和最后一次访问时间。
type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

//AloneBucket 单机版限流器
type AloneBucket struct {
	visitors map[string]*visitor
	mux      sync.Mutex
}

var (
	alone    *AloneBucket
	once     sync.Once
)

//NewAloneBucket 实例化单机版限流器
func NewAloneBucket() *AloneBucket {
	once.Do(func() {
		alone = &AloneBucket{
			visitors: make(map[string]*visitor),
		}
		//回收内存
		go gc()
	})
	return alone
}

//gc 将长时间没有访问的对象删除
func gc() {
	ticker := time.NewTicker(time.Minute)
	for range ticker.C {
		alone.mux.Lock()
		for key, v := range alone.visitors {
			if time.Since(v.lastSeen) > 3 * time.Minute {
				delete(alone.visitors, key)
			}
		}
		alone.mux.Unlock()
	}
}

//Check 检测请求是否超额
// key : 用来标识限流的对象, 可为单个 Ip, 单个Route, 或者 Ip+Route 的组合
// format : 即访问频率的控制, 格式如下：
//	5 reqs/second: "5-S"
//	10 reqs/minute: "10-M"
//	1000 reqs/hour: "1000-H"
//	2000 reqs/day: "2000-D"

func (a *AloneBucket) Check(key string, format string) error {
	a.mux.Lock()
	defer a.mux.Unlock()

	limit,everyDuration,err := ParseFormat(format)
	if err != nil {
		return err
	}

	var limiter *rate.Limiter

	v, exists := a.visitors[key]
	if !exists {
		limiter = rate.NewLimiter(rate.Every(everyDuration), limit)
		// 在创建新访问者时, 记录当前时间
		a.visitors[key] = &visitor{limiter, time.Now()}
	} else {
		v.lastSeen = time.Now()
		limiter = v.limiter
	}

	if limiter.Allow() == false {
		return errors.New("访问太频繁")
	}
	return nil
}

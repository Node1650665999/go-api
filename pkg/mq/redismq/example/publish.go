package main

import (
	"context"
	"fmt"
	"gin-api/pkg/mq/redismq"
	"math/rand"
	"time"
)

func main() {
	ctx := context.Background()
	current := time.Now().Format("2006-01-02 15:04:05")
	msg := fmt.Sprintf("rand_num:%v, 消息发送时间:%v", RandomIntRange(1000, 10000), current)

	//通过消息队列发送消息
	queueConfig := redismq.Config{
		QueueName: "msg.worker",
	}
	redismq.NewRedisMQ(nil, queueConfig).SendNormalMsg(ctx, []byte(msg))
	fmt.Println(msg)
}

func RandomIntRange(min, max int64) int64 {
	return min + rand.Int63n(max-min)
}

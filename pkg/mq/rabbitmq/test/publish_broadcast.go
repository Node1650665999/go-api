package main

import (
	"context"
	"fmt"
	"gin-api/pkg/mq/rabbitmq"
	"math/rand"
	"time"
)

func main() {
	ctx := context.Background()
	name := "msg.broadcast"
	current := time.Now().Format("2006-01-02 15:04:05")
	randNum := RandomIntRange(1000, 10000)
	msg := fmt.Sprintf("rand_num:%v, 消息发送时间:%v", randNum, current)

	//通过消息队列发送消息
	dsn := ""
	rabbitmq.NewRabbitMq(dsn).UseBroadCast(name).SendNormalMsg(ctx, []byte(msg))

	fmt.Println(msg)
}

func RandomIntRange(min, max int64) int64 {
	return min + rand.Int63n(max-min)
}

package main

import (
	"context"
	"fmt"
	"gin-api/pkg/mq"
	"gin-api/pkg/mq/rabbitmq"
	"math/rand"
	"time"
)

func main() {
	current := time.Now().Format("2006-01-02 15:04:05")
	msg := fmt.Sprintf("rand_num:%v, 消息发送时间:%v", RandomIntRange(1000, 10000), current)

	//通过消息队列发送消息
	dsn := ""
	queueConfig := rabbitmq.Config{
		QueueName: "msg.worker",
		Mode:      rabbitmq.Worker,
		LogPrintf: mq.Printf,
	}
	rabbitmq.NewRabbitMQ(dsn, queueConfig).SendNormalMsg(context.Background(), []byte(msg))
	fmt.Println(msg)
}

func RandomIntRange(min, max int64) int64 {
	return min + rand.Int63n(max-min)
}

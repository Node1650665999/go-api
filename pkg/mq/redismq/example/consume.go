package main

import (
	"context"
	"fmt"
	"gin-api/pkg/mq/redismq"
	"time"
)

func main() {

	ctx := context.Background()
	queueConfig := redismq.Config{
		QueueName: "msg.worker",
		LogPrintf: func(s string, i ...interface{}) {
			fmt.Printf(s, i...)
		},
	}
	redismq.NewRedisMQ(nil, queueConfig).ReceiveNormalMsg(ctx, HandlerMessage)

	fmt.Println("注意：这里永远无法输出，因为consume()会挂起进程以待消息")
}

func HandlerMessage(ctx context.Context, msgId string, msg []byte, extra interface{}) error {
	current := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("打印消息：msgid: %v, msg: %v，消息接收时间:%v \n", msgId, string(msg), current)
	return nil
}

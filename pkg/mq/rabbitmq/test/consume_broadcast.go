package main

import (
	"context"
	"fmt"
	"time"

	"gin-api/pkg/mq/rabbitmq"
)

func main() {

	ctx := context.Background()
	dsn := ""
	name := "msg.broadcast"
	rabbitmq.NewRabbitMq(dsn).UseBroadCast(name).ReceiveNormalMsg(ctx, HandlerBroadCastMessage)

	fmt.Println("注意：这里永远无法输出，因为consume()会挂起进程以待消息")
}

func HandlerBroadCastMessage(ctx context.Context, msgId string, msg []byte, extra interface{}) error {
	current := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("打印消息：msgid: %v, msg: %v，消息接收时间:%v \n", msgId, string(msg), current)
	return nil
}

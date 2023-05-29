package redismq

import (
	"context"
	"fmt"
	"time"

	"gin-api/pkg/mq"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/cast"
)

//确保 RedisMQ 实现了 mq.Contracts
var _ mq.Contracts = (*RedisMQ)(nil)

//Config 定义了redis-stream 配置信息,多协程下不同的 groupName和consumerName可以模拟出三种消费模式：
//   1.worker模式(一对一, 启动一个goroutine): groupName和consumerName传空就可以了
//   2.负载均衡模式(一对多, 启动多个goroutine均摊消息): queueName,groupName 必须相同,consumerName 必须不相同
//   3.广播模式(一对多, 启动多个goroutine广播式消费): queueName,consumerName 必须相同,groupName 必须不相同
type Config struct {
	QueueName    string
	GroupName    string
	ConsumerName string
	LogPrintf    func(string, ...interface{})
}

type RedisMQ struct {
	rdb  *redis.Client
	err  chan error
	done chan bool
	Config
}

func NewRedisMQ(rdb *redis.Client, queueConfig Config) *RedisMQ {
	if queueConfig.QueueName == "" {
		panic("please provide the queue name")
	}

	if queueConfig.LogPrintf == nil {
		queueConfig.LogPrintf = mq.Printf
	}

	redisMQ := &RedisMQ{
		rdb:  rdb,
		err:  make(chan error),
		done: make(chan bool),
	}

	return redisMQ.setConfig(queueConfig)
}

func (r *RedisMQ) setConfig(queueConfig Config) *RedisMQ {
	queueName := queueConfig.QueueName
	groupName := queueConfig.GroupName
	consumerName := queueConfig.ConsumerName
	if groupName == "" {
		groupName = queueConfig.QueueName
	}
	if consumerName == "" {
		consumerName = queueConfig.QueueName
	}

	queueConfig.QueueName = fmt.Sprintf("queue:stream.%v", queueName)
	queueConfig.GroupName = fmt.Sprintf("queue:stream.%v:group.%v", queueName, groupName)
	queueConfig.ConsumerName = fmt.Sprintf("queue:stream.%v:group.%v:consumer.%v", queueName, groupName, consumerName)
	r.Config = queueConfig

	go r.log()

	return r
}

func (r *RedisMQ) log() {
	for {
		select {
		case err, ok := <-r.err:
			if !ok {
				return
			}
			r.LogPrintf(
				"%v, consumer_name:%v, group_name:%v, queue_name:%v",
				err,
				r.Config.ConsumerName,
				r.Config.GroupName,
				r.Config.QueueName,
			)
		}
	}
}

func (r *RedisMQ) SendNormalMsg(ctx context.Context, msg []byte) (error, string) {
	return r.publish(ctx, msg, 0)
}

func (r *RedisMQ) ReceiveNormalMsg(ctx context.Context, callback mq.ConsumeCallBack) {
	r.consume(ctx, callback)
}

func (r *RedisMQ) SendDelayMsg(ctx context.Context, msg []byte, delay int64) (error, string) {
	return r.publish(ctx, msg, delay)
}

func (r *RedisMQ) ReceiveDelayMsg(ctx context.Context, callback mq.ConsumeCallBack) {
	r.consume(ctx, callback)
}

func (r *RedisMQ) publish(ctx context.Context, msg []byte, delay int64) (error, string) {
	values := map[string]interface{}{
		"message": msg,
	}
	if delay > 0 {
		timestamp := time.Now().Add(time.Duration(delay) * time.Millisecond).Format(time.RFC3339Nano)
		values["timestamp"] = timestamp
	}

	msgId, err := r.rdb.XAdd(ctx, &redis.XAddArgs{
		Stream:       r.Config.QueueName,
		ID:           "*",
		MaxLenApprox: 100000,
		Values:       values,
	}).Result()

	return err, msgId
}

func (r *RedisMQ) consume(ctx context.Context, callback mq.ConsumeCallBack) {
	go r.processFailAckMsg(ctx, callback)

	retryCount := 0
	r.rdb.XGroupCreateMkStream(ctx, r.Config.QueueName, r.Config.GroupName, "$")

	for {
		time.Sleep(time.Duration(2*retryCount) * time.Second)
		notify := true

		result, err := r.rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    r.Config.GroupName,
			Consumer: r.Config.ConsumerName,
			Streams:  []string{r.Config.QueueName, ">"},
			Block:    0,
		}).Result()

		if err != nil && err != redis.Nil {
			r.err <- fmt.Errorf("XReadGroup failed : %v", err)
			retryCount++
			continue
		}

		if len(result) <= 0 {
			retryCount++
			continue
		}

		for _, stream := range result {
			for _, message := range stream.Messages {
				timestamp, ok := message.Values["timestamp"].(string)
				if ok {
					deliveryTime, err := time.Parse(time.RFC3339Nano, timestamp)
					if err != nil {
						r.err <- fmt.Errorf("delay queue parse time failed: %v", err)
						continue
					}

					// 如果消息还没到投递时间，则重新加入延时队列
					if time.Now().Before(deliveryTime) {
						//重新入队的消息应从 consumer pending 和 stream 列表中移除
						r.rdb.XAck(ctx, r.Config.QueueName, r.Config.GroupName, message.ID).Result()
						r.rdb.XDel(ctx, r.Config.QueueName, message.ID)
						time.Sleep(deliveryTime.Sub(time.Now()))

						_, err := r.rdb.XAdd(ctx, &redis.XAddArgs{
							Stream: r.Config.QueueName,
							ID:     "*",
							Values: message.Values,
						}).Result()

						if err != nil {
							r.err <- fmt.Errorf("delay queue re-entry queue faild: %v", err)
						}

						notify = false
						continue
					}
				}

				if err := callback(ctx, message.ID, []byte(cast.ToString(message.Values["message"])), r.Config.ConsumerName); err != nil {
					r.err <- fmt.Errorf("callback exec failed, callbackName: %+v, callbackResult: %v", mq.GetFuncName(callback), err)
					continue
				}

				if _, err := r.rdb.XAck(ctx, r.Config.QueueName, r.Config.GroupName, message.ID).Result(); err != nil {
					r.err <- fmt.Errorf("consume ack faild: %v", err)
					continue
				}

				r.rdb.XDel(ctx, r.Config.QueueName, message.ID)
			}
		}

		retryCount = 0
		if notify {
			r.done <- notify
		}
	}

}

func (r *RedisMQ) processFailAckMsg(ctx context.Context, callback mq.ConsumeCallBack) {
	r.rdb.XGroupCreateMkStream(ctx, r.Config.QueueName, r.Config.GroupName, "0")

	for {
		<-r.done

		result, err := r.rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    r.Config.GroupName,
			Consumer: r.Config.ConsumerName,
			Streams:  []string{r.Config.QueueName, "0"},
			Block:    0,
		}).Result()

		if err != nil && err != redis.Nil {
			r.err <- fmt.Errorf("XReadGroup failed in pendings: %v", err)
			continue
		}

		if len(result) <= 0 {
			continue
		}

		for _, stream := range result {
			for _, message := range stream.Messages {
				if err := callback(ctx, message.ID, []byte(cast.ToString(message.Values["message"])), r.Config.ConsumerName); err != nil {
					r.err <- fmt.Errorf("callback exec failed in pendings, callbackName: %v, callbackResult: %v", mq.GetFuncName(callback), err)
					continue
				}

				if _, err := r.rdb.XAck(ctx, r.Config.QueueName, r.Config.GroupName, message.ID).Result(); err != nil {
					r.err <- fmt.Errorf("consume ack faild in pendings: %v", err)
					continue
				}

				r.rdb.XDel(ctx, r.Config.QueueName, message.ID)
			}
		}
	}
}

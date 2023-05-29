package rabbitmq

import (
	"context"
	"fmt"
	"strings"
	"time"

	"gin-api/pkg/hash"
	"gin-api/pkg/mq"
	amqp "github.com/rabbitmq/amqp091-go"
)

type WorkMode string

const (
	Worker    WorkMode = "worker"
	Topic     WorkMode = "topic"
	BroadCast WorkMode = "broadcast"
)

var (
	retryLimit           = 20
	normalExchangeDirect = "normal.exchange.direct" //直连交换机,用来投递单播消息
	normalExchangeTopic  = "normal.exchange.topic"  //主题交换机,用来投递topic消息
)

//确保 RabbitMQ 实现了 mq.Contracts
var _ mq.Contracts = (*RabbitMQ)(nil)

type Config struct {
	QueueName string
	Mode      WorkMode
	LogPrintf func(string, ...interface{})

	exchangeName       string
	exchangeKind       string //交换机类型(direct/topic)
	normalExchangeName string //交换机名称
	normalRoutingKey   string //路由名称
	deadExchangeName   string //死信交换机名称
	deadRoutingKey     string //死信路由名称
}

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	dsn     string
	err     chan error
	Config
}

func NewRabbitMQ(dsn string, queueConfig Config) *RabbitMQ {

	if queueConfig.QueueName == "" {
		panic("please provide the queue name")
	}

	if queueConfig.Mode != Worker && queueConfig.Mode != BroadCast && queueConfig.Mode != Topic {
		panic("please provide a legal working mode")
	}

	if queueConfig.LogPrintf == nil {
		queueConfig.LogPrintf = mq.Printf
	}

	r := &RabbitMQ{
		dsn: dsn,
		err: make(chan error),
	}

	return r.setConfig(queueConfig).connect()
}

func (r *RabbitMQ) setConfig(queueConfig Config) *RabbitMQ {
	if queueConfig.Mode == Worker || queueConfig.Mode == BroadCast {
		queueConfig.exchangeName = normalExchangeDirect
	} else if r.Mode == Topic {
		queueConfig.exchangeName = normalExchangeTopic
	}

	routingKey := addPrefix(queueConfig.QueueName)
	queueConfig.exchangeKind = getExchangeKind(queueConfig.exchangeName)
	queueConfig.normalExchangeName = queueConfig.exchangeName
	queueConfig.normalRoutingKey = routingKey
	queueConfig.deadExchangeName = replaceNormalToDead(queueConfig.exchangeName)
	queueConfig.deadRoutingKey = replaceNormalToDead(routingKey)
	r.Config = queueConfig
	return r
}

func (r *RabbitMQ) connect() *RabbitMQ {
	go r.log()
	if r.conn != nil && !r.conn.IsClosed() {
		return r
	}

	retryCount := 0
	for {
		var err error
		r.conn, err = amqp.Dial(r.dsn)
		if err != nil || r.conn == nil || r.conn.IsClosed() {
			if retryCount > retryLimit {
				panic("RabbitMQ Service Unavailable")
			}
			r.err <- fmt.Errorf("connect failed : %v", err)
			time.Sleep(time.Duration(retryCount) * time.Second)
			retryCount++
			continue
		}
		retryCount = 0
		return r
	}
}

func (r *RabbitMQ) log() {
	for {
		select {
		case err, ok := <-r.err:
			if !ok {
				return
			}
			r.LogPrintf("%v, queue_name: %v", err, r.QueueName)
		}
	}
}

func (r *RabbitMQ) destroy() {
	r.channel.Close()
	r.conn.Close()
}

//SendNormalMsg 发送消息，delay 为延时投递时间(单位毫秒)
func (r *RabbitMQ) SendNormalMsg(ctx context.Context, message []byte) (error, string) {
	return r.publish(ctx, message, 0)
}

//SendDelayMsg 发送延时消息，delay 为延时投递时间(单位毫秒)
func (r *RabbitMQ) SendDelayMsg(ctx context.Context, message []byte, delay int64) (error, string) {
	if delay <= 0 {
		panic("SendDelayMsg must provide delay ms")
	}
	return r.publish(ctx, message, delay)
}

//publish 发送消息
//参数说明： message 为需要投递的消息，delay 如果大于0则延时投递消息
//返回值说明： errChan 为消息投递中出现的错误
func (r *RabbitMQ) publish(ctx context.Context, message []byte, delay int64) (error, string) {
	defer r.destroy()
	if len(message) > 4*1024*1024 {
		return fmt.Errorf("message size cannot exceed 4M"), ""
	}

	//创建channel
	r.connect()
	channel, _ := r.conn.Channel()
	r.channel = channel

	//死信交换机属性
	argsQue := make(map[string]interface{})
	if delay > 0 {
		argsQue["x-message-ttl"] = delay                        //消息过期时间,毫秒
		argsQue["x-dead-letter-exchange"] = r.deadExchangeName  //指定死信交换机
		argsQue["x-dead-letter-routing-key"] = r.deadRoutingKey //指定死信routing-key
	}

	//创建队列,队列名称为空让 broker 自动生成
	q, err := r.channel.QueueDeclare(
		r.normalRoutingKey,
		true,
		false,
		false,
		false,
		argsQue)
	if err != nil {
		return fmt.Errorf("failed to declare queue, err:%v", err), ""
	}

	//队列绑定（将队列、routing-key、交换机三者绑定到一起）
	err = r.channel.QueueBind(
		q.Name,
		r.normalRoutingKey,
		r.normalExchangeName,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to QueueBind, err:%v", err), ""
	}

	msgId := r.GenMsgId(message)
	err = r.channel.PublishWithContext(
		ctx,
		r.normalExchangeName,
		r.normalRoutingKey, //routing-key
		false,
		false,
		amqp.Publishing{
			MessageId:   r.GenMsgId(message),
			ContentType: "text/plain",
			Body:        message,
		})
	if err != nil {
		return fmt.Errorf("failed to PublishWithContext, err:%v", err), ""
	}

	//非延时广播应删掉队列
	if r.Mode == BroadCast && delay <= 0 {
		r.channel.QueueDelete(q.Name, false, false, false)
	}

	return nil, msgId
}

//ReceiveNormalMsg 接收消息 ，callback 为消费者提供的回调，消息会写入到该回调中供消费者自行处理，callback() 没有错误则会自动ack，反之broker会再次投递。
func (r *RabbitMQ) ReceiveNormalMsg(ctx context.Context, callback mq.ConsumeCallBack) {
	r.consume(ctx, r.normalExchangeName, r.normalRoutingKey, callback)
}

//ReceiveDelayMsg 接收延时消息, callback 为消费者提供的回调，消息会写入到该回调中供消费者自行处理，callback() 没有错误则会自动ack，反之broker会再次投递。
func (r *RabbitMQ) ReceiveDelayMsg(ctx context.Context, callback mq.ConsumeCallBack) {
	r.consume(ctx, r.deadExchangeName, r.deadRoutingKey, callback)
}

func (r *RabbitMQ) GenMsgId(message []byte) string {
	return fmt.Sprintf("%v-%v-%v", removePrefix(r.normalRoutingKey), time.Now().Format("20060102150405"), hash.HashByMd5(string(message)))
}

//consume 消费消息，需要传入交换机、routing-key以及回调函数
func (r *RabbitMQ) consume(ctx context.Context, exchangeName, routingKey string, callback mq.ConsumeCallBack) {
	for {
		//创建 channel
		r.connect()
		channel, _ := r.conn.Channel()
		r.channel = channel
		notifyClose := r.conn.NotifyClose(make(chan *amqp.Error))

		//创建交换机
		err := r.channel.ExchangeDeclare(
			exchangeName,
			r.exchangeKind,
			true,
			false,
			false,
			false,
			nil,
		)

		if err != nil {
			r.err <- fmt.Errorf("failed to declare exchange, err:%v", err)
			continue
		}

		//手动设置队列名和自动生成队列名的区别.
		//手动设置：多worker队列名相同，所以会在worker间均衡派发消息
		//自动生成：多worker队列名不相同，消息会派发给每一位worker
		queueName := routingKey
		autoDelete := false
		if r.Mode == BroadCast {
			queueName = ""
			autoDelete = true
		}
		q, err := r.channel.QueueDeclare(
			queueName,
			true,
			autoDelete,
			false,
			false,
			nil,
		)

		if err != nil {
			r.err <- fmt.Errorf("failed to declare queue, err:%v", err)
			continue
		}

		//队列绑定（将队列、routing-key、交换机三者绑定到一起）
		err = r.channel.QueueBind(q.Name, routingKey, exchangeName, false, nil)
		if err != nil {
			r.err <- fmt.Errorf("failed to QueueBind, err:%v", err)
			continue
		}

		//负载均衡策略-根据负载量公平调度
		//if r.exchangeKind == "direct" && ! r.isBroadCast {
		//	err := r.channel.Qos(1, 0, false)
		//	if err != nil {
		//		logs.CtxError(ctx, "failed to Qos, err:%v", err)
		//		panic(fmt.Errorf("failed to Qos, err:%v", err))
		//	}
		//}

		//消费消息
		msgs, err := r.channel.Consume(
			q.Name,
			"",
			false,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			r.err <- fmt.Errorf("failed to ReceiveNormalMsg, err:%v", err)
			continue
		}

		quit := false
		for !quit {
			select {
			case <-ctx.Done():
				r.err <- fmt.Errorf("context done")
				r.destroy()
				return
			case e := <-notifyClose:
				r.err <- fmt.Errorf("receive notifyClose: %v", e)
				quit = true
			case msg, ok := <-msgs:
				if !ok {
					r.err <- fmt.Errorf("consume chan has been closed")
					quit = true
					continue
				}
				err := callback(ctx, msg.MessageId, msg.Body, nil)
				if err != nil {
					r.err <- fmt.Errorf("callback exec failed, callbackName:%v, callbackResult:%v", mq.GetFuncName(callback), err)
					continue
				}
				msg.Ack(true)
			}
		}

		r.destroy()
	}
}

func addPrefix(routingKey string) string {
	return fmt.Sprintf("normal.%v", routingKey)
}

func removePrefix(routingKey string) string {
	return strings.TrimPrefix(routingKey, "normal.")
}

func replaceNormalToDead(str string) string {
	if !strings.HasPrefix(str, "normal.") {
		return str
	}
	return strings.Replace(str, "normal.", "dead.", -1)
}

func getExchangeKind(exchangeName string) string {
	arr := strings.Split(exchangeName, ".")
	return arr[len(arr)-1]
}

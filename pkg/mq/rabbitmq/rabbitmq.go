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

var (
	DefaultConnectTimeOut = time.Second
	RecoverTimeOut        = 30 * time.Second
	RetryCount            = 0
	RetryLimit            = 20

	NormalExchangeDirect = "normal.exchange.direct" //直连交换机,用来投递单播消息
	NormalExchangeTopic  = "normal.exchange.topic"  //主题交换机,用来投递topic消息
)

//确保 RabbitMQ 实现了 mq.Contracts
var _ mq.Contracts = (*RabbitMQ)(nil)

type RabbitMQ struct {
	conn         *amqp.Connection
	channel      *amqp.Channel
	ExchangeKind string //交换机类型(direct/topic)

	NormalExchangeName string //交换机名称
	NormalRoutingKey   string //路由名称，routing_key必须是一个由.分隔开的词语列表。例如"erp.order.stat"，erp.order.sync","erp.bill.sync" 词语个数可以随意，但是不要超过255字节。

	DeadExchangeName string //死信交换机名称
	DeadRoutingKey   string //死信路由名称

	isBroadCast bool //是否为广播消息
	Dsn         string
}

func NewRabbitMq(dsn string) *RabbitMQ {
	r := &RabbitMQ{Dsn: dsn}
	return r.connect()
}

func (r *RabbitMQ) connect() *RabbitMQ {
	var err error

	if r.conn != nil && !r.conn.IsClosed() {
		return r
	}

	for {
		if RetryCount > RetryLimit && err != nil {
			RetryCount = 0
			panic(err)
		}
		r.conn, err = amqp.Dial(r.Dsn)
		if err != nil || r.conn == nil || r.conn.IsClosed() {
			fmt.Println("rabbit-mq重新连接中...")
			RetryCount++
			time.Sleep(time.Duration(RetryCount) * DefaultConnectTimeOut)
			continue
		} else {
			fmt.Println("已连接...")
			RetryCount = 0
			break
		}
	}
	return r
}

//UseWorker 使用 worker 模式
func (r *RabbitMQ) UseWorker(name string) *RabbitMQ {
	r.setExchange(name, NormalExchangeDirect)
	return r
}

//UseBroadCast 使用 broadcast 模式
func (r *RabbitMQ) UseBroadCast(name string) *RabbitMQ {
	r.setExchange(name, NormalExchangeDirect)
	r.isBroadCast = true
	return r
}

//UseTopic 使用 topic 模式
func (r *RabbitMQ) UseTopic(name string) *RabbitMQ {
	r.setExchange(name, NormalExchangeTopic)
	return r
}

func (r *RabbitMQ) setExchange(name, exchangeName string) *RabbitMQ {
	routingKey := addPrefix(name)
	r.ExchangeKind = getExchangeKind(exchangeName)
	r.NormalExchangeName = exchangeName
	r.NormalRoutingKey = routingKey
	r.DeadExchangeName = replaceNormalToDead(exchangeName)
	r.DeadRoutingKey = replaceNormalToDead(routingKey)

	return r
}

func (r *RabbitMQ) Destroy() {
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
		panic("SendDelayMsg must provide delay")
	}
	return r.publish(ctx, message, delay)
}

//publish 发送消息
//参数说明： message 为需要投递的消息，delay 如果大于0则延时投递消息
//返回值说明： errChan 为消息投递中出现的错误
func (r *RabbitMQ) publish(ctx context.Context, message []byte, delay int64) (error, string) {
	defer r.Destroy()

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
		argsQue["x-dead-letter-exchange"] = r.DeadExchangeName  //指定死信交换机
		argsQue["x-dead-letter-routing-key"] = r.DeadRoutingKey //指定死信routing-key
	}

	//创建队列,队列名称为空让 broker 自动生成
	q, err := r.channel.QueueDeclare(
		r.NormalRoutingKey,
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
		r.NormalRoutingKey,
		r.NormalExchangeName,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to QueueBind, err:%v", err), ""
	}

	msgId := r.GenMsgId(message)
	err = r.channel.PublishWithContext(
		ctx,
		r.NormalExchangeName,
		r.NormalRoutingKey, //routing-key
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
	if r.isBroadCast && delay <= 0 {
		r.channel.QueueDelete(q.Name, false, false, false)
	}

	return nil, msgId
}

//ReceiveNormalMsg 接收消息 ，fn 为消费者提供的回调，消息会写入到该回调中供消费者自行处理，fn() 没有错误则会自动ack，反之broker会再次投递。
func (r *RabbitMQ) ReceiveNormalMsg(ctx context.Context, fn mq.ConsumeCallBack) {
	r.consume(ctx, r.NormalExchangeName, r.NormalRoutingKey, fn)
}

//ReceiveDelayMsg 接收延时消息, fn 为消费者提供的回调，消息会写入到该回调中供消费者自行处理，fn() 没有错误则会自动ack，反之broker会再次投递。
func (r *RabbitMQ) ReceiveDelayMsg(ctx context.Context, fn mq.ConsumeCallBack) {
	r.consume(ctx, r.DeadExchangeName, r.DeadRoutingKey, fn)
}

func (r *RabbitMQ) GenMsgId(message []byte) string {
	return fmt.Sprintf("%v-%v-%v", removePrefix(r.NormalRoutingKey), time.Now().Format("20060102150405"), hash.HashByMd5(string(message)))
}

//consume 消费消息，需要传入交换机、routing-key以及回调函数
func (r *RabbitMQ) consume(ctx context.Context, exchangeName, routingKey string, fn mq.ConsumeCallBack) {

	//需保证consume不能崩溃，待服务恢复后，使consume从recover中再次拉起消费进程
	//@notice 如果一直连接不上，这里将有可能导致栈溢出
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("consume recover err:%v", err)
			time.Sleep(RecoverTimeOut)
			r.consume(ctx, exchangeName, routingKey, fn)
		}
	}()

	for {
		//创建 channel
		r.connect()
		channel, _ := r.conn.Channel()
		r.channel = channel
		notifyClose := r.conn.NotifyClose(make(chan *amqp.Error))

		//创建交换机
		err := r.channel.ExchangeDeclare(
			exchangeName,
			r.ExchangeKind,
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			fmt.Printf("failed to declare exchange, err:%v", err)
			panic(fmt.Errorf("failed to declare exchange, err:%v", err))
		}

		//手动设置队列名和自动生成队列名的区别.
		//手动设置：多worker队列名相同，所以会在worker间均衡派发消息
		//自动生成：多worker队列名不相同，消息会派发给每一位worker
		queueName := routingKey
		autoDelete := false
		if r.isBroadCast {
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
			fmt.Printf("failed to declare queue, err:%v", err)
			panic(fmt.Errorf("failed to declare queue, err:%v", err))
		}

		//队列绑定（将队列、routing-key、交换机三者绑定到一起）
		err = r.channel.QueueBind(q.Name, routingKey, exchangeName, false, nil)
		if err != nil {
			fmt.Printf("failed to QueueBind, err:%v", err)
			panic(fmt.Errorf("failed to QueueBind, err:%v", err))
		}

		//负载均衡策略-根据负载量公平调度
		//if r.ExchangeKind == "direct" && ! r.isBroadCast {
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
			fmt.Printf("failed to Consume, err:%v", err)
			panic(fmt.Errorf("failed to Consume, err:%v", err))
		}

		forever := make(chan bool)
		go func() {
			for {
				select {
				case e := <-notifyClose:
					fmt.Printf("notifyClose receive err:%v", e.Error())
					forever <- false
					goto loop
				case msg := <-msgs:
					if len(msg.MessageId) > 0 {
						err := fn(ctx, msg.MessageId, msg.Body, nil)
						if err == nil {
							msg.Ack(true)
						} else {
							fmt.Printf("CallBack fn fail, err:%v", err)
						}
					}
				}
			}
		loop:
		}()
		<-forever
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

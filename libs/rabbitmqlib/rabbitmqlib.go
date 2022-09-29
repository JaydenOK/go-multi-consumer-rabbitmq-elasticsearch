package rabbitmqlib

import (
	"app/utils"
	"fmt"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
	"time"
)

type RabbitMQ struct {
	connection   *amqp.Connection
	channel      *amqp.Channel
	url          string
	exchangeName string
	queueName    string
	routeKey     string
	bindKey      string
	exchangeType string
}

func NewRabbitMQ() *RabbitMQ {
	host := viper.GetString("rabbitmq.host")
	port := viper.GetString("rabbitmq.port")
	username := viper.GetString("rabbitmq.username")
	password := viper.GetString("rabbitmq.password")
	url := "amqp://" + username + ":" + password + "@" + host + ":" + port + "/"
	rabbitMQ := &RabbitMQ{
		url: url,
	}
	_ = rabbitMQ.InitRabbitMQ()
	return rabbitMQ
}

// 关闭
func (rabbitMQ *RabbitMQ) close() {
	defer rabbitMQ.connection.Close()
	defer rabbitMQ.channel.Close()
}

// 直接发送消息，direct模式，路由设置与队列名一致
// direct – 直接匹配
// 直接匹配交换器
// 用于支持路由模式（Routing）
// 直接匹配交换器会对比路由键和绑定键，如果路由键和绑定键完全相同，则把消息转发到绑定键所对应的队列中。
func (rabbitMQ *RabbitMQ) SendMessage(message string, exchangeName string, queueName string) error {
	var v interface{}
	rabbitMQ.exchangeName = exchangeName
	rabbitMQ.queueName = queueName
	rabbitMQ.routeKey = queueName
	rabbitMQ.exchangeType = "direct"
	//name:交换器的名称，对应图中exchangeName。
	//kind:也叫作type，表示交换器的类型。有四种常用类型：direct、fanout、topic、headers。
	//durable:是否持久化，true表示是。持久化表示会把交换器的配置存盘，当RMQ Server重启后，会自动加载交换器。
	//autoDelete:是否自动删除，true表示是。至少有一条绑定才可以触发自动删除，当所有绑定都与交换器解绑后，会自动删除此交换器。
	//internal:是否为内部，true表示是。客户端无法直接发送msg到内部交换器，只有交换器可以发送msg到内部交换器。
	//noWait:是否非阻塞，true表示是。阻塞：表示创建交换器的请求发送后，阻塞等待RMQ Server返回信息。非阻塞：不会阻塞等待RMQ Server的返回信息，而RMQ Server也不会返回信息。（不推荐使用）
	//args:直接写nil
	if err := rabbitMQ.channel.ExchangeDeclare(
		rabbitMQ.exchangeName,
		rabbitMQ.exchangeType,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		fmt.Println("创建exchange异常：", err.Error())
		rabbitMQ.close()
		v = "创建exchange异常"
		panic(v)
	}
	//不存在创建
	//name：队列名称
	//durable：是否持久化，true为是。持久化会把队列存盘，服务器重启后，不会丢失队列以及队列内的信息。（注：1、不丢失是相对的，如果宕机时有消息没来得及存盘，还是会丢失的。2、存盘影响性能。）
	//autoDelete：是否自动删除，true为是。至少有一个消费者连接到队列时才可以触发。当所有消费者都断开时，队列会自动删除。
	//exclusive：是否设置排他，true为是。如果设置为排他，则队列仅对首次声明他的连接可见，并在连接断开时自动删除。（注意，这里说的是连接不是信道，相同连接不同信道是可见的）。
	//nowait：是否非阻塞，true表示是。阻塞：表示创建交换器的请求发送后，阻塞等待RMQ Server返回信息。非阻塞：不会阻塞等待RMQ Server的返回信息，而RMQ Server也不会返回信息。（不推荐使用）
	//args：直接写nil
	if _, err := rabbitMQ.channel.QueueDeclare(
		rabbitMQ.queueName,
		true,
		false,
		false,
		true,
		nil,
	); err != nil {
		rabbitMQ.close()
		v = "rabbitmq连接异常"
		panic(v)
	}

	if err := rabbitMQ.channel.QueueBind(
		rabbitMQ.queueName,
		rabbitMQ.routeKey,
		rabbitMQ.exchangeName,
		true,
		nil,
	); err != nil {
		rabbitMQ.close()
		v = "绑定队列异常"
		panic(v)
	}
	//exchange：要发送到的交换机名称，对应图中exchangeName。
	//key：路由键，对应图中RoutingKey。
	//mandatory：直接false，不建议使用，后面有专门章节讲解。
	//immediate ：直接false，不建议使用，后面有专门章节讲解。
	//msg：要发送的消息，msg对应一个Publishing结构，Publishing结构里面有很多参数
	err := rabbitMQ.channel.Publish(
		rabbitMQ.exchangeName,
		rabbitMQ.routeKey,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: 2,
			ContentType:  "application/json",
			Body:         []byte(message),
			Timestamp:    time.Now(),
		},
	)
	return err
}

// topic – 模式匹配
// 与直接匹配相对应，可以用一些模式来代替字符串的完全匹配。
// 规则：
// 以 ‘.’ 来分割单词。
// ‘#’ 表示一个或多个单词。
// ‘*’ 表示一个单词。
// 如：
// RoutingKey为：
// aaa.bbb.ccc
// BindingKey可以为：
// *.bbb.ccc
// aaa.#
func (rabbitMQ *RabbitMQ) sendMessageTopic(message string, exchangeName string, queueName string) {
	rabbitMQ.exchangeName = exchangeName
	rabbitMQ.queueName = queueName
	//不存在创建
	//name：队列名称
	//durable：是否持久化，true为是。持久化会把队列存盘，服务器重启后，不会丢失队列以及队列内的信息。（注：1、不丢失是相对的，如果宕机时有消息没来得及存盘，还是会丢失的。2、存盘影响性能。）
	//autoDelete：是否自动删除，true为是。至少有一个消费者连接到队列时才可以触发。当所有消费者都断开时，队列会自动删除。
	//exclusive：是否设置排他，true为是。如果设置为排他，则队列仅对首次声明他的连接可见，并在连接断开时自动删除。（注意，这里说的是连接不是信道，相同连接不同信道是可见的）。
	//nowait：是否非阻塞，true表示是。阻塞：表示创建交换器的请求发送后，阻塞等待RMQ Server返回信息。非阻塞：不会阻塞等待RMQ Server的返回信息，而RMQ Server也不会返回信息。（不推荐使用）
	//args：直接写nil

}

// fanout – 扇出型
// 扇出交换器
// 用于支持发布、订阅模式（pub/sub）
// 交换器把消息转发到所有与之绑定的队列中。
// 扇出类型交换器会屏蔽掉路由键、绑定键的作用。
func (rabbitMQ *RabbitMQ) sendMessageFanOut(message string, exchangeName string, queueName string) {
	rabbitMQ.exchangeName = exchangeName
	rabbitMQ.queueName = queueName
}

// 初始化rabbitMQ，连接及创建通道，连接失败，返回Error错误
func (rabbitMQ *RabbitMQ) InitRabbitMQ() error {
	var err error
	var v interface{}
	if rabbitMQ.connection, err = amqp.Dial(rabbitMQ.url); err != nil {
		fmt.Println("rabbitMQ连接异常：", err.Error())
		v = "rabbitMQ连接异常，执行终止"
		panic(v)
	}
	//2,创建通道
	if rabbitMQ.channel, err = rabbitMQ.connection.Channel(); err != nil {
		if err := rabbitMQ.connection.Close(); err != nil {
			fmt.Println(err)
		}
		fmt.Println("rabbitMQ创建channel异常：", err.Error())
		v = "rabbitMQ创建channel异常，执行终止"
		panic(v)
	}
	return err
}

// 初始化消费服务
// 直接发送消息，direct模式，路由设置与队列名一致
func (rabbitMQ *RabbitMQ) ConsumeInit(exchangeName string, queueName string) error {
	rabbitMQ.exchangeName = exchangeName
	rabbitMQ.queueName = queueName
	rabbitMQ.bindKey = queueName
	rabbitMQ.exchangeType = "direct"
	if err := rabbitMQ.channel.ExchangeDeclare(
		rabbitMQ.exchangeName,
		rabbitMQ.exchangeType,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		fmt.Println("绑定exchange异常：", err.Error())
		rabbitMQ.close()
		panic(utils.StringToInterface("绑定exchange异常" + err.Error()))
	}

	if _, err := rabbitMQ.channel.QueueDeclare(
		rabbitMQ.queueName,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		fmt.Println("绑定queue异常：", err.Error())
		rabbitMQ.close()
		panic(utils.StringToInterface("创建queue异常" + err.Error()))
	}

	if err := rabbitMQ.channel.QueueBind(
		rabbitMQ.queueName,
		rabbitMQ.bindKey,
		rabbitMQ.exchangeName,
		false,
		nil,
	); err != nil {
		fmt.Println("绑定queue异常：", err.Error())
		rabbitMQ.close()
		panic(utils.StringToInterface("绑定queue异常" + err.Error()))
	}

	return nil
}

// 消费
func (rabbitMQ *RabbitMQ) Consume(consumerTag string) (<-chan amqp.Delivery, error) {
	//只读信道，接收队列数据
	var delivery <-chan amqp.Delivery
	var err error
	if delivery, err = rabbitMQ.channel.Consume(
		rabbitMQ.queueName,
		consumerTag,
		false,
		false,
		false,
		false,
		nil,
	); err != nil {
		rabbitMQ.close()
		return nil, err
	}
	return delivery, nil
}

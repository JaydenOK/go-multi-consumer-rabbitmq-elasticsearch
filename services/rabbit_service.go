package services

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	connection   *amqp.Connection
	channel      *amqp.Channel
	url          string
	exchangeName string
	queueName    string
	routeKey     string
	bindKey      string
}

func NewRabbitMQ() *RabbitMQ {
	host := viper.GetString("rabbitmq.host")
	port := viper.GetString("mysql.port")
	username := viper.GetString("rabbitmq.username")
	password := viper.GetString("rabbitmq.password")
	url := "amqp://" + username + ":" + password + "@" + host + ":" + port + "/"
	rabbitMQ := &RabbitMQ{
		url: url,
	}
	rabbitMQ.InitRabbitMQ()
	return rabbitMQ
}

// direct模式
// direct – 直接匹配
// 直接匹配交换器
// 用于支持路由模式（Routing）
// 直接匹配交换器会对比路由键和绑定键，如果路由键和绑定键完全相同，则把消息转发到绑定键所对应的队列中。
func (rabbitMQ *RabbitMQ) sendMessage(message string, exchangeName string, queueName string) {
	rabbitMQ.exchangeName = exchangeName
	rabbitMQ.queueName = queueName
	//不存在创建
	//name：队列名称
	//durable：是否持久化，true为是。持久化会把队列存盘，服务器重启后，不会丢失队列以及队列内的信息。（注：1、不丢失是相对的，如果宕机时有消息没来得及存盘，还是会丢失的。2、存盘影响性能。）
	//autoDelete：是否自动删除，true为是。至少有一个消费者连接到队列时才可以触发。当所有消费者都断开时，队列会自动删除。
	//exclusive：是否设置排他，true为是。如果设置为排他，则队列仅对首次声明他的连接可见，并在连接断开时自动删除。（注意，这里说的是连接不是信道，相同连接不同信道是可见的）。
	//nowait：是否非阻塞，true表示是。阻塞：表示创建交换器的请求发送后，阻塞等待RMQ Server返回信息。非阻塞：不会阻塞等待RMQ Server的返回信息，而RMQ Server也不会返回信息。（不推荐使用）
	//args：直接写nil，没研究过，不解释。
	queue, err := rabbitMQ.channel.QueueDeclare(rabbitMQ.queueName, true, false, true, true, nil)
	if err != nil {
		fmt.Println("rabbitmq连接异常：", err.Error())
		panic(err)
	}
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
	//args：直接写nil，没研究过，不解释。
	queue, err := rabbitMQ.channel.QueueDeclare(rabbitMQ.queueName, true, false, true, true, nil)
	if err != nil {
		fmt.Println("rabbitmq连接异常：", err.Error())
		panic(err)
	}
}

// fanout – 扇出型
// 扇出交换器
// 用于支持发布、订阅模式（pub/sub）
// 交换器把消息转发到所有与之绑定的队列中。
// 扇出类型交换器会屏蔽掉路由键、绑定键的作用。
func (rabbitMQ *RabbitMQ) sendMessageFanOut(message string, exchangeName string, queueName string) {
	rabbitMQ.exchangeName = exchangeName
	rabbitMQ.queueName = queueName
	//不存在创建
	//name：队列名称
	//durable：是否持久化，true为是。持久化会把队列存盘，服务器重启后，不会丢失队列以及队列内的信息。（注：1、不丢失是相对的，如果宕机时有消息没来得及存盘，还是会丢失的。2、存盘影响性能。）
	//autoDelete：是否自动删除，true为是。至少有一个消费者连接到队列时才可以触发。当所有消费者都断开时，队列会自动删除。
	//exclusive：是否设置排他，true为是。如果设置为排他，则队列仅对首次声明他的连接可见，并在连接断开时自动删除。（注意，这里说的是连接不是信道，相同连接不同信道是可见的）。
	//nowait：是否非阻塞，true表示是。阻塞：表示创建交换器的请求发送后，阻塞等待RMQ Server返回信息。非阻塞：不会阻塞等待RMQ Server的返回信息，而RMQ Server也不会返回信息。（不推荐使用）
	//args：直接写nil，没研究过，不解释。
	queue, err := rabbitMQ.channel.QueueDeclare(rabbitMQ.queueName, true, false, true, true, nil)
	if err != nil {
		fmt.Println("rabbitmq连接异常：", err.Error())
		panic(err)
	}
}

// 初始化rabbitMQ，连接及创建通道
func (rabbitMQ *RabbitMQ) InitRabbitMQ() {
	defer func() {
		_ = rabbitMQ.connection.Close()
	}()
	var err error
	rabbitMQ.connection, err = amqp.Dial(rabbitMQ.url)
	if err != nil {
		fmt.Println("rabbitmq连接异常：", err.Error())
		panic(err)
	}
	//2,创建通道
	rabbitMQ.channel, err = rabbitMQ.connection.Channel()
	if err != nil {
		fmt.Println("rabbitmq创建通道异常：", err.Error())
		panic(err)
	}
}

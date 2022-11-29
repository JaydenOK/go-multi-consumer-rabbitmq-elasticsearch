package tasks

import (
	"app/libs/rabbitmqlib"
	"github.com/streadway/amqp"
)

type ConsumerUint struct {
	//监听交换机名称
	exchangeName string
	//监听队列名称
	queueName string
	//启动的消费者数量，默认一个
	taskNum int
	//rabbitMQ连接对象
	rabbitMQ *rabbitmqlib.RabbitMQ
}

func (consumerUint *ConsumerUint) getName() string {
	return "unit_consumer"
}

func (consumerUint *ConsumerUint) start() error {

	return nil
}

// 业务逻辑
func (consumerUint *ConsumerUint) handler(delivery amqp.Delivery) error {

	return nil
}

func (consumerUint *ConsumerUint) pushOrderToElasticSearch(orderIds []string) error {
	return nil
}

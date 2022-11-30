package tasks

import (
	"app/libs/rabbitmqlib"
	"fmt"
	"github.com/streadway/amqp"
)

type StockConsumer struct {
	exchangeName string
	queueName    string
	taskNum      int
	rabbitMQ     *rabbitmqlib.RabbitMQ
}

func (stockConsumer *StockConsumer) getName() string {
	return "stock_consumer"
}

func (stockConsumer *StockConsumer) getRabbitMQ() *rabbitmqlib.RabbitMQ {
	return stockConsumer.rabbitMQ
}

func (stockConsumer *StockConsumer) start() error {
	var err error
	stockConsumer.rabbitMQ = rabbitmqlib.NewRabbitMQ()
	stockConsumer.rabbitMQ.SetConsumerConfig(
		stockConsumer.exchangeName,
		stockConsumer.queueName,
		stockConsumer.taskNum,
		stockConsumer.handler,
	)
	if err = stockConsumer.rabbitMQ.ConsumeStart(); err != nil {
		fmt.Println("启动失败：", err.Error())
	}
	return nil
}

// 业务逻辑
func (stockConsumer *StockConsumer) handler(delivery amqp.Delivery) error {
	fmt.Printf("stock_consumer[%s]接收到数据：%s", delivery.ConsumerTag, string(delivery.Body))
	return nil
}

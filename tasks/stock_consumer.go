package tasks

import (
	"app/libs/rabbitmqlib"
	"app/utils"
	"fmt"
	"github.com/streadway/amqp"
	"strconv"
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

func (stockConsumer *StockConsumer) start() error {
	if stockConsumer.taskNum == 0 {
		stockConsumer.taskNum = 1
	}
	for index := 0; index < stockConsumer.taskNum; index++ {
		if err := stockConsumer.rabbitMQ.ConsumeInit(
			stockConsumer.exchangeName,
			stockConsumer.queueName,
		); err != nil {
			fmt.Println(utils.StringToInterface(err.Error()))
		}
		consumerTag := stockConsumer.queueName + strconv.Itoa(index)
		delivery, err := stockConsumer.rabbitMQ.Consume(consumerTag)
		if err != nil {
			fmt.Println("消费异常：", utils.StringToInterface(err.Error()))
		}
		go func() {
			for d := range delivery {
				if err := stockConsumer.handle(d); err == nil {
					_ = d.Ack(false)
				} else {
					_ = d.Ack(false)
				}
			}
		}()
	}

	return nil
}

// 业务逻辑
func (stockConsumer *StockConsumer) handle(delivery amqp.Delivery) error {
	fmt.Printf("stock_consumer[%s]接收到数据：%s", delivery.ConsumerTag, string(delivery.Body))
	return nil
}

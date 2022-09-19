package tasks

import (
	"app/libs/rabbitmqlib"
	"app/utils"
	"fmt"
	"github.com/streadway/amqp"
	"strconv"
)

type OrderConsumer struct {
	//监听交换机
	exchangeName string
	//监听队列
	queueName string
	//启动的消费者数量
	taskNum int
	//rabbitMQ连接对象
	rabbitMQ *rabbitmqlib.RabbitMQ
}

func (orderConsumer *OrderConsumer) getName() string {
	return "order_consumer"
}

func (orderConsumer *OrderConsumer) start() error {
	if orderConsumer.taskNum == 0 {
		orderConsumer.taskNum = 1
	}
	for index := 0; index < orderConsumer.taskNum; index++ {
		go func(i int) {
			if err := orderConsumer.rabbitMQ.ConsumeInit(
				orderConsumer.exchangeName,
				orderConsumer.queueName,
			); err != nil {
				fmt.Println(utils.StringToInterface(err.Error()))
			}
			delivery, err := orderConsumer.rabbitMQ.Consume(orderConsumer.queueName + strconv.Itoa(i))
			if err != nil {
				fmt.Println("消费异常：", utils.StringToInterface(err.Error()))
			}
			go func() {
				for d := range delivery {
					//只读，没数据时，处于阻塞状态
					if err := orderConsumer.handle(d); err == nil {
						//false 单条确认，true多条确认
						_ = d.Ack(false)
					} else {
						//当 requeue 为真时，请求服务器将此消息传递给不同的
						//消费者。 如果不可能或 requeue 为 false，则消息将是
						//丢弃或交付到服务器配置的死信队列。
						//
						//此方法不得用于选择或重新排队客户端希望的消息
						//不去处理，而是通知服务端客户端无能力
						//在这个时候处理这个消息。
						//_ = d.Nack(false, false)
						_ = d.Ack(false)
					}
				}
			}()
		}(index)
	}

	//if err := orderConsumer.rabbitMQ.ConsumeInit(
	//	orderConsumer.exchangeName,
	//	orderConsumer.queueName,
	//); err != nil {
	//	fmt.Println(utils.StringToInterface(err.Error()))
	//}
	//delivery, err := orderConsumer.rabbitMQ.Consume()
	//if err != nil {
	//	fmt.Println("消费异常：", utils.StringToInterface(err.Error()))
	//}
	//go func() {
	//	for d := range delivery {
	//		//只读，没数据时，处于阻塞状态
	//		if err := orderConsumer.handle(d.Body); err == nil {
	//			//false 单条确认，true多条确认
	//			_ = d.Ack(false)
	//		} else {
	//			//当 requeue 为真时，请求服务器将此消息传递给不同的
	//			//消费者。 如果不可能或 requeue 为 false，则消息将是
	//			//丢弃或交付到服务器配置的死信队列。
	//			//
	//			//此方法不得用于选择或重新排队客户端希望的消息
	//			//不去处理，而是通知服务端客户端无能力
	//			//在这个时候处理这个消息。
	//			//_ = d.Nack(false, false)
	//			_ = d.Ack(false)
	//		}
	//	}
	//}()
	return nil
}

// 业务逻辑
func (orderConsumer *OrderConsumer) handle(delivery amqp.Delivery) error {

	fmt.Printf("order_consumer[%s]接收到数据：%s", delivery.ConsumerTag, string(delivery.Body))
	//@todo do something

	return nil
}

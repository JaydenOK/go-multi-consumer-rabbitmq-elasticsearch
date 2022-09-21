package tasks

import (
	"app/libs/elasticsearchlib"
	"app/libs/mysqllib"
	"app/libs/rabbitmqlib"
	"app/models"
	"app/utils"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/streadway/amqp"
	"log"
	"strconv"
	"sync"
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
	var mqData map[string]interface{}
	if err := json.Unmarshal(delivery.Body, &mqData); err != nil {
		fmt.Println("队列json数据解析异常:", string(delivery.Body), err.Error())
		return err
	}
	key := "order_id"
	orderId, ok := mqData[key]
	if ok != true {
		msg := fmt.Sprintf("推送的数据错误，%s 不存在", key)
		fmt.Println(msg)
		return errors.New(msg)
	}
	if err := orderConsumer.pushOrderToElasticSearch(orderId.(string)); err != nil {
		msg := fmt.Sprintf("推送es异常，%s", orderId)
		fmt.Println(msg)
		return errors.New(msg)
	}
	return nil
}

func (orderConsumer *OrderConsumer) pushOrderToElasticSearch(orderId string) error {
	//@todo do something 将订单变更信息推送es
	mysqlClient := mysqllib.GetMysqlClient()
	var orderModel models.OrderModel
	mysqlClient.Where("order_id = ?", orderId).Find(&orderModel)
	if orderModel.OrderId == "" {
		msg := fmt.Sprintf("推送的数据错误，order_id不存在")
		fmt.Println(msg)
		return errors.New(msg)
	}
	var orderModels []models.OrderModel
	orderModels = append(orderModels, orderModel)
	esClient := elasticsearchlib.GetClient()
	var wg sync.WaitGroup
	for i, orderModel := range orderModels {
		wg.Add(1)
		go func(orderModel models.OrderModel) {
			defer wg.Done()
			// 构建请求json数据
			data, err := json.Marshal(orderModel)
			if err != nil {
				log.Fatalf("Error marshaling document: %s", err)
			}
			//设置请求对象。
			//索引创建或更新索引中的文档。
			request := esapi.IndexRequest{
				Index:      "order",
				DocumentID: strconv.Itoa(int(orderModel.Id) + 1),
				Body:       bytes.NewReader(data),
				Refresh:    "true",
			}

			// 与客户端执行请求。
			res, err := request.Do(context.Background(), esClient)
			if err != nil {
				log.Fatalf("Error getting response: %s", err)
			}
			defer res.Body.Close()

			if res.IsError() {
				log.Printf("[%s] Error indexing document ID=%d", res.Status(), i+1)
			} else {
				// 将响应反序列化为映射。
				var r map[string]interface{}
				if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
					log.Printf("Error parsing the response body: %s", err)
				} else {
					// 打印响应状态和索引文档版本。
					log.Printf("[%s] %s; version=%d", res.Status(), r["result"], int(r["_version"].(float64)))
				}
			}
		}(orderModel)
	}
	wg.Wait()
	return nil
}

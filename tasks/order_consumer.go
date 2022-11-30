package tasks

import (
	"app/libs/elasticsearchlib"
	"app/libs/mysqllib"
	"app/libs/rabbitmqlib"
	"app/models"
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
	//监听交换机名称
	exchangeName string
	//监听队列名称
	queueName string
	//启动的消费者数量，默认一个
	taskNum int
	//rabbitMQ连接对象
	rabbitMQ *rabbitmqlib.RabbitMQ
}

func (orderConsumer *OrderConsumer) getName() string {
	return "order_consumer"
}

func (orderConsumer *OrderConsumer) getRabbitMQ() *rabbitmqlib.RabbitMQ {
	return orderConsumer.rabbitMQ
}

func (orderConsumer *OrderConsumer) start() error {
	var err error
	orderConsumer.rabbitMQ = rabbitmqlib.NewRabbitMQ()
	orderConsumer.rabbitMQ.SetConsumerConfig(
		orderConsumer.exchangeName,
		orderConsumer.queueName,
		orderConsumer.taskNum,
		orderConsumer.handler,
	)
	if err = orderConsumer.rabbitMQ.ConsumeStart(); err != nil {
		fmt.Println("启动失败：", err.Error())
	}
	return nil
}

// 业务逻辑
func (orderConsumer *OrderConsumer) handler(delivery amqp.Delivery) error {
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
	orderIds := []string{orderId.(string)}
	if err := orderConsumer.pushOrderToElasticSearch(orderIds); err != nil {
		msg := fmt.Sprintf("推送es异常，%s", orderIds)
		fmt.Println(msg)
		return errors.New(msg)
	}
	return nil
}

func (orderConsumer *OrderConsumer) pushOrderToElasticSearch(orderIds []string) error {
	//@todo do something 将订单变更信息推送es
	mysqlClient := mysqllib.GetMysqlClient()
	var orderModels []models.OrderModel
	mysqlClient.Where("order_id in ?", orderIds).Find(&orderModels)
	if orderModels == nil {
		msg := fmt.Sprintf("推送的数据错误，订单没找到:%v", orderIds)
		fmt.Println(msg)
		return errors.New(msg)
	}
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
				DocumentID: strconv.Itoa(int(orderModel.Id)),
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

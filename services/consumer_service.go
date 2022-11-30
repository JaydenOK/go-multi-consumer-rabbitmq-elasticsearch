package services

import (
	"app/tasks"
	"github.com/gin-gonic/gin"
)

type ConsumerService struct {
}

func (service ConsumerService) StopAll(ctx *gin.Context) interface{} {
	//exchangeName := ctx.PostForm("exchange_name")
	taskManager := &tasks.TaskManager{}
	taskManager.StopAll()
	return nil
}

func (service ConsumerService) StartConsumer(ctx *gin.Context) interface{} {
	//exchangeName := ctx.PostForm("exchange_name")
	//queneName := ctx.PostForm("queue_name")
	//bindRoute := ctx.PostForm("bind_route")
	//rabbitMQ := rabbitmqlib.NewRabbitMQ()

	return nil
}

func (service ConsumerService) StopConsumer(ctx *gin.Context) interface{} {
	exchangeName := ctx.PostForm("exchange_name")
	//todo 不属于同一个taskManager，不能控制已启动的任务，改用信号方式
	taskManager := &tasks.TaskManager{}
	//taskManager.StopTask(&tasks.ConsumerUint{
	//	exchangeName: constants.EventStockChange,
	//	queueName:    constants.QueueStockChange,
	//	taskNum:      1,
	//})

	taskManager.StopConsumer(exchangeName)
	return nil
}

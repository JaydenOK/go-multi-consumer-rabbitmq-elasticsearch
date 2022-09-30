package tasks

import (
	"app/constants"
	"app/utils"
	"fmt"
)

type TaskManager struct {
	taskPool []Consumer
}

func Run() {
	taskManager := &TaskManager{}
	//添加订单任务程序，
	taskManager.AddTask(&OrderConsumer{
		exchangeName: constants.EventOrderChange,
		queueName:    constants.QueueOrderChange,
		taskNum:      3,
	})

	//添加库存任务程序
	taskManager.AddTask(&StockConsumer{
		exchangeName: constants.EventStockChange,
		queueName:    constants.QueueStockChange,
		taskNum:      1,
	})

	taskManager.start()
}

// 添加任务类型
func (taskManager *TaskManager) AddTask(consumer Consumer) {
	isExist := false
	for _, task := range taskManager.taskPool {
		if consumer.getName() == task.getName() {
			isExist = true
			break
		}
	}
	if !isExist {
		fmt.Println("添加消费者：", consumer.getName())
		taskManager.taskPool = append(taskManager.taskPool, consumer)
	}
}

// 启动
func (taskManager *TaskManager) start() {
	for _, task := range taskManager.taskPool {
		if err := task.start(); err != nil {
			fmt.Println(utils.StringToInterface(err.Error()))
		}
	}
}

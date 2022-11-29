package tasks

import (
	"app/libs/rabbitmqlib"
	"sync"
)

type ConsumerProcess struct {
	lock         sync.Mutex
	consumerList []rabbitmqlib.RabbitMQ
}

func NewConsumerProcess() *ConsumerProcess {
	return &ConsumerProcess{}
}

func (ths *ConsumerProcess) isExist() bool {
	isExist := false
	//for _, consumer := range ths.consumerList {
	//
	//}
	return isExist
}

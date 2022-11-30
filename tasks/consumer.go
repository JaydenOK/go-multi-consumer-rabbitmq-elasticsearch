package tasks

import (
	"app/libs/rabbitmqlib"
	"github.com/streadway/amqp"
)

type Consumer interface {
	getName() string
	start() error
	handler(amqp.Delivery) error
	getRabbitMQ() *rabbitmqlib.RabbitMQ
}

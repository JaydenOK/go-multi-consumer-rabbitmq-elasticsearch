package tasks

import "github.com/streadway/amqp"

type Consumer interface {
	getName() string
	start() error
	handle(amqp.Delivery) error
}

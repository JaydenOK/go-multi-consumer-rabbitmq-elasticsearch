package tasks

import "github.com/streadway/amqp"

type Consumer interface {
	getName() string
	start() error
	handler(amqp.Delivery) error
}

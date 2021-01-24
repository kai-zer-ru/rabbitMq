package rabbitMqGolang

import "github.com/streadway/amqp"

type RabbitMQ struct {
	Host               string
	Port               int64
	UserName           string
	Password           string
	rabbitMqConnection *amqp.Connection
	isConnected        bool
}

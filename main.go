package rabbitMqGolang

import (
	"fmt"
	"github.com/streadway/amqp"
)

func (r *RabbitMQ) Connect() error {
	if r.Host == "" {
		r.Host = "localhost"
	}
	if r.Port == 0 {
		r.Port = 5672
	}
	if r.UserName == "" {
		r.UserName = "guest"
	}
	if r.Password == "" {
		r.Password = "guest"
	}
	var err error
	r.rabbitMqConnection, err = amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%v/", r.UserName, r.Password, r.Host, r.Port))
	if err == nil {
		r.isConnected = true
	}
	return err
}

func (r *RabbitMQ) Close() error {
	if !r.isConnected {
		return nil
	}
	r.isConnected = false
	return r.rabbitMqConnection.Close()
}

func (r *RabbitMQ) Send(queue string, data []byte) error {
	if !r.isConnected {
		return fmt.Errorf("RabbitMQ is not connected")
	}
	ch, err := r.rabbitMqConnection.Channel()
	if err != nil {
		return err
	}
	defer func() {
		_ = ch.Close()
	}()
	q, err := ch.QueueDeclare(queue, false, false, false, false, nil)
	if err != nil {
		return err
	}
	err = ch.Publish("", q.Name, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        data,
	})
	return err
}

func (r *RabbitMQ) Listen(channels ...string) (map[string]<-chan amqp.Delivery, error) {
	ch, err := r.rabbitMqConnection.Channel()
	if err != nil {
		return nil, err
	}
	queues := make(map[string]<-chan amqp.Delivery)
	for _, channel := range channels {
		q, err := ch.QueueDeclare(channel, false, false, false, false, nil)
		if err != nil {
			return nil, err
		}
		messages, err := ch.Consume(q.Name, "", true, false, false, false, nil)
		if err != nil {
			return nil, err
		}
		queues[channel] = messages
	}
	return queues, nil
}
package rabbitmq

import (
	"os"
	"sync"

	"github.com/streadway/amqp"
)

var (
	amqpConnection     *amqp.Connection
	amqpConnectionOnce sync.Once
)

func AmqpConn() *amqp.Connection {
	amqpConnectionOnce.Do(func() {
		var err error
		amqpConnection, err = amqp.Dial(os.Getenv("RABBITMQ_URL"))
		if err != nil {
			panic(err)
		}
	})
	return amqpConnection
}

func AmqpClose() {
	if amqpConnection != nil {
		amqpConnection.Close()
	}
}

func AmqpChannel() (*amqp.Channel, error) {
	return AmqpConn().Channel()
}

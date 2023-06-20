package rabbitmq

import (
	"errors"
	"os"
	"sync"
	"time"

	"github.com/streadway/amqp"
	"github.com/tinkler/mqttadmin/pkg/logger"
)

var (
	RetryDuration      = time.Second * 5
	amqpConnection     *amqp.Connection
	amqpConnectionOnce sync.Once
)

func AmqpConn() *amqp.Connection {
	amqpConnectionOnce.Do(func() {
		dialRabbiMQ()
	})
	return amqpConnection
}

func dialRabbiMQ() {
	var err error
	amqpConnection, err = amqp.Dial(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		panic(err)
	}
}

func AmqpClose() {
	if amqpConnection != nil {
		amqpConnection.Close()
	}
}

func AmqpChannel() (*amqp.Channel, error) {
	ch, err := AmqpConn().Channel()
	if errors.Is(amqp.ErrClosed, err) {
		logger.Warn("retry to dial")
		// retry to dial
	RETYR:
		dialRabbiMQ()
		ch, err = AmqpConn().Channel()
		if err != nil {
			logger.Error(err)
			time.Sleep(RetryDuration)
			goto RETYR
		}
	}

	return ch, err
}

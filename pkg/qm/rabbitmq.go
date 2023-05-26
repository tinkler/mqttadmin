package qm

import (
	"errors"
	"sync"
	"time"

	"github.com/streadway/amqp"
	"github.com/tinkler/mqttadmin/pkg/rabbitmq"
)

type RabbitMQ struct {
	sync.Mutex
	channels map[string]*amqp.Channel
}

func (d *RabbitMQ) Publish(channel string, message string) (string, error) {
	d.Lock()
	defer d.Unlock()
	if _, ok := d.channels[channel]; !ok {
		ch, err := rabbitmq.AmqpChannel()
		if err != nil {
			return "", err
		}
		d.channels[channel] = ch
	}
	ch := d.channels[channel]
	correlationId := channel + "-qm"
	messageId := GetNextVal()
	p := amqp.Publishing{
		ContentType:   "text/plain",
		Body:          []byte(message),
		ReplyTo:       channel + "-reply",
		CorrelationId: correlationId,
		MessageId:     messageId,
	}
	acks := ch.NotifyPublish(make(chan amqp.Confirmation))
	err := ch.Publish(
		"",      // exchange
		channel, // routing key
		false,   // mandatory
		false,   // immediate
		p,
	)
	if err != nil {
		return "", err
	}

	channelReplyQueue, err := ch.QueueDeclare(
		channel+"-reply", // name
		false,            // durable
		false,            // delete when unused
		false,            // exclusive
		false,            // no-wait
		nil,              // arguments
	)
	if err != nil {
		return "", err
	}

	consumer := channel + "-consumer"
	msgs, err := ch.Consume(
		channelReplyQueue.Name, // Queue name
		consumer,               // Consumer tag
		true,                   // Auto-ack
		false,                  // Exclusive
		false,                  // No-local
		false,                  // No-wait
		nil,                    // Arguments
	)
	if err != nil {
		return "", err
	}
	defer ch.Cancel(consumer, false)

	for {
		select {
		case msg := <-msgs:
			// Wait for the response
			if msg.CorrelationId == "" {
				// Ignore messages that don't have a correlation ID set
				continue
			}
			if msg.CorrelationId == correlationId {
				if msg.MessageId == p.MessageId {
					return string(msg.Body), nil
				}
			}
		case ack := <-acks:
			if ack.Ack {
				return "", nil
			}
			return "", errors.New("processing failed")
		case <-time.NewTimer(time.Second * 5).C:
			return "", errors.New("processing timeout")
		}
	}

}

func (d *RabbitMQ) Close() error {
	d.Lock()
	defer d.Unlock()
	for _, ch := range d.channels {
		ch.Close()
	}
	return nil
}

func NewRabbitMQ() *RabbitMQ {
	return &RabbitMQ{
		channels: make(map[string]*amqp.Channel),
	}
}

package qm

import (
	"errors"
	"sync"
	"time"

	"github.com/streadway/amqp"
	"github.com/tinkler/mqttadmin/pkg/logger"
	"github.com/tinkler/mqttadmin/pkg/rabbitmq"
)

type RabbitMQ struct {
	sync.Mutex
	channels map[string]*amqp.Channel
}

func (d *RabbitMQ) Publish(channel string, message string) error {
	d.Lock()
	defer d.Unlock()
	if _, ok := d.channels[channel]; !ok {
		ch, err := rabbitmq.AmqpChannel()
		if err != nil {
			logger.Error(err)
			return ErrPublish
		}
		d.channels[channel] = ch
	}
	ch := d.channels[channel]

	err := ch.Confirm(false)
	if err != nil {
		logger.Error(err)
		return ErrPublish
	}
	p := amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(message),
	}
	err = ch.Publish(
		"",      // exchange
		channel, // routing key
		false,   // mandatory
		false,   // immediate
		p,
	)
	if err != nil {
		logger.Error(err)
		return ErrPublish
	}

	acks := ch.NotifyPublish(make(chan amqp.Confirmation, 1))

	timeout := time.NewTimer(PublishTimeout)
	for {
		select {
		case ack, ok := <-acks:
			if !ok {
				logger.Error("ask channel closed")
				return errors.New("ask channel closed")
			}
			if ack.Ack {
				return nil
			} else {
				continue
			}
		case <-timeout.C:
			logger.Error("Processing timeout")
			return errors.New("processing timeout")
		default:
			time.Sleep(time.Millisecond * 100)
		}
	}
}

func (d *RabbitMQ) PublishAndReceive(channel string, message string) (string, error) {
	d.Lock()
	defer d.Unlock()
	if _, ok := d.channels[channel]; !ok {
		ch, err := rabbitmq.AmqpChannel()
		if err != nil {
			logger.Error(err)
			return "", ErrPublish
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
	err := ch.Confirm(false)
	if err != nil {
		logger.Error(err)
		return "", ErrPublish
	}
	err = ch.Publish(
		"",      // exchange
		channel, // routing key
		false,   // mandatory
		false,   // immediate
		p,
	)
	if err != nil {
		logger.Error(err)
		return "", ErrPublish
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
		logger.Error(err)
		return "", ErrPublish
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
		logger.Error(err)
		return "", ErrPublish
	}
	defer ch.Cancel(consumer, false)

	timeout := time.NewTimer(PublishTimeout)
	for {
		select {
		case msg, ok := <-msgs:
			if !ok {
				logger.Error("replay channel closed")
				return "", errors.New("replay channel closed")
			}
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
		case <-timeout.C:
			logger.Error("Processing timeout")
			return "", errors.New("processing timeout")
		default:
			time.Sleep(time.Millisecond * 100)
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

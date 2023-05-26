/*
A queue manager
*/
package qm

import "sync"

type QueueManager struct {
	Driver Driver
}

var (
	qmInstance     *QueueManager
	qmInstanceOnce sync.Once
)

func Qm() *QueueManager {
	qmInstanceOnce.Do(func() {
		qmInstance = &QueueManager{
			Driver: NewRabbitMQ(),
		}
	})
	return qmInstance
}

func Publish(channel string, message string) (string, error) {
	return Qm().Driver.Publish(channel, message)
}

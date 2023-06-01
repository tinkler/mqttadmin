package qm

import "io"

type Driver interface {
	Publish(channel string, message string) (err error)
	PublishAndReceive(channel string, message string) (reply string, err error)
	io.Closer
}

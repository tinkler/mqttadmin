package gs

import "golang.org/x/net/context"

// Stream is a generic interface for stream
// It is used in the following way:
// @stream(bidi) or @stream(c2s) @ stream(s2c)
// It is used for sending or receiving T message.
// when the Recv function is triggered,besides receiving message,
// arguments of struct will also be updated.
type Stream[T any] interface {
	Context() context.Context
	Send(m T) error
	Recv() (T, error)
}

// Null is a placeholder type for gs.Stream[*gs.Null]
type Null struct{}

// NullStream is a placeholder type for gs.Stream[*gs.Null]
// It is useful for generating a stream that does not send or receive any message.
// but only to trigger the stream function.
// when the Recv function is triggered,
// arguments of struct will be updated.
// It is used in the following way:
// @stream(bidi) or @stream(c2s) @ stream(s2c)
type NullStream Stream[*Null]

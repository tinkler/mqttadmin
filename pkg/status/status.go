// Package status provides a simple way to get the status of a service.
package status

import "fmt"

type Status struct {
	Code    int32
	Message string
}

func Ok(msg string) *Status {
	return &Status{Code: 200, Message: msg}
}

func New(c int32, msg string) *Status {
	return &Status{Code: c, Message: msg}
}

func Newf(c int32, format string, a ...interface{}) *Status {
	return New(c, fmt.Sprintf(format, a...))
}

func (s *Status) Err() error {
	if s.Code == 200 {
		return nil
	}
	return s
}

func (s *Status) Error() string {
	return s.String()
}

func (s *Status) String() string {
	return fmt.Sprintf("rpc error: code = %d desc = %s", s.Code, s.Message)
}

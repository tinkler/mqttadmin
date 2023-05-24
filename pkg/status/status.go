// Package status provides a simple way to get the status of a service.
package status

import "fmt"

type Status struct {
	Code      int32
	Message   string
	CnMessage string
}

func Ok(msg string) *Status {
	return &Status{Code: 200, Message: msg}
}

func OkCn(msg, cnMsg string) *Status {
	return &Status{Code: 200, Message: msg, CnMessage: cnMsg}
}

func New(c int32, msg string) *Status {
	return &Status{Code: c, Message: msg}
}

func Newf(c int32, format string, a ...interface{}) *Status {
	return New(c, fmt.Sprintf(format, a...))
}

func NewCn(c int32, msg string, cnMsg string) *Status {
	return &Status{Code: c, Message: msg, CnMessage: cnMsg}
}

func NewCnf(c int32, format, formatCn string, a ...interface{}) *Status {
	return NewCn(c, fmt.Sprintf(format, a...), fmt.Sprintf(formatCn, a...))
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

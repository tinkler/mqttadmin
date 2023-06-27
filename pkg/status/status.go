// Package status provides a simple way to get the status of a service.
package status

import (
	"fmt"
	"net/http"

	"google.golang.org/grpc/codes"
)

// CodeType is the type of the code.
// 0 is HTTP status codes as registered with IANA.
// 1 is gRPC status codes
type Status struct {
	CodeType  uint8
	Code      int32
	Message   string
	CnMessage string
}

func Ok(msg string) *Status {
	return &Status{Code: 200, CodeType: 0, Message: msg}
}

func OkCn(msg, cnMsg string) *Status {
	return &Status{Code: 200, CodeType: 0, Message: msg, CnMessage: cnMsg}
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

func (s *Status) Type(codeType uint8) *Status {
	s.CodeType = codeType
	return s
}

func (s *Status) ConvToType(codeType uint8) *Status {
	if s.CodeType == codeType {
		return s
	}
	if s.CodeType == 0 && codeType == 1 {
		s.CodeType = 1
		switch s.Code {
		case http.StatusOK:
			s.Code = int32(codes.OK)
		}
		return s
	}
	if s.CodeType == 1 && codeType == 0 {
		s.CodeType = 0
		switch s.Code {
		case int32(codes.OK):
			s.Code = http.StatusOK
		}
		return s
	}
	return s
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

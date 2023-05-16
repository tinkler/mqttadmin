package errz

import (
	errzpb "github.com/tinkler/mqttadmin/errz/v1"
	"github.com/tinkler/mqttadmin/pkg/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Error struct {
	s *status.Status
}

func (e *Error) Err() error {
	return e.s.Err()
}

func (e *Error) Error() string {
	return e.s.String()
}

func Err(err error) error {
	if err == nil {
		return nil
	}
	switch e := err.(type) {
	case *errzpb.AuthError:
		st := status.New(codes.Unauthenticated, codes.Unauthenticated.String())
		ds, err := st.WithDetails(e)
		if err != nil {
			return &Error{s: st}
		}
		return &Error{s: ds}
	case *errzpb.ServerError:
		st := status.New(codes.Internal, codes.Internal.String())
		ds, err := st.WithDetails(e)
		if err != nil {
			return &Error{s: st}
		}
		return &Error{s: ds}
	case *errzpb.ValidateError:
		st := status.New(codes.InvalidArgument, codes.InvalidArgument.String())
		ds, err := st.WithDetails(e)
		if err != nil {
			return &Error{s: st}
		}
		return &Error{s: ds}
	case *Error:
		return e
	default:
		logger.Debug("%T", err)
		return e
	}
}

func ErrDb(err error) error {
	if err == nil {
		return nil
	}
	logger.Error(err)
	return &errzpb.ServerError{
		Type:      errzpb.ServerError_DB,
		Message:   "数据库错误",
		EnMessage: "Database error",
	}
}

func ErrInternal(err error) error {
	if err == nil {
		return nil
	}
	logger.Error(err)
	return &errzpb.ServerError{
		Type:      errzpb.ServerError_INTERNAL,
		Message:   "内部错误",
		EnMessage: "Internal error",
	}
}

func ErrInternalM(message string, enMessage string) error {
	if message == "" {
		return nil
	}
	logger.Error(message)
	return &errzpb.ServerError{
		Type:      errzpb.ServerError_INTERNAL,
		Message:   message,
		EnMessage: enMessage,
	}
}

func ErrVd(err error) error {
	if err == nil {
		return nil
	}
	logger.Error(err)
	switch e := err.(type) {
	case *errzpb.ValidateError:
		return e
	default:
		return &errzpb.ValidateError{}
	}
}

func ErrVdM(field string, message string, enMessage string) error {
	logger.Debug("非法输入:%s -> %s", field, message)
	return &errzpb.ValidateError{
		Field:     field,
		Message:   message,
		EnMessage: enMessage,
	}
}

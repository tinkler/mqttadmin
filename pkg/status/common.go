package status

import (
	"context"
	"errors"

	"github.com/tinkler/mqttadmin/pkg/logger"
	"google.golang.org/grpc/codes"
)

var (
	ErrUnauthorized   = NewCn(401, "unauthorized", "未授权")
	ErrForbidden      = NewCn(403, "forbidden", "禁止访问")
	ErrInternalServer = NewCn(500, "internal server error", "服务器内部错误")
	ErrCanceled       = NewCn(int32(codes.Canceled), "canceled", "调用取消").Type(1)
)

func StatusUnauthorized() error {
	logger.Error("未授权")
	return ErrUnauthorized
}
func StatusForbidden() error {
	logger.Error("禁止访问")
	return ErrForbidden
}
func StatusInternalServer(args ...interface{}) error {
	for i := 0; i < len(args); i++ {
		if e, ok := args[i].(error); ok {
			switch v := e.(type) {
			case *Status:
				return v
			default:
				if errors.Is(context.Canceled, e) {
					logger.Warn(e)
				} else {
					logger.Error(e)
				}

			}
		}
	}
	if len(args) == 0 {
		logger.Error("服务器内部错误")
	}
	return ErrInternalServer
}

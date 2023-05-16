package errz

import (
	errzpb "github.com/tinkler/mqttadmin/errz/v1"
	"github.com/tinkler/mqttadmin/pkg/logger"
)

func ErrAuth(typ errzpb.AuthError_ErrorType) error {
	var message [2]string
	switch typ {
	case errzpb.AuthError_TOKEN_EXPIRED:
		message = [2]string{"密钥过期", "token expired"}
	case errzpb.AuthError_TOKEN_INVALID:
		message = [2]string{"密钥非法", "token invalid"}
	case errzpb.AuthError_TOKEN_NOT_FOUND:
		message = [2]string{"密钥为空", "token empty"}
	case errzpb.AuthError_DEVICE_NOT_MATCH:
		message = [2]string{"手机设备不匹配", "device not match"}
	case errzpb.AuthError_DEVICE_NOT_FOUND:
		message = [2]string{"手机设备不存在", "device not found"}
	case errzpb.AuthError_PASSWORD_INVALID:
		message = [2]string{"密码错误", "password invalid"}
	default:
		return nil
	}
	logger.Debug("验证错误:%s", message[0])
	return &errzpb.AuthError{
		Type:      typ,
		Message:   message[0],
		EnMessage: message[1],
	}
}

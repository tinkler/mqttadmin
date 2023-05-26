package status

import "github.com/tinkler/mqttadmin/pkg/logger"

func StatusUnauthorized() error {
	logger.Error("未授权")
	return NewCn(401, "unauthorized", "未授权")
}
func StatusForbidden() error {
	logger.Error("禁止访问")
	return NewCn(403, "forbidden", "禁止访问")
}
func StatusInternalServer() error {
	logger.Error("服务器内部错误")
	return NewCn(500, "internal server error", "服务器内部错误")
}

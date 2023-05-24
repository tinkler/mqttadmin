package status

var (
	StatusUnauthorized   = NewCn(401, "unauthorized", "未授权")
	StatusForbidden      = NewCn(403, "forbidden", "禁止访问")
	StatusInternalServer = NewCn(500, "internal server error", "服务器内部错误")
)

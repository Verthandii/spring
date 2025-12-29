package ierrors

var (
	InternalServerError = NewStd(500, "服务器内部错误: %s")
	Unauthorized        = NewStd(401, "登录过期")
)

package ierrors

// Info 获取错误信息
func Info(err error) (int, string) {
	if err == nil {
		return _em.InternalServerError, ""
	}
	if e, ok := err.(*E); ok {
		var msg string
		if e.args != nil {
			msg = e.String()
		} else {
			msg = e.Error()
		}
		return e.code, msg
	}
	return _em.InternalServerError, err.Error()
}

// Code read code from _em
func Code(err error) int {
	if err == nil {
		return 0
	}
	if e, ok := err.(*E); ok {
		return e.code
	} else {
		return _em.InternalServerError
	}
}

// Message 给开发者看的错误信息
func Message(err error) string {
	return err.Error()
}

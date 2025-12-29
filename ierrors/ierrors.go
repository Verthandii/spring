package ierrors

import (
	"fmt"
)

// E implement Error interface
type E struct {
	code   int
	stderr bool
	msg    string
	args   []interface{}
}

// New 添加一个不会标准输出的 error
func New(code int, msg string) *E {
	return newError(code, false, msg)
}

// NewStd 添加一个会标准输出的 error
func NewStd(code int, msg string) *E {
	return newError(code, true, msg)
}

// newError create a error
func newError(code int, stderr bool, msg string) *E {
	e := &E{
		code:   code,
		stderr: stderr,
		msg:    msg,
	}
	_em.add(e)
	return e
}

func (e *E) String() string {
	if e.args != nil {
		return fmt.Sprintf(e.msg, e.args...)
	}
	return ""
}

// Error return code in string form
func (e *E) Error() string {
	return e.msg
}

// Format 格式化提示
func (e *E) Format(args ...interface{}) *E {
	clone := e.Clone()
	if clone.args != nil {
		clone.msg = fmt.Sprintf(e.msg, e.args)
	}
	if err, ok := args[0].(error); ok {
		clone.args = []interface{}{err.Error()}
		return &clone
	}
	clone.args = args
	return &clone
}

// Replace 替换原有模板
func (e *E) Replace(args ...interface{}) *E {
	clone := e.Clone()
	clone.msg = "%v"
	return clone.Format(args)
}

func (e *E) Clone() E {
	return E{
		code:   e.code,
		stderr: e.stderr,
		msg:    e.msg,
		args:   e.args,
	}
}

func (e *E) GetArgs() []interface{} {
	return e.args
}

func (e *E) SetMsg(msg string) {
	e.msg = msg
}

func (e *E) SetArgs(args []interface{}) {
	e.args = args
}

// List 获取错误码列表,用户判断错误码是否配置完成
func List() ([]int, map[int]string) {
	source := _em.m
	list := make(map[int]string, len(source))
	var ids []int
	for _, item := range source {
		list[item.code] = item.msg
		ids = append(ids, item.code)
	}
	return ids, list
}

package ierrors

import (
	"fmt"
	"sync"
)

var _em = initemap()

func initemap() *emap {
	return &emap{
		m:                   make(map[int]*E),
		mutex:               &sync.RWMutex{},
		InternalServerError: 500,
	}
}

type emap struct {
	m map[int]*E // 仅用来判断是否有重复的 code

	mutex               *sync.RWMutex
	InternalServerError int // 服务器内部错误 错误码 默认500
}

func (em *emap) add(e *E) {
	em.mutex.Lock()
	defer em.mutex.Unlock()

	if _, ok := em.m[e.code]; ok {
		panic(fmt.Sprintf("ecode error [%d] has exist", e.code))
	}

	em.m[e.code] = e
}

// SetInternalServerError 设置自定义的服务器内部错误
func SetInternalServerError(e *E) {
	InternalServerError = e
	_em.InternalServerError = Code(e)
}

// GetInternalServerErrorCode 获取服务器内部错误码
func GetInternalServerErrorCode() int {
	return _em.InternalServerError
}

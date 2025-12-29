package ihttp

import (
	"net/http"
)

var (
	_empty = make(map[string]any)
)

type response interface {
	setCode(code int)
	setData(data any)
	setMsg(msg string)
	reset()
}

type upperResponse struct {
	Code int
	Data any
	Msg  string
}

func (u *upperResponse) setCode(code int) {
	u.Code = code
}

func (u *upperResponse) setData(data any) {
	u.Data = data
}

func (u *upperResponse) setMsg(msg string) {
	u.Msg = msg
}

func (u *upperResponse) reset() {
	u.Code = 0
	u.Data = nil
	u.Msg = ""
}

type lowerResponse struct {
	Code int    `json:"code"`
	Data any    `json:"data"`
	Msg  string `json:"msg"`
}

func (l *lowerResponse) setCode(code int) {
	l.Code = code
}

func (l *lowerResponse) setData(data any) {
	l.Data = data
}

func (l *lowerResponse) setMsg(msg string) {
	l.Msg = msg
}

func (l *lowerResponse) reset() {
	l.Code = 0
	l.Data = nil
	l.Msg = ""
}

type SuccessEncodeFunc func(http.ResponseWriter, *http.Request, any)

type ErrorEncodeFunc func(http.ResponseWriter, *http.Request, error, any, int)

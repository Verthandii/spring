package ilogger

var _builtinPairs []interface{}

func Register(key string, value Valuer) {
	_builtinPairs = append(_builtinPairs, key, value)
}

package utils

import "golang.org/x/sync/singleflight"

var sf = singleflight.Group{}

func SF(key string, f func() (any, error)) (any, error, bool) {
	return sf.Do(key, f)
}

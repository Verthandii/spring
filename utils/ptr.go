package utils

func Ptr[T any](v T) *T {
	return &v
}

// Deref 安全地从不确定是否为nil的指针中获取(nil默认零值)值而不是 panic
func Deref[T any](in *T) T {
	if in == nil {
		var zero T
		return zero
	}
	return *in
}

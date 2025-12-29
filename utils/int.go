package utils

// BoolToInt64 bool->int64
func BoolToInt64(b bool) int64 {
	if b {
		return 1
	}
	return 0
}

func Int64ToBool(i int64) bool {
	return i > 0
}

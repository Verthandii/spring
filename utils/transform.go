package utils

import (
	"fmt"
	"strconv"
)

func ToInt(v any) int {
	switch t := v.(type) {
	case int:
		return t
	case int8:
		return int(t)
	case int16:
		return int(t)
	case int32:
		return int(t)
	case int64:
		return int(t)
	case bool:
		if t {
			return 1
		}
		return 0
	case float32:
		return int(t)
	case float64:
		return int(t)
	case uint8:
		return int(t)
	case uint16:
		return int(t)
	case uint32:
		return int(t)
	case uint64:
		return int(t)
	case string:
		i, err := strconv.Atoi(t)
		if err != nil {
			return 0
		}
		return i
	default:
		return 0
	}
}

func ToString(v any) string {
	switch t := v.(type) {
	case int:
		return strconv.Itoa(t)
	case int8:
		return fmt.Sprintf("%d", t)
	case int16:
		return fmt.Sprintf("%d", t)
	case int32:
		return fmt.Sprintf("%d", t)
	case int64:
		return strconv.FormatInt(t, 10)
	case bool:
		if t {
			return "true"
		}
		return "false"
	case float32:
		return strconv.FormatFloat(float64(t), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 32)
	case uint8:
		return fmt.Sprintf("%d", t)
	case uint16:
		return fmt.Sprintf("%d", t)
	case uint32:
		return fmt.Sprintf("%d", t)
	case uint64:
		return strconv.FormatUint(t, 10)
	case string:
		return t
	default:
		return ""
	}
}

func ToFloat64(v any) float64 {
	switch t := v.(type) {
	case int:
		return float64(t)
	case int8:
		return float64(t)
	case int16:
		return float64(t)
	case int32:
		return float64(t)
	case int64:
		return float64(t)
	case bool:
		if t {
			return 1
		}
		return 0
	case float32:
		return float64(t)
	case float64:
		return t
	case uint8:
		return float64(t)
	case uint16:
		return float64(t)
	case uint32:
		return float64(t)
	case uint64:
		return float64(t)
	case string:
		tt, err := strconv.ParseFloat(t, 64)
		if err != nil {
			return 0
		}
		return tt
	default:
		return 0
	}
}

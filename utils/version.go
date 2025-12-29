package utils

import (
	"strconv"
	"strings"
)

// IsGreaterThan 判断 v1 是否大于 v2
func IsGreaterThan(v1, v2 string) bool {
	return compareVersions(v1, v2) > 0
}

// IsGreaterThanOrEqual 判断 v1 是否大于等于 v2
func IsGreaterThanOrEqual(v1, v2 string) bool {
	return compareVersions(v1, v2) >= 0
}

// IsLessThan 判断 v1 是否小于 v2
func IsLessThan(v1, v2 string) bool {
	return compareVersions(v1, v2) < 0
}

// IsLessThanOrEqual 判断 v1 是否小于等于 v2
func IsLessThanOrEqual(v1, v2 string) bool {
	return compareVersions(v1, v2) <= 0
}

// compareVersions 比较两个版本号
func compareVersions(v1, v2 string) int {
	parts1 := strings.Split(v1, ".")
	parts2 := strings.Split(v2, ".")
	maxLength := maxInt(len(parts1), len(parts2))
	for i := 0; i < maxLength; i++ {
		val1 := 0
		if i < len(parts1) {
			val1, _ = strconv.Atoi(parts1[i])
		}
		val2 := 0
		if i < len(parts2) {
			val2, _ = strconv.Atoi(parts2[i])
		}
		if val1 > val2 {
			return 1
		} else if val1 < val2 {
			return -1
		}
	}
	return 0
}

// 获取最大长度
func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

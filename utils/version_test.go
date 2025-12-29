package utils

import (
	"testing"
)

func TestVersion(t *testing.T) {
	// 测试用例
	tests := []struct {
		v1     string
		v2     string
		v3     func(string, string) bool
		result bool
	}{
		{"1.2.3", "1.2.4", IsGreaterThan, false},
		{"1.2.3", "1.2.4", IsGreaterThanOrEqual, false},
		{"1.2.3", "1.2.4", IsLessThan, true},
		{"1.2.3", "1.2.4", IsLessThanOrEqual, true},
		{"1.10.1", "1.2.9", IsGreaterThan, true},
		{"1.10.1", "1.2.9", IsGreaterThanOrEqual, true},
		{"1.10.1", "1.2.9", IsLessThan, false},
		{"1.10.1", "1.2.9", IsLessThanOrEqual, false},
		{"1.0.0", "1.0.0", IsGreaterThan, false},
		{"1.0.0", "1.0.0", IsGreaterThanOrEqual, true},
		{"1.0.0", "1.0.0", IsLessThan, false},
		{"1.0.0", "1.0.0", IsLessThanOrEqual, true},
		{"2.0", "1.9.9", IsGreaterThan, true},
		{"2.0", "1.9.9", IsGreaterThanOrEqual, true},
		{"2.0", "1.9.9", IsLessThan, false},
		{"2.0", "1.9.9", IsLessThanOrEqual, false},
		{"1.2", "1.2.0", IsGreaterThan, false},
		{"1.2", "1.2.0", IsGreaterThanOrEqual, true},
		{"1.2", "1.2.0", IsLessThan, false},
		{"1.2", "1.2.0", IsLessThanOrEqual, true},
		{"2.0", "10.0.9", IsGreaterThan, false},
		{"2.0", "10.0.9", IsGreaterThanOrEqual, false},
		{"2.0", "10.0.9", IsLessThan, true},
		{"2.0", "10.0.9", IsLessThanOrEqual, true},
	}
	for _, test := range tests {
		if test.v3(test.v1, test.v2) != test.result {
			t.Fatalf("版本 %s ? %s: %v", test.v1, test.v2, test.result)
		}
	}
}

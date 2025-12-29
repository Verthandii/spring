package utils

import (
	"fmt"
	"testing"
)

func TestCurrentMonth(t *testing.T) {
	start := MonthStar()
	end := MonthEnd()
	t.Logf("当月开始日期:%s 本月结束日期:%s", start.Format("20060102"), end.Format("20060102"))

	ba, ea := CuMonth()
	t.Logf("当月开始时间:%d 本月结束时间:%d", ba, ea)
}

func TestIDCardSafe(t *testing.T) {
	for _, id := range []string{"360402198611133850", ""} {
		safe := IDCardSafe(id)
		t.Logf("safe = %s", safe)
		t.Logf("len(safe) = %d", len(safe))
	}
}

func TestIDCardSafeMainlandAndNotMainland(t *testing.T) {
	t.Logf("safe = %s", IDCardSafeMainlandAndNotMainland("360402198611133850", 0))
	t.Logf("safe = %s", IDCardSafeMainlandAndNotMainland("360402198611133850", 1))
	t.Logf("safe = %s", IDCardSafeMainlandAndNotMainland("A1234567", 2))
	t.Logf("safe = %s", IDCardSafeMainlandAndNotMainland("A1234567", 3))
}

func TestEncodeToShiftJIS(t *testing.T) {
	productName := "中文A商abcd商商品名称QA商品名称QA商品名称QA商品名称QA商品名称QA商品名称QA商品名称QA商品名称QA商品名称QA商品名称QA商品名称QA"
	for i := 0; i < 5; i++ {
		productName += productName
	}
	jis, err := EncodeToShiftJISLength(productName, 80)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(len(jis))
	t.Log(jis)
	utf8, err := ShiftJISToUTF8([]byte(jis))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(utf8)
}

func TestEmailSafe(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		// 长度=1
		{"a@qq.com", "a***@qq.com"},
		// 长度=2
		{"ab@gmail.com", "a***b@gmail.com"},
		// 长度=3
		{"abc@163.com", "a***bc@163.com"},
		// 长度=4
		{"abcd@yahoo.com", "ab***cd@yahoo.com"},
		// 长度=5
		{"abcde@domain.com", "ab***de@domain.com"},
		// 长度>5
		{"username@mail.com", "us***me@mail.com"},
		{"helloworld@test.com", "he***ld@test.com"},
		// 无效邮箱格式
		{"notanemail", "notanemail"},
		{"missing@domain@parts.com", "missing@domain@parts.com"},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			got := EmailSafeINTL(tc.input)
			if got != tc.expected {
				t.Errorf("EmailSafe(%q) = %q, want %q", tc.input, got, tc.expected)
			} else {
				fmt.Printf("✓ PASS: %-25s => %s\n", tc.input, got)
			}
		})
	}
}

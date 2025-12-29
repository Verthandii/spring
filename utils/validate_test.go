package utils

import (
	"testing"
)

func TestValIDCard(t *testing.T) {
	for _, s := range []string{"360402198611133850", "110101199003075293"} {
		if ValIDCard(s) != nil {
			t.Fatal("ValIDCard Error")
		}
	}
	for _, s := range []string{
		"310105202304082810",
		"310105202304082811",
		"310105202304082812",
		"310105202304082813",
		"310105202304082815",
		"310105202304082816",
		"310105202304082817",
		"310105202304082818",
		"310105202304082819",
		"310105202304082810",
		"31010520230408281X",
	} {
		if ValIDCard(s) == nil {
			t.Fatal("ValIDCard Error")
		}
	}
}

func TestCalculateIDCardChecksum(t *testing.T) {
	for _, s := range []string{"360402198611133850", "110101199003075293"} {
		sum := CalculateIDCardChecksum(s)
		if string(sum) != s[17:] {
			t.Fatal("IdNum CheckSum Wrong")
		}
	}

	for _, s := range []string{"310105202304082814"} {
		sum := CalculateIDCardChecksum(s)
		if string(sum) != s[17:] {
			t.Fatal("IdNum CheckSum Wrong")
		}
	}

	for _, s := range []string{
		"310105202304082810",
		"310105202304082811",
		"310105202304082812",
		"310105202304082813",
		"310105202304082815",
		"310105202304082816",
		"310105202304082817",
		"310105202304082818",
		"310105202304082819",
		"310105202304082810",
		"31010520230408281X",
	} {
		sum := CalculateIDCardChecksum(s)
		if string(sum) == s[17:] {
			t.Fatal("IdNum CheckSum Wrong")
		}
	}
}

func TestValBirthDate(t *testing.T) {
	// 测试 ValBirthDate 函数
	dates := []string{"20230605", "20240229", "19901215", "30000101", "19901301", "20230631"}
	for _, date := range dates {
		err := ValBirthDate(date)
		if err != nil {
			t.Logf("日期 %s 无效: %s\n", date, err)
		} else {
			t.Logf("日期 %s 有效\n", date)
		}
	}
}

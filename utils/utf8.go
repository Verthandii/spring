package utils

func IsChinese(r rune) bool {
	return r >= 0x4E00 && r <= 0x9FA5
}

func IsEnglish(r rune) bool {
	return r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z'
}

func IsNumber(r rune) bool {
	return r >= '0' && r <= '9'
}

package utils

import (
	"errors"
	"regexp"
	"time"
)

var (
	errDateFormat   = errors.New("日期格式错误:yyyymmdd")
	errMobileFormat = errors.New("手机号码格式错误")
	errFormat       = errors.New("格式有误")
	errNickname     = errors.New("昵称格式有误")
)

// ValDate 验证日期格式
func ValDate(s string) (err error) {
	matched, _ := regexp.MatchString("^\\d{4}\\d{2}\\d{2}$", s)
	if !matched {
		err = errDateFormat
	}
	return
}

// ValMobile 验证中国手机号
func ValMobile(s string) error {
	matched, _ := regexp.MatchString(`^(1\d{10})$`, s)
	if !matched {
		return errMobileFormat
	}
	return nil
}

// ValRealName 验证真实姓名(包括少数名族)
func ValRealName(s string) error {
	matched, err := regexp.MatchString(`^.{2,20}$`, s)
	if err != nil {
		return err
	}
	if !matched {
		return errFormat
	}
	return nil
}

// ValIDCard 验证身份证
func ValIDCard(s string) error {
	matched, err := regexp.MatchString(`(^\d{18}$)|(^\d{17}(\d|X|x)$)`, s)
	if err != nil {
		return err
	}
	if !matched {
		return errFormat
	}
	sum := CalculateIDCardChecksum(s)
	if string(sum) != s[17:] {
		return errFormat
	}
	return nil
}

// ValBirthDate 验证出生日期
func ValBirthDate(s string) error {
	// 判断是否YYYYMMDD格式的日期
	matched, err := regexp.MatchString(`^\d{4}\d{2}\d{2}$`, s)
	if err != nil {
		return err
	}
	if !matched {
		return errors.New("日期格式不正确，正确格式为 YYYYMMDD")
	}
	// 将字符串转为 YYYY-MM-DD 格式以便使用 time.Parse 进行验证
	formattedDate := s[:4] + "-" + s[4:6] + "-" + s[6:]
	const layout = "2006-01-02"
	// 解析输入的日期字符串
	_, err = time.Parse(layout, formattedDate)
	if err != nil {
		return errors.New("无效的日期")
	}
	// 获取当前日期
	now := time.Now()
	birthDate, _ := time.Parse(layout, formattedDate)
	// 判断日期是否在未来
	if birthDate.After(now) {
		return errors.New("出生日期不能在未来")
	}
	return nil
}

// CalculateIDCardChecksum 计算身份证18位检验码
func CalculateIDCardChecksum(id string) rune {
	var weight = [17]int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
	var checksum = [11]rune{'1', '0', 'X', '9', '8', '7', '6', '5', '4', '3', '2'}
	sum := 0
	for i, c := range id[:17] {
		digit := int(c - '0')
		sum += digit * weight[i]
	}
	return checksum[sum%11]
}

// ValAlphanum 验证字母+数字组合
func ValAlphanum(s string) error {
	matched, err := regexp.MatchString("^[A-Za-z0-9]+$", s)
	if err != nil || !matched {
		return errFormat
	}
	return nil
}

// NickNameReg 2-16 个字符，支持数字、大小写字母和汉字
var NickNameReg = "^[A-Za-z0-9\u4e00-\u9fa5]{2,16}$"

// ValNickname 验证昵称格式
func ValNickname(s string) error {
	matched, err := regexp.MatchString(NickNameReg, s)
	if err != nil || !matched {
		return errNickname
	}
	return nil
}

var (
	ErrVersion  = errors.New("版本号格式错误")
	ErrAlphanum = errors.New("Only allow alphanumeric. ")
	ErrDate     = errors.New("日期格式错误:yyyymmdd")
)

// ValVersion 验证版本号
func ValVersion(s string) error {
	matched, err := regexp.MatchString("^(\\d)+(\\.\\d+){0,99}$", s)
	if err != nil || !matched {
		return ErrVersion
	}
	return nil
}

// 美服:2-20/字母/数字
var USNickname = "^[A-Za-z0-9]{2,20}$"
var regUsNickname = regexp.MustCompile(USNickname)

// 日服:2-16/字母/数字/日语/汉字
var JPNickname = "^[A-Za-z0-9\u3040-\u309F\u30A0-\u30FF\u4E00-\u9FA5\uF900-\uFA2D]{2,16}$"
var regJpNickname = regexp.MustCompile(JPNickname)

// 韩服:2-16/字母/数字/韩文
var KrNickname = "^[A-Za-z0-9\uAC00-\uD7A3]{2,16}$"
var regKrNickname = regexp.MustCompile(KrNickname)

// GetNicknameReg 获取区服昵称正则
func GetNicknameReg(country string) string {
	switch country {
	case "US":
		return USNickname
	case "JP":
		return JPNickname
	case "KR":
		return KrNickname
	default:
		return USNickname
	}
}

// ValINTLNickname 验证通行证昵称
func ValINTLNickname(country, s string) bool {
	switch country {
	case "US":
		return ValUsNickname(s)
	case "JP":
		return ValJpNickname(s)
	case "KR":
		return ValKrNickname(s)
	}
	return false
}

// ValUsNickname 美服:2-20/字母/数字
func ValUsNickname(s string) bool {
	return regUsNickname.MatchString(s)
}

// ValJpNickname 日服:2-16/字母/数字/日语/汉字
func ValJpNickname(s string) bool {
	return regJpNickname.MatchString(s)
}

// ValKrNickname 韩服:2-16/字母/数字/韩文
func ValKrNickname(s string) bool {
	return regKrNickname.MatchString(s)
}

func ValEmail(email string) bool {
	pattern := `^(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$` // 匹配电子邮箱
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

// ValDefaultNickname 验证是否为默认昵称
func ValDefaultNickname(s string) bool {
	pattern := `^Y[\d]+$`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(s)
}

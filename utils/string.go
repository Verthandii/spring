package utils

import (
	"bytes"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unicode/utf8"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

// RandStr 获取随机字符串
func RandStr(n int) string {
	src := []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = src[rand.Intn(len(src))]
	}
	return string(b)
}

// TranscodeRandStr 获取引继码字符串
func TranscodeRandStr(n int) string {
	src := []rune("123456789abcdefghjkmnpqrstuvwxyzABCDEFGHJKMNPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = src[rand.Intn(len(src))]
	}
	return string(b)
}

// RandInt 获取随机数 含最小值,不含最大值
func RandInt(min int, max int) int64 {
	return int64(rand.Intn(max-min) + min)
}

// RandSha2 根据时间获取随机sha2字符串
func RandSha2() string {
	s := RandStr(32) + strconv.Itoa(int(time.Now().UnixNano()))
	return Sha2(s)
}

// RandSha1 根据时间获取随机sha1字符串
func RandSha1() string {
	s := RandStr(32) + strconv.Itoa(int(time.Now().UnixNano()))
	return Sha1(s)
}

// RandMd5 根据时间获取随机md5字符串
func RandMd5() string {
	s := RandStr(32) + strconv.Itoa(int(time.Now().UnixNano()))
	return Md5Lower(s)
}

// KeySort 将map的键值按ASCII码排序
func KeySort(params map[string]string) string {
	var dataParams string
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for i, k := range keys {
		dataParams += k + "=" + params[k]
		if i+1 != len(keys) {
			dataParams += "&"
		}
	}
	return dataParams
}

// KeySortUrlEncode 将map的键值按ASCII码排序,value进行urlEncode
func KeySortUrlEncode(params map[string]any) string {
	var dataParams string
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for i, k := range keys {
		dataParams += k + "=" + url.QueryEscape(fmt.Sprintf("%v", params[k]))
		if i+1 != len(keys) {
			dataParams += "&"
		}
	}
	return dataParams
}

// ErrorToString 把recover捕获的异常,转成错误字符串
func ErrorToString(r any) string {
	switch v := r.(type) {
	case error:
		// 管道错误,一般是客户端关闭了tcp连接引起的异常,一般出现在返回数据较大的接口上,忽略即可
		if errors.Is(v, syscall.EPIPE) || errors.Is(v, syscall.ECONNRESET) {
			return "ignore"
		}
		return v.Error()
	case string:
		if v == "" {
			v = "服务器异常"
		}
		return v
	default:
		return "服务器异常"
	}
}

// CheckPwd 检查密码强度
func CheckPwd(minLength, maxLength, minLevel int, pwd string) error {
	if len(pwd) < minLength {
		return fmt.Errorf("密码不能短于%d个字符", minLength)
	}
	if len(pwd) > maxLength {
		return fmt.Errorf("密码长度不能超过%d个字符", maxLength)
	}
	var level = 0
	patternList := []string{`[0-9]+`, `[a-z]+`, `[A-Z]+`, `[~!@#$%^&*?_-]+`}
	for _, pattern := range patternList {
		match, _ := regexp.MatchString(pattern, pwd)
		if match {
			level++
		}
	}
	if level < minLevel {
		return fmt.Errorf("密码复杂度不符合安全要求")
	}
	return nil
}

// VersionOrdinal 获取版本号序数
func VersionOrdinal(version string) string {
	// ISO/IEC 14651:2011
	const maxByte = 1<<8 - 1
	vo := make([]byte, 0, len(version)+8)
	j := -1
	for i := 0; i < len(version); i++ {
		b := version[i]
		if '0' > b || b > '9' {
			vo = append(vo, b)
			j = -1
			continue
		}
		if j == -1 {
			vo = append(vo, 0x00)
			j = len(vo) - 1
		}
		if vo[j] == 1 && vo[j+1] == '0' {
			vo[j+1] = b
			continue
		}
		if vo[j]+1 > maxByte {
			return ""
		}
		vo = append(vo, b)
		vo[j]++
	}
	return string(vo)
}

// EmailSafe 邮箱脱敏
func EmailSafe(s string) string {
	nn := strings.Split(s, "@")
	if len(nn) != 2 {
		return s
	}
	prefixLen := len(nn[0])
	limit := prefixLen
	if prefixLen > 3 {
		limit = 3
	}
	return nn[0][0:limit] + "***@" + nn[1]
}

// EmailSafeINTL 邮箱脱敏-海外SDK
func EmailSafeINTL(s string) string {
	parts := strings.Split(s, "@")
	if len(parts) != 2 {
		return s
	}
	prefix := parts[0]
	domain := parts[1]
	var safePrefix string
	switch {
	case len(prefix) == 1:
		safePrefix = prefix + "***" // 单字符：a → a***
	case len(prefix) >= 5:
		safePrefix = prefix[:2] + "***" + prefix[len(prefix)-2:] // 长度≥5：前2 + *** + 后2
	default:
		half := len(prefix) / 2
		safePrefix = prefix[:half] + "***" + prefix[half:] // 长度2-4：前半 + *** + 后半
	}
	return safePrefix + "@" + domain
}

// MobileSafe 手机脱敏
func MobileSafe(s string) string {
	if len(s) != 11 {
		return s
	}
	return s[:3] + "****" + s[7:]
}

// RealNameSafe 姓名脱敏
func RealNameSafe(s string) string {
	ns := []rune(s)
	length := len(ns)
	if length <= 1 {
		return s
	}
	last := string(ns[length-1:])
	return getStar(length-1, 5) + last
}

// IDCardSafe 身份证脱敏
func IDCardSafe(s string) string {
	l := len(s)
	if l == 0 {
		return ""
	}
	return s[:3] + getStar(9, 0) + s[l-6:l-1] + getStar(1, 0)
}

// IDCardSafeMainlandAndNotMainland 身份证脱敏大陆和非大陆
func IDCardSafeMainlandAndNotMainland(s string, idType int64) string {
	n := len(s)
	switch idType {
	case 0, 1: // 大陆&手动认证
		return IDCardSafe(s)
	case 2: // 港澳居民来往内地通行证
		if len(s) < 6 {
			break
		}
		prefix := s[:3]
		suffix := s[n-3:n-2] + s[n-2:n-1]
		mid := strings.Repeat("*", n-6)
		result := prefix + mid + suffix + "*"
		return result
	case 3: // 台湾居民来往大陆通行证
		if len(s) < 6 {
			break
		}
		prefix := s[:3]
		suffix := s[n-3:n-2] + s[n-2:n-1]
		mid := strings.Repeat("*", n-6)
		result := prefix + mid + suffix + "*"
		return result
	}
	return ""
}

// AppleIDSafe 苹果ID脱敏 前六位+******
func AppleIDSafe(s string) string {
	l := len(s)
	if l <= 6 {
		return s
	}
	return s[:6] + "******"
}

// getStar 获取星星
func getStar(n int, max int) string {
	var s string
	for i := 1; i <= n; i++ {
		s += "*"
		if i == max {
			break
		}
	}
	return s
}

// GetLang 获取yostar服务的语言
func GetLang(lang string) string {
	return "zh-cn"
	if lang == "" {
		return "zh-cn"
	}
	langList := map[string]string{
		"zh-cn": "zh-cn",
		"cn":    "zh-cn",
		"zh":    "zh-cn",
		"ja":    "ja",
		"jp":    "ja",
		"kr":    "kr",
		"ko":    "kr",
		"us":    "en",
		"en":    "en",
	}
	lang = langList[strings.ToLower(lang)]
	if lang == "" {
		lang = "zh-cn"
	}
	return lang
}

// GetIntLang 获取yostar国际服务的语言
func GetIntLang(lang string) string {
	if lang == "" {
		return "en"
	}
	langList := map[string]string{
		"ja":      "ja",
		"jp":      "ja",
		"kr":      "kr",
		"ko":      "kr",
		"us":      "en",
		"en":      "en",
		"fr":      "fr",
		"de":      "de",
		"zh-hant": "zh-Hant",
		"en_mn":   "en_mn",
	}
	lang = langList[strings.ToLower(lang)]
	if lang == "" {
		lang = "en"
	}
	return lang
}

// LeftUpper 首字母转大写
func LeftUpper(s string) string {
	if len(s) > 0 {
		return strings.ToUpper(string(s[0])) + s[1:]
	}
	return s
}

// LeftLower 首字母转小写
func LeftLower(s string) string {
	if len(s) > 0 {
		return strings.ToLower(string(s[0])) + s[1:]
	}
	return s
}

func HeaderToSlice(header []byte) []string {
	hs := bytes.Split(header, []byte("\r\n"))
	var data = make([]string, 0)
	for _, h := range hs {
		if len(h) == 0 {
			continue
		}
		data = append(data, string(h))
	}
	return data
}

// StrDecodeToCert 证书内容解码成证书，不带前缀和后缀
func StrDecodeToCert(s string) (*x509.Certificate, error) {
	x5cBytes, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	cert, err := x509.ParseCertificate(x5cBytes)
	if err != nil {
		return nil, err
	}
	return cert, nil
}

// PemDecodeToCert pem格式解码成证书，带 -----BEGIN CERTIFICATE----- 前后缀
func PemDecodeToCert(s string) (*x509.Certificate, error) {
	block, _ := pem.Decode([]byte(s))
	if block == nil {
		return nil, errors.New("Failed to decode PEM block containing the certificate")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}
	return cert, nil
}

// DefaultLang 获取项目默认语言
func DefaultLang(pid string) string {
	langList := map[string]string{
		"jp": "ja",
		"kr": "kr",
		"us": "en",
		"tw": "zh-Hant",
		"mn": "en_mn",
	}
	if lang, ok := langList[strings.ToLower(pid[0:2])]; ok {
		return lang
	}
	return "en"
}

func IsLetter(r rune) bool {
	if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') {
		return false
	}
	return true
}

// GetFirst 获取第一个字符
func GetFirst(s string) string {
	return s[0:1]
}

// Ternary 三元表达式
func Ternary(exp bool, t, f string) string {
	if exp {
		return t
	}
	return f
}

// 该函数比较两个版本号是否相等，是否大于或小于的关系
// 返回值：0表示v1与v2相等；1表示v1大于v2；2表示v1小于v2
func verCompare(v1, v2 string) int {
	// 替换一些常见的版本符号
	replaceMap := map[string]string{"V": "", "v": "", "-": "."}
	for k, v := range replaceMap {
		if strings.Contains(v1, k) {
			v1 = strings.ReplaceAll(v1, k, v)
		}
		if strings.Contains(v2, k) {
			v2 = strings.ReplaceAll(v2, k, v)
		}
	}
	verStr1 := strings.Split(v1, ".")
	verStr2 := strings.Split(v2, ".")
	ver1 := strSlice2IntSlice(verStr1)
	ver2 := strSlice2IntSlice(verStr2)
	// 找出v1和v2哪一个最短
	shorter := len(ver1)
	if len(ver1) > len(ver2) {
		shorter = len(ver2)
	}
	// 循环比较
	for i := 0; i < shorter; i++ {
		if ver1[i] == ver2[i] {
			if shorter-1 == i {
				if len(ver1) == len(ver2) {
					return 0
				} else {
					if len(ver1) > len(ver2) {
						return 1
					} else {
						return 2
					}
				}
			}
		} else if ver1[i] > ver2[i] {
			return 1
		} else {
			return 2
		}
	}
	return -1
}

// strSlice2IntSlice 版本转int切片
func strSlice2IntSlice(version []string) []int64 {
	if len(version) == 0 {
		return []int64{}
	}
	retInt := make([]int64, 0, len(version))
	for _, str := range version {
		i, err := strconv.ParseInt(str, 10, 64)
		if err == nil {
			retInt = append(retInt, i)
		}
	}
	return retInt
}

// VersionCompare 版本号对比
// operator 操作符 ==,>,>=,<,<=
func VersionCompare(v1, operator, v2 string) bool {
	com := verCompare(v1, v2)
	switch operator {
	case "==":
		return com == 0
	case ">":
		return com == 1
	case ">=":
		return com == 0 || com == 1
	case "<":
		return com == 2
	case "<=":
		return com == 0 || com == 2
	}
	return false
}

func IsEmpty(str string) bool {
	return TrimSpace(str) == ""
}

func IsSpace(r rune) bool {
	return r == 0x3000 || r == 0x0020 || r == 0x00A0
}

func TrimSpace(str string) string {
	return strings.TrimFunc(str, IsSpace)
}

// TruncateBytes 截取前 n 个字节，并确保不截断 UTF-8 字符
func TruncateBytes(s string, n int) string {
	if len(s) <= n {
		return s
	}
	// 截断到 n 字节的位置，并确保不会截断 UTF-8 字符
	for i := n; i >= 0; i-- {
		if utf8.ValidString(s[:i]) {
			return s[:i]
		}
	}
	return ""
}

// EncodeToShiftJIS 将字符串编码为 Shift-JIS
func EncodeToShiftJIS(input string) (string, error) {
	encoded, err := io.ReadAll(transform.NewReader(strings.NewReader(input), japanese.ShiftJIS.NewEncoder()))
	if err != nil {
		return "", err
	}
	return string(encoded), nil
}

// EncodeToShiftJISLength 将字符串编码为 Shift-JIS，截取前 maxLen 字节
func EncodeToShiftJISLength(input string, maxLen int) (string, error) {
	for {
		jis, err := EncodeToShiftJIS(input)
		if err != nil {
			return "", err
		}
		if maxLen <= 0 || len(jis) <= maxLen {
			return jis, nil
		}
		diff := maxLen - len(jis)
		if diff < 1 {
			diff = 1
		}
		input = TruncateBytes(input, len(input)-diff)
	}
}

// ShiftJISToUTF8 将 Shift-JIS 编码的字符串转换为 UTF-8
func ShiftJISToUTF8(input []byte) (string, error) {
	reader := transform.NewReader(bytes.NewReader(input), japanese.ShiftJIS.NewDecoder())
	// 读取并转换为 UTF-8 字符串
	utf8Bytes, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(utf8Bytes), nil
}

// BoolToString bool->string
func BoolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

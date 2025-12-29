package utils

import (
	"database/sql/driver"
	"fmt"
	"time"
)

var (
	DefaultZone = time.FixedZone("UTC+8", 8*3600)
	Timezone, _ = time.LoadLocation(timezone)
)

const (
	// ModelTime MarshalJSON
	timeFormat = "2006/01/02 15:04:05"
	timezone   = "Asia/Shanghai"
)

// Now 获取当前时间
func Now() int64 {
	return time.Now().Unix()
}

// NowMill 当前毫秒数
func NowMill() int64 {
	return time.Now().UnixNano() / 1000 / 1000
}

// Midnight 获取 t 时的零点零分零秒
func Midnight(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
}

// Date 当前当前日期
func Date() string {
	return Time().Format("20060102")
}

// DateTimeYmdHis 获取当前时间
func DateTimeYmdHis() string {
	return Time().Format("20060102150405")
}

// DateTimeYmdHis2 获取当前时间
func DateTimeYmdHis2() string {
	return Time().Format("2006-01-02 15:04:05")
}

// Time 获取当前时区时间
func Time() time.Time {
	return time.Now().In(time.Local)
}

// DateTime 时间戳转换成当前项目对应时区时间
func DateTime(t int64) string {
	if t <= 0 {
		return ""
	}
	return time.Unix(t, 0).In(time.Local).Format("2006-01-02 15:04:05")
}

// DateTimeToTime 时期时间转为Time
func DateTimeToTime(s string) (time.Time, error) {
	return time.ParseInLocation("2006-01-02 15:04:05", s, time.Local)
}

// IsPm89 是否为晚8点到晚9点
func IsPm89(time time.Time) bool {
	nowTime := time.Format("1504")
	if nowTime >= "2000" && nowTime <= "2059" {
		return true
	}
	return false
}

// CuMonth 获取本月第一天和最后一天的时间戳
func CuMonth() (int64, int64) {
	now := time.Now().In(time.Local)
	currentYear, currentMonth, _ := now.Date()
	begin := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, time.Local)
	end := begin.AddDate(0, 1, 0)
	return begin.Unix(), end.Unix() - 1
}

// MonthStar 当月第一天时间
func MonthStar() time.Time {
	now := time.Now().In(time.Local)
	currentYear, currentMonth, _ := now.Date()
	return time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, time.Local)
}

// MonthEnd 当月最后一天时间
func MonthEnd() time.Time {
	now := time.Now().In(time.Local)
	currentYear, currentMonth, _ := now.Date()
	begin := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, time.Local)
	return begin.AddDate(0, 1, 0).Add(-time.Second)
}

// CuToday  获取今天开始/结束时间戳
func CuToday() (int64, int64) {
	timeStr := time.Now().In(time.Local).Format("2006-01-02")
	begin, _ := time.ParseInLocation("2006-01-02 15:04:05", timeStr+" 00:00:00", time.Local)
	end, _ := time.ParseInLocation("2006-01-02", timeStr, time.Local)
	return begin.Unix(), end.Add(time.Hour*24).Unix() - 1
}

// GetAge 根据身份证计算年龄
func GetAge(idCard string) int {
	if len(idCard) != 18 {
		return 0
	}
	birthDate := idCard[6:14]
	if birthDate == "" {
		return 0
	}
	birthday, err := time.ParseInLocation("20060102", birthDate, time.Local)
	if err != nil {
		return 0
	}
	now := time.Now()
	age := now.Year() - birthday.Year()
	if now.Format("0102") < birthday.Format("0102") {
		age--
	}
	return age
}

// GetAgeByBirthDate 根据出生日期计算年龄
func GetAgeByBirthDate(idCard string, birthDate string) int {
	if len(birthDate) == 8 {
		return GetAge(fmt.Sprintf("******%s****", birthDate))
	}
	return GetAge(idCard)
}

// ModelTime 时间格式
type ModelTime struct {
	time.Time
}

// MarshalJSON 转 JSON
func (t *ModelTime) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(timeFormat)+2)
	b = append(b, '"')
	b = t.AppendFormat(b, timeFormat)
	b = append(b, '"')
	return b, nil
}

// UnmarshalJSON JSON 转 ModelTime
func (t *ModelTime) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(`"`+timeFormat+`"`, string(data), time.Local)
	*t = ModelTime{Time: now}
	return
}

func (t *ModelTime) String() string {
	return t.Format(timeFormat)
}

func (t *ModelTime) Local() ModelTime {
	loc, _ := time.LoadLocation(timezone)

	return ModelTime{t.In(loc)}
}

// Value insert timestamp into mysql need this function.
func (t *ModelTime) Value() (driver.Value, error) {
	return t.Time, nil
}

// Scan valueof time.Time
func (t *ModelTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		t.Time = value
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

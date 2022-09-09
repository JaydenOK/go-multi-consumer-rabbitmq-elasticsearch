package utils

import (
	"app/constant"
	"database/sql/driver"
	"fmt"
	"time"
)

type LocalTime time.Time

//实现MarshalJSON接口，格式化数据
func (t LocalTime) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format(constant.TimeFormat))
	return []byte(stamp), nil
}

// 写入 mysql 时调用
func (t LocalTime) Value() (driver.Value, error) {
	// 0001-01-01 00:00:00 属于空值，遇到空值解析成 null 即可
	if t.String() == "0001-01-01 00:00:00" {
		return nil, nil
	}
	return []byte(time.Time(t).Format(constant.TimeFormat)), nil
}

// 检出 mysql 时调用
func (t *LocalTime) Scan(v interface{}) error {
	// mysql 内部日期的格式可能是 2006-01-02 15:04:05 +0800 CST 格式，所以检出的时候还需要进行一次格式化
	tTime, _ := time.Parse("2006-01-02 15:04:05 +0800 CST", v.(time.Time).String())
	*t = LocalTime(tTime)
	return nil
}

// 用于 fmt.Println 和后续验证场景
func (t LocalTime) String() string {
	return time.Time(t).Format(constant.TimeFormat)
}

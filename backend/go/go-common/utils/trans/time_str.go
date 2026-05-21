package trans

import (
	"go-common/utils/catch"
	"time"
)

const (
	YyyyMmDdHhMmSs  = "2006-01-02 15:04:05"
	YyyyMmDdHhMmSsT = "2006-01-02T15:04:05+08:00"
	YyyyMmDd        = "2006-01-02"
)

func Str2Milli(value string, format string) (int64, error) {
	location, err := time.ParseInLocation(format, value, time.Local)
	if err != nil {
		return 0, err
	}
	return location.UnixMilli(), nil
}

func TryStr2Milli(value string, format string) int64 {
	return catch.Try1(Str2Milli(format, value))
}

func Milli2Str(value int64, format string) string {
	return time.UnixMilli(value).Format(format)
}

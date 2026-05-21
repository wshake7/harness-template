package str

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

func GetSuffix(s string) string {
	ext := filepath.Ext(s)
	if ext == "" {
		return ""
	}
	// 去掉开头的点号
	return ext[1:]
}

func GetSuffixWithDot(s string) string {
	ext := filepath.Ext(s)
	if ext == "" {
		return ""
	}
	// 去掉开头的点号
	return ext
}

var printfReg = regexp.MustCompile(`\{(\w+)}`)

func Join(ss ...string) string {
	return strings.Join(ss, "")
}

func TrimStringPtr(value *string, transform func(string) string) {
	if value == nil {
		return
	}
	next := strings.TrimSpace(*value)
	if transform != nil {
		next = transform(next)
	}
	*value = next
}

func Sprintf(format string, params map[string]any) string {
	return printfReg.ReplaceAllStringFunc(format, func(str string) string {
		key := printfReg.FindStringSubmatch(str)[1]
		if val, ok := params[key]; ok {
			return fmt.Sprint(val)
		}
		return str
	})
}

type ValueF struct {
	Value string
}

func (v *ValueF) Sprintf(arg ...any) string {
	return fmt.Sprintf(v.Value, arg...)
}

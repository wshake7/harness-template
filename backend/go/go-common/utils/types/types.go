package types

import (
	"reflect"
)

func IsPointer(v any) bool {
	// 处理 nil 的情况：如果 v 本身是 nil interface{}，ValueOf 返回的 Value 是无效的
	if v == nil {
		return false
	}

	// 获取反射值
	rv := reflect.ValueOf(v)

	// 检查 Kind 是否为 Ptr
	return rv.Kind() == reflect.Pointer
}

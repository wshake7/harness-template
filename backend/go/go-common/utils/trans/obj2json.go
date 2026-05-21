package trans

import (
	"github.com/bytedance/sonic"
	"go-common/utils/catch"
)

func Obj2Json(obj any) (string, error) {
	bytes, err := sonic.Marshal(obj)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func TryObj2Json(obj any) string {
	return catch.Try1(Obj2Json(obj))
}

func Json2Obj[T any](m string) (T, error) {
	t := new(T)
	err := sonic.Unmarshal([]byte(m), t)
	if err != nil {
		return *new(T), err
	}
	return *t, nil
}

func TryJson2Obj[T any](m string) T {
	return catch.Try1(Json2Obj[T](m))
}

package trans

import (
	"github.com/bytedance/sonic"
	"go-common/utils/catch"
)

func Obj2Map[T any](obj any) (map[string]T, error) {
	bytes, err := sonic.Marshal(obj)
	if err != nil {
		return nil, err
	}
	temp := new(map[string]T)
	err = sonic.Unmarshal(bytes, temp)
	if err != nil {
		return nil, err
	}
	return *temp, nil
}

func TryObj2Map[T any](obj any) map[string]T {
	return catch.Try1(Obj2Map[T](obj))
}

func Map2Obj[T any](m map[string]any) (T, error) {
	t := new(T)
	data, err := sonic.Marshal(m)
	if err != nil {
		return *t, err
	}
	err = sonic.Unmarshal(data, t)
	if err != nil {
		return *t, err
	}
	return *t, nil
}

func TryMap2Obj[T any](m map[string]any) T {
	return catch.Try1(Map2Obj[T](m))
}

package trans

import (
	"go-common/utils/catch"
	"strconv"
)

func TryStr2U8(str string) (uint8, error) {
	v, err := strconv.ParseUint(str, 10, 8)
	if err != nil {
		return 0, err
	}
	return uint8(v), nil
}

func Str2U8(str string) uint8 {
	return catch.Try1(TryStr2U8(str))
}

func TryStr2I8(str string) (int8, error) {
	v, err := strconv.ParseInt(str, 10, 8)
	if err != nil {
		return 0, err
	}
	return int8(v), nil
}

func Str2I8(str string) int8 {
	return catch.Try1(TryStr2I8(str))
}

func TryStr2U16(str string) (uint16, error) {
	v, err := strconv.ParseUint(str, 10, 16)
	if err != nil {
		return 0, err
	}
	return uint16(v), nil
}

func Str2U16(str string) uint16 {
	return catch.Try1(TryStr2U16(str))
}

func TryStr2I16(str string) (int16, error) {
	v, err := strconv.ParseInt(str, 10, 16)
	if err != nil {
		return 0, err
	}
	return int16(v), nil
}

func Str2I16(str string) int16 {
	return catch.Try1(TryStr2I16(str))
}

func TryStr2U32(str string) (uint32, error) {
	v, err := strconv.ParseUint(str, 10, 32)
	if err != nil {

	}
	return uint32(v), nil
}

func Str2U32(str string) uint32 {
	return catch.Try1(TryStr2U32(str))
}

func TryStr2I32(str string) (int32, error) {
	v, err := strconv.ParseInt(str, 10, 32)
	if err != nil {

	}
	return int32(v), nil
}

func Str2I32(str string) int32 {
	return catch.Try1(TryStr2I32(str))
}

func TryStr2Int(str string) (int, error) {
	v, err := strconv.ParseInt(str, 10, 32)
	if err != nil {

	}
	return int(v), nil
}

func Str2Int(str string) int {
	return catch.Try1(TryStr2Int(str))
}

func TryStr2U64(str string) (uint64, error) {
	v, err := strconv.ParseUint(str, 10, 64)
	if err != nil {

	}
	return v, nil
}

func Str2U64(str string) uint64 {
	return catch.Try1(TryStr2U64(str))
}

func TryStr2I64(str string) (int64, error) {
	v, err := strconv.ParseInt(str, 10, 64)
	if err != nil {

	}
	return v, nil
}

func Str2I64(str string) int64 {
	return catch.Try1(TryStr2I64(str))
}

func TryStr2F32(str string) (float32, error) {
	v, err := strconv.ParseFloat(str, 32)
	if err != nil {

	}
	return float32(v), nil
}

func Str2F32(str string) float32 {
	return catch.Try1(TryStr2F32(str))
}

func TryStr2F64(str string) (float64, error) {
	v, err := strconv.ParseFloat(str, 64)
	if err != nil {

	}
	return v, nil
}

func Str2F64(str string) float64 {
	return catch.Try1(TryStr2F64(str))
}

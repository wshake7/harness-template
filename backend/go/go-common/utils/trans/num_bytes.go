package trans

import (
	"fmt"
	"go-common/types"
	"slices"
	"unsafe"
)

func Number2Bytes[T types.Number](t T) []byte {
	size := int(unsafe.Sizeof(*new(T)))
	u := UnsafeTrans[uint64](t)
	bytes := make([]byte, size)
	for i := range size {
		bytes[i] |= byte(u >> (i * 8))
	}
	return bytes
}

func Bytes2Number[T types.Number](bytes []byte) T {
	size := int(unsafe.Sizeof(*new(T)))
	if len(bytes) < size {
		return T(0)
	}
	t := uint64(0)
	for i, b := range bytes[0:size] {
		t |= uint64(b) << (i * 8)
	}
	return UnsafeTrans[T](t)
}

func Numbers2Bytes[T types.Number](numbers []T) []byte {
	size := int(unsafe.Sizeof(*new(T)))
	numbers = slices.Clone(numbers)
	bytes := UnsafeTrans[[]byte](numbers)
	ptr := (*[3]int)(unsafe.Pointer(&bytes))
	ptr[1], ptr[2] = ptr[1]*size, ptr[2]*size
	return bytes
}

func Bytes2Numbers[T types.Number](bytes []byte) []T {
	size := int(unsafe.Sizeof(*new(T)))
	if len(bytes)%size > 0 {
		panic(fmt.Errorf("bad bytes len: %v", len(bytes)))
	}
	bytes = slices.Clone(bytes)
	numbers := UnsafeTrans[[]T](bytes)
	ptr := (*[3]int)(unsafe.Pointer(&numbers))
	ptr[1], ptr[2] = ptr[1]/size, ptr[2]/size
	return numbers
}

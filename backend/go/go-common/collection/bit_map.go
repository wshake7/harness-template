package collection

import (
	"go-common/types"
	"math/bits"
	"unsafe"
)

type BitMap[T types.Integer] struct {
	value T
}

func BitMapNew[T types.Integer](value T) BitMap[T] {
	return BitMap[T]{value}
}

func (bm *BitMap[T]) Value() T {
	return bm.value
}

// Set 把T的第index位置为b
func (bm *BitMap[T]) Set(index int, b bool) {
	// 添加边界检查，避免无效操作
	if index < 0 || index >= int(unsafe.Sizeof(bm.value))*8 {
		return
	}

	mask := T(1) << index
	if b {
		bm.value |= mask
	} else {
		bm.value &^= mask // 使用 &^ 操作符更简洁
	}
}

// Get 获取T中的第index位置的值
func (bm *BitMap[T]) Get(index int) bool {
	if index < 0 || index >= int(unsafe.Sizeof(bm.value))*8 {
		return false
	}
	return bm.get(index)
}

// Count 使用 bits 包优化计数
func (bm *BitMap[T]) Count() int {
	// 根据类型大小选择对应的 OnesCount 函数
	switch any(bm.value).(type) {
	case uint8:
		return bits.OnesCount8(uint8(bm.value))
	case uint16:
		return bits.OnesCount16(uint16(bm.value))
	case uint32:
		return bits.OnesCount32(uint32(bm.value))
	case uint64:
		return bits.OnesCount64(uint64(bm.value))
	default:
		// 使用 uint64 作为通用处理
		return bits.OnesCount64(uint64(bm.value))
	}
}

func (bm *BitMap[T]) get(index int) bool {
	return (bm.value>>index)&1 == 1
}

type BytesBitMap struct {
	value []byte
}

func BytesBitMapNew(value []byte) BytesBitMap {
	return BytesBitMap{value}
}

// NewBytesBitMapWithCapacity 预分配容量
func NewBytesBitMapWithCapacity(bitCount int) BytesBitMap {
	byteCount := (bitCount + 7) / 8
	return BytesBitMap{make([]byte, byteCount)}
}

func (bm *BytesBitMap) Value() []byte {
	return bm.value
}

func (bm *BytesBitMap) String() string {
	return string(bm.value)
}

// Set 把T的第index位置为b
func (bm *BytesBitMap) Set(index int, b bool) {
	if index < 0 {
		return
	}

	byteIndex := index >> 3 // 等价于 index / 8，位运算更快
	bitIndex := index & 7   // 等价于 index % 8

	if byteIndex >= len(bm.value) {
		// 使用更合理的增长策略
		newLen := max(byteIndex+1, len(bm.value)*2)
		value := make([]byte, newLen)
		copy(value, bm.value)
		bm.value = value
	}

	mask := byte(1) << bitIndex
	if b {
		bm.value[byteIndex] |= mask
	} else {
		bm.value[byteIndex] &^= mask
	}
}

// Get 获取T中的第index位置的值
func (bm *BytesBitMap) Get(index int) bool {
	if index < 0 || index >= len(bm.value)*8 {
		return false
	}
	return bm.get(index)
}

// Count 使用 bits.OnesCount8 优化
func (bm *BytesBitMap) Count() int {
	count := 0
	for _, b := range bm.value {
		count += bits.OnesCount8(b)
	}
	return count
}

func (bm *BytesBitMap) Len() int {
	return len(bm.value) * 8
}

func (bm *BytesBitMap) get(index int) bool {
	byteIndex := index >> 3
	bitIndex := index & 7
	return (bm.value[byteIndex]>>bitIndex)&1 == 1
}

func (bm *BytesBitMap) ForEach(fn func(index int, b bool)) {
	length := len(bm.value) * 8
	for i := range length {
		fn(i, bm.get(i))
	}
}

// ForEachSet 只遍历设置为 true 的位，性能更好
func (bm *BytesBitMap) ForEachSet(fn func(index int)) {
	for byteIndex, b := range bm.value {
		if b == 0 {
			continue
		}
		baseIndex := byteIndex << 3
		for bitIndex := range 8 {
			if (b>>bitIndex)&1 == 1 {
				fn(baseIndex + bitIndex)
			}
		}
	}
}

// ToMap 优化：预先计算容量
func ToMap[T comparable](bm *BytesBitMap, fn func(v int) T) map[T]types.Unit {
	count := bm.Count()
	res := make(map[T]types.Unit, count)

	bm.ForEachSet(func(index int) {
		res[fn(index)] = types.Unit{}
	})

	return res
}

// ToSlice 优化：预先计算容量
func ToSlice[T any](bm *BytesBitMap, fn func(v int) T) []T {
	count := bm.Count()
	res := make([]T, 0, count)

	bm.ForEachSet(func(index int) {
		res = append(res, fn(index))
	})

	return res
}

// Clear 清空位图
func (bm *BytesBitMap) Clear() {
	for i := range bm.value {
		bm.value[i] = 0
	}
}

// Clone 克隆位图
func (bm *BytesBitMap) Clone() BytesBitMap {
	value := make([]byte, len(bm.value))
	copy(value, bm.value)
	return BytesBitMap{value}
}

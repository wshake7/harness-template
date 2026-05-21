package slices_utils

import (
	"slices"
)

// Map 遍历 slice，将每个元素映射为一个新值，生成新的 slice。
// predicate 接收元素下标和元素本身。
func Map[T any, TS ~[]T, R any, RS ~[]R](
	slice TS,
	predicate func(index int, item T) R,
) RS {
	result := make(RS, len(slice))
	for i, v := range slice {
		result[i] = predicate(i, v)
	}
	return result
}

// MapFn 是 Map 的简化版本，不关心元素下标。
func MapFn[T any, TS ~[]T, R any, RS ~[]R](
	slice TS,
	predicate func(T) R,
) RS {
	return Map[T, TS, R, RS](slice, func(_ int, item T) R {
		return predicate(item)
	})
}

// FlatMap 对 slice 中的每个元素进行映射，并将结果打平成一个新的 slice。
func FlatMap[T any, R any, TS ~[]T, RS ~[]R](
	slice TS,
	fn func(T) RS,
) RS {
	result := make(RS, 0)
	for _, v := range slice {
		result = append(result, fn(v)...)
	}
	return result
}

// Reduce 从左到右遍历 slice，将元素聚合为一个值。
func Reduce[T any, TS ~[]T, R any](
	slice TS,
	initial R,
	predicate func(agg R, index int, item T) R,
) R {
	for i, v := range slice {
		initial = predicate(initial, i, v)
	}
	return initial
}

// ReduceRight 从右到左遍历 slice，将元素聚合为一个值。
func ReduceRight[T any, TS ~[]T, R any](
	slice TS,
	initial R,
	predicate func(agg R, index int, item T) R,
) R {
	for i := len(slice) - 1; i >= 0; i-- {
		initial = predicate(initial, i, slice[i])
	}
	return initial
}

// GroupBy 按指定 key 函数对 slice 中的元素进行分组。
func GroupBy[K comparable, T any, TS ~[]T](
	slice TS,
	keyFn func(T) K,
) map[K]TS {
	result := make(map[K]TS)
	for _, v := range slice {
		k := keyFn(v)
		result[k] = append(result[k], v)
	}
	return result
}

// GroupCount 按 key 统计 slice 中元素的数量。
func GroupCount[K comparable, T any, TS ~[]T](
	slice TS,
	keyFn func(T) K,
) map[K]int {
	result := make(map[K]int)
	for _, v := range slice {
		result[keyFn(v)]++
	}
	return result
}

// GroupByWithMapper 在分组的同时，将元素映射为新的值。
func GroupByWithMapper[K comparable, V any, T any, TS ~[]T, VS ~[]V](
	items TS,
	keyMapper func(T) K,
	valueMapper func(T) V,
) map[K]VS {
	result := make(map[K]VS)
	for _, item := range items {
		k := keyMapper(item)
		result[k] = append(result[k], valueMapper(item))
	}
	return result
}

// ToMapKV 将 slice 转换为 map，key/value 由 predicate 提供。
// 若 key 重复，后者会覆盖前者。
func ToMapKV[T any, TS ~[]T, K comparable, V any](
	slice TS,
	predicate func(item T) (K, V),
) map[K]V {
	result := make(map[K]V, len(slice))
	for _, v := range slice {
		k, val := predicate(v)
		result[k] = val
	}
	return result
}

// ToMap 将 slice 转换为 map，value 仍为元素本身。
func ToMap[T any, TS ~[]T, K comparable](
	slice TS,
	predicate func(item T) (K, T),
) map[K]T {
	result := make(map[K]T, len(slice))
	for _, v := range slice {
		k, val := predicate(v)
		result[k] = val
	}
	return result
}

// Index 返回元素在 slice 中第一次出现的位置。
func Index[T comparable](slice []T, element T) (int, bool) {
	i := slices.Index(slice, element)
	return i, i >= 0
}

// FindIndex 返回第一个满足条件的元素下标。
func FindIndex[T any, TS ~[]T](
	slice TS,
	predicate func(T) bool,
) (int, bool) {
	i := slices.IndexFunc(slice, predicate)
	return i, i >= 0
}

// FindLastIndex 从右向左查找第一个满足条件的元素下标。
func FindLastIndex[T any, TS ~[]T](
	slice TS,
	predicate func(T) bool,
) (int, bool) {
	for i := len(slice) - 1; i >= 0; i-- {
		if predicate(slice[i]) {
			return i, true
		}
	}
	return -1, false
}

// Find 查找 slice 中第一个满足条件的元素。
func Find[T any, TS ~[]T](
	slice TS,
	predicate func(T) bool,
) (T, bool) {
	for _, item := range slice {
		if predicate(item) {
			return item, true
		}
	}
	var zero T
	return zero, false
}

// FindPtr 返回第一个满足条件的元素指针。
func FindPtr[T any, TS ~[]T](
	slice TS,
	predicate func(T) bool,
) *T {
	for i := range slice {
		if predicate(slice[i]) {
			return &slice[i]
		}
	}
	return nil
}

// First 返回 slice 的第一个元素。
func First[T any, TS ~[]T](slice TS) (T, bool) {
	if len(slice) == 0 {
		var zero T
		return zero, false
	}
	return slice[0], true
}

// FirstPtr 返回 slice 第一个元素的指针。
func FirstPtr[T any, TS ~[]T](slice TS) *T {
	if len(slice) == 0 {
		return nil
	}
	return &slice[0]
}

// Last 返回 slice 的最后一个元素。
func Last[T any, TS ~[]T](slice TS) (T, bool) {
	if len(slice) == 0 {
		var zero T
		return zero, false
	}
	return slice[len(slice)-1], true
}

// LastPtr 返回 slice 最后一个元素的指针。
func LastPtr[T any, TS ~[]T](slice TS) *T {
	if len(slice) == 0 {
		return nil
	}
	return &slice[len(slice)-1]
}

// Count 统计满足条件的元素数量。
func Count[T any, TS ~[]T](
	slice TS,
	predicate func(T) bool,
) int {
	n := 0
	for _, v := range slice {
		if predicate(v) {
			n++
		}
	}
	return n
}

// Every 判断 slice 中的所有元素是否都满足条件。
func Every[T any, TS ~[]T](
	slice TS,
	predicate func(index int, item T) bool,
) bool {
	for i, v := range slice {
		if !predicate(i, v) {
			return false
		}
	}
	return true
}

// Some 判断 slice 中是否存在至少一个满足条件的元素。
func Some[T any, TS ~[]T](
	slice TS,
	predicate func(index int, item T) bool,
) bool {
	for i, v := range slice {
		if predicate(i, v) {
			return true
		}
	}
	return false
}

// Filter 过滤 slice 中满足条件的元素。
func Filter[T any, TS ~[]T](
	slice TS,
	predicate func(index int, item T) bool,
) TS {
	result := make(TS, 0, len(slice))
	for i, v := range slice {
		if predicate(i, v) {
			result = append(result, v)
		}
	}
	return result
}

// Difference 返回 slice 中存在但 compared 中不存在的元素。
func Difference[T comparable, TS ~[]T](slice, compared TS) TS {
	set := make(map[T]struct{}, len(compared))
	for _, v := range compared {
		set[v] = struct{}{}
	}

	result := make(TS, 0, len(slice))
	for _, v := range slice {
		if _, ok := set[v]; !ok {
			result = append(result, v)
		}
	}
	return result
}

// Intersect 返回多个 slice 的交集（保持第一个 slice 的顺序）。
func Intersect[T comparable, TS ~[]T](slices ...TS) TS {
	if len(slices) == 0 {
		return TS{}
	}

	m := make(map[T]struct{})
	for _, v := range slices[0] {
		m[v] = struct{}{}
	}

	for i := 1; i < len(slices); i++ {
		next := make(map[T]struct{})
		for _, v := range slices[i] {
			if _, ok := m[v]; ok {
				next[v] = struct{}{}
			}
		}
		m = next
		if len(m) == 0 {
			break
		}
	}

	result := make(TS, 0, len(m))
	for _, v := range slices[0] {
		if _, ok := m[v]; ok {
			result = append(result, v)
		}
	}
	return result
}

// Union 返回多个 slice 的并集，保持首次出现的顺序。
func Union[T comparable, TS ~[]T](slices ...TS) TS {
	seen := make(map[T]struct{})
	result := make(TS, 0)

	for _, sc := range slices {
		for _, v := range sc {
			if _, ok := seen[v]; ok {
				continue
			}
			seen[v] = struct{}{}
			result = append(result, v)
		}
	}
	return result
}

// Distinct 去重，保留首次出现的元素。
func Distinct[T comparable, TS ~[]T](slice TS) TS {
	seen := make(map[T]struct{}, len(slice))
	result := make(TS, 0, len(slice))

	for _, v := range slice {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		result = append(result, v)
	}
	return result
}

// DistinctFn 根据 keyFn 结果去重。
func DistinctFn[T any, K comparable, TS ~[]T](
	slice TS,
	keyFn func(T) K,
) TS {
	seen := make(map[K]struct{}, len(slice))
	result := make(TS, 0, len(slice))

	for _, item := range slice {
		key := keyFn(item)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		result = append(result, item)
	}
	return result
}

// ChunkList 将 slice 按指定大小拆分为多个子 slice。
func ChunkList[T any, TS ~[]T, RS ~[]TS](slice TS, size int) RS {
	if size <= 0 || len(slice) == 0 {
		return RS{}
	}

	if size >= len(slice) {
		return RS{slice}
	}

	n := (len(slice) + size - 1) / size
	result := make(RS, 0, n)

	for i := 0; i < len(slice); i += size {
		end := min(i+size, len(slice))
		result = append(result, slice[i:end])
	}
	return result
}

// ChunkForeach 按 chunk 遍历 slice，callback 返回 false 可提前终止。
func ChunkForeach[T any, TS ~[]T](
	slice TS,
	size int,
	callback func(ts TS) bool,
) {
	if size <= 0 || len(slice) == 0 {
		return
	}
	slices.Chunk(slice, size)(callback)
}

// Partition 根据条件将 slice 拆分为两个 slice。
func Partition[T any, TS ~[]T](
	slice TS,
	predicate func(T) bool,
) (TS, TS) {
	left := make(TS, 0)
	right := make(TS, 0)

	for _, v := range slice {
		if predicate(v) {
			left = append(left, v)
		} else {
			right = append(right, v)
		}
	}
	return left, right
}

// ReverseCopy 返回反转后的新 slice。
func ReverseCopy[T any, TS ~[]T](slice TS) TS {
	n := len(slice)
	result := make(TS, n)
	for i, v := range slice {
		result[n-1-i] = v
	}
	return result
}

// Reverse 原地反转 slice。
func Reverse[T any, TS ~[]T](slice TS) {
	slices.Reverse(slice)
}

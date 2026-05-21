package maps_utils

import (
	"net/url"
	"slices"
)

// Keys 返回 map 中的所有 key，顺序不保证
func Keys[M ~map[K]V, K comparable, V any](m M) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// Values 返回 map 中的所有 value，顺序不保证
func Values[M ~map[K]V, K comparable, V any](m M) []V {
	values := make([]V, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

// Entry 表示 map 的一个键值对
type Entry[K comparable, V any] struct {
	Key   K
	Value V
}

// Entries 将 map 转换为 Entry 切片，顺序不保证
func Entries[M ~map[K]V, K comparable, V any](m M) []Entry[K, V] {
	entries := make([]Entry[K, V], 0, len(m))
	for k, v := range m {
		entries = append(entries, Entry[K, V]{Key: k, Value: v})
	}
	return entries
}

// MapToSlice 将 map 映射为任意 slice
// 常用于 map → DTO / map → struct / map → string 等场景
func MapToSlice[M ~map[K]V, K comparable, V any, R any, RS ~[]R](m M, fn func(K, V) R) RS {
	result := make([]R, 0, len(m))
	for k, v := range m {
		result = append(result, fn(k, v))
	}
	return result
}

// Url 将 map 转为 url.Values，每个 key 只保留一个 value（使用 Set）
func Url[M ~map[K]V, K comparable, V any](m M, keyMapper func(K) string, valueMapper func(V) string) url.Values {
	result := url.Values{}
	for k, v := range m {
		result.Set(keyMapper(k), valueMapper(v))
	}
	return result
}

// SortUrl 在生成 url.Values 前按 key 排序（常用于签名、缓存 key）
func SortUrl[M ~map[K]V, K comparable, V any](m M, sortFunc func(K, K) int, keyMapper func(K) string, valueMapper func(V) string) url.Values {
	keys := Keys(m)
	slices.SortFunc(keys, sortFunc)
	result := url.Values{}
	for _, key := range keys {
		result.Set(keyMapper(key), valueMapper(m[key]))
	}
	return result
}

// UrlMulti 将 map[K][]V 转为 url.Values
// 每个 slice 元素使用 Add 追加
func UrlMulti[M ~map[K]VS, K comparable, V any, VS ~[]V](m M, keyMapper func(K) string, valueMapper func(V) string) url.Values {
	result := url.Values{}
	for k, v := range m {
		for _, t := range v {
			result.Add(keyMapper(k), valueMapper(t))
		}
	}
	return result
}

// SortUrlMulti 与 UrlMulti 类似，但会先对 key 排序
func SortUrlMulti[M ~map[K]VS, K comparable, V any, VS ~[]V](m M, sortFunc func(K, K) int, keyMapper func(K) string, valueMapper func(V) string) url.Values {
	keys := Keys(m)
	slices.SortFunc(keys, sortFunc)
	result := url.Values{}
	for _, key := range keys {
		v := m[key]
		for _, t := range v {
			result.Add(keyMapper(key), valueMapper(t))
		}
	}
	return result
}

// MergeWith 将 src 合并到 dst
// 若 key 已存在，使用 fn(old, new) 决定最终值
func MergeWith[M ~map[K]V, K comparable, V any](dst M, src M, fn func(old V, new V) V) M {
	if dst == nil {
		dst = make(M, len(src))
	}
	for k, v := range src {
		if old, ok := dst[k]; ok {
			dst[k] = fn(old, v)
		} else {
			dst[k] = v
		}
	}
	return dst
}

// Group 按 keyFn 对 map 进行分组
// 返回 map[分组key]map[原key]value
func Group[K comparable, V any, G comparable](m map[K]V, keyFn func(K, V) G) map[G]map[K]V {
	result := make(map[G]map[K]V)
	for k, v := range m {
		g := keyFn(k, v)
		if _, ok := result[g]; !ok {
			result[g] = make(map[K]V)
		}
		result[g][k] = v
	}
	return result
}

// Filter 返回满足 predicate 的元素集合
func Filter[M ~map[K]V, K comparable, V any](m M, predicate func(K, V) bool) M {
	result := make(M)
	for k, v := range m {
		if predicate(k, v) {
			result[k] = v
		}
	}
	return result
}

// Partition 根据 predicate 将 map 分为两部分
// true → left，false → right
func Partition[M ~map[K]V, K comparable, V any](m M, predicate func(K, V) bool) (M, M) {
	left := make(M)
	right := make(M)
	for k, v := range m {
		if predicate(k, v) {
			left[k] = v
		} else {
			right[k] = v
		}
	}
	return left, right
}

// MapKeys 对 map 的 key 进行映射，value 保持不变
// 若 key 冲突，后写入的值会覆盖
func MapKeys[M ~map[K]V, K comparable, V any, NK comparable](m M, fn func(K) NK) map[NK]V {
	result := make(map[NK]V, len(m))
	for k, v := range m {
		result[fn(k)] = v
	}
	return result
}

// MapValues 对 map 的 value 进行映射，key 保持不变
func MapValues[M ~map[K]V, K comparable, V any, NV any](m M, fn func(V) NV) map[K]NV {
	result := make(map[K]NV, len(m))
	for k, v := range m {
		result[k] = fn(v)
	}
	return result
}

// Reduce 对 map 进行归约，返回结果
func Reduce[M ~map[K]V, K comparable, V any, R any](
	m M,
	init R,
	fn func(R, K, V) R,
) R {
	result := init
	for k, v := range m {
		result = fn(result, k, v)
	}
	return result
}

// Diff 比较两个 map 的差异
// added   : b 中新增的 key
// removed : a 中被删除的 key
// changed : key 存在于 a 和 b 中，但 value 不相等（使用 b 的值）
func Diff[M ~map[K]V, K comparable, V comparable](
	a, b M,
) (added, removed, changed M) {

	added = make(M)
	removed = make(M)
	changed = make(M)

	// 遍历 a：找 removed / changed
	for k, av := range a {
		if bv, ok := b[k]; !ok {
			removed[k] = av
		} else if av != bv {
			changed[k] = bv
		}
	}

	// 遍历 b：找 added
	for k, bv := range b {
		if _, ok := a[k]; !ok {
			added[k] = bv
		}
	}

	return
}

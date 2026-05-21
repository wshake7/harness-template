package collection

import (
	"fmt"
	"iter"
)

// Deque 双端队列，使用环形缓冲区实现
type Deque[E any] struct {
	data []E
	head int // 第一个元素的索引
	tail int // 最后一个元素的索引
	len  int // 当前元素数量
}

// New 创建指定容量的 Deque
func New[E any](capacity int) Deque[E] {
	if capacity < 0 {
		capacity = 0
	}
	return Deque[E]{
		data: make([]E, capacity),
		head: 0,
		tail: 0,
		len:  0,
	}
}

// NewFromSlice 从切片创建 Deque
func NewFromSlice[E any](slice []E) Deque[E] {
	if len(slice) == 0 {
		return New[E](0)
	}
	data := make([]E, len(slice))
	copy(data, slice)
	return Deque[E]{
		data: data,
		head: 0,
		tail: len(slice) - 1,
		len:  len(slice),
	}
}

// Len 返回元素数量
func (d *Deque[E]) Len() int {
	return d.len
}

// Cap 返回当前容量
func (d *Deque[E]) Cap() int {
	return cap(d.data)
}

// Empty 判断是否为空
func (d *Deque[E]) Empty() bool {
	return d.len == 0
}

// PushFront 在前端添加元素
func (d *Deque[E]) PushFront(value E) {
	if d.len == cap(d.data) {
		d.grow()
	}

	if d.len == 0 {
		// 首次添加元素
		d.data[0] = value
		d.head = 0
		d.tail = 0
		d.len = 1
	} else {
		// head 向前移动
		d.head = (d.head - 1 + cap(d.data)) % cap(d.data)
		d.data[d.head] = value
		d.len++
	}
}

// PushBack 在后端添加元素
func (d *Deque[E]) PushBack(value E) {
	if d.len == cap(d.data) {
		d.grow()
	}

	if d.len == 0 {
		// 首次添加元素
		d.data[0] = value
		d.head = 0
		d.tail = 0
		d.len = 1
	} else {
		// tail 向后移动
		d.tail = (d.tail + 1) % cap(d.data)
		d.data[d.tail] = value
		d.len++
	}
}

// PopFront 从前端弹出元素
func (d *Deque[E]) PopFront() (E, bool) {
	if d.len == 0 {
		var zero E
		return zero, false
	}

	value := d.data[d.head]

	// 清理引用，帮助 GC
	var zero E
	d.data[d.head] = zero

	d.len--
	if d.len > 0 {
		d.head = (d.head + 1) % cap(d.data)
	}

	return value, true
}

// PopBack 从后端弹出元素
func (d *Deque[E]) PopBack() (E, bool) {
	if d.len == 0 {
		var zero E
		return zero, false
	}

	value := d.data[d.tail]

	// 清理引用，帮助 GC
	var zero E
	d.data[d.tail] = zero

	d.len--
	if d.len > 0 {
		d.tail = (d.tail - 1 + cap(d.data)) % cap(d.data)
	}

	return value, true
}

// Front 查看前端元素
func (d *Deque[E]) Front() (E, bool) {
	if d.len == 0 {
		var zero E
		return zero, false
	}
	return d.data[d.head], true
}

// Back 查看后端元素
func (d *Deque[E]) Back() (E, bool) {
	if d.len == 0 {
		var zero E
		return zero, false
	}
	return d.data[d.tail], true
}

// Get 通过索引获取元素（0 为第一个元素）
func (d *Deque[E]) Get(index int) (E, bool) {
	if index < 0 || index >= d.len {
		var zero E
		return zero, false
	}
	realIdx := (d.head + index) % cap(d.data)
	return d.data[realIdx], true
}

// Set 通过索引设置元素
func (d *Deque[E]) Set(index int, value E) bool {
	if index < 0 || index >= d.len {
		return false
	}
	realIdx := (d.head + index) % cap(d.data)
	d.data[realIdx] = value
	return true
}

// Clear 清空队列
func (d *Deque[E]) Clear() {
	// 清理所有元素引用
	var zero E
	for i := 0; i < d.len; i++ {
		realIdx := (d.head + i) % cap(d.data)
		d.data[realIdx] = zero
	}
	d.head = 0
	d.tail = 0
	d.len = 0
}

// ForEach 遍历所有元素
func (d *Deque[E]) ForEach(fn func(E)) {
	for i := 0; i < d.len; i++ {
		realIdx := (d.head + i) % cap(d.data)
		fn(d.data[realIdx])
	}
}

// ForEachWithIndex 遍历所有元素（带索引）
func (d *Deque[E]) ForEachWithIndex(fn func(int, E)) {
	for i := 0; i < d.len; i++ {
		realIdx := (d.head + i) % cap(d.data)
		fn(i, d.data[realIdx])
	}
}

// ToSeq 转换为迭代器
func (d *Deque[E]) ToSeq() iter.Seq[E] {
	return func(yield func(E) bool) {
		for i := 0; i < d.len; i++ {
			realIdx := (d.head + i) % cap(d.data)
			if !yield(d.data[realIdx]) {
				return
			}
		}
	}
}

// ToSeq2 转换为带索引的迭代器
func (d *Deque[E]) ToSeq2() iter.Seq2[int, E] {
	return func(yield func(int, E) bool) {
		for i := 0; i < d.len; i++ {
			realIdx := (d.head + i) % cap(d.data)
			if !yield(i, d.data[realIdx]) {
				return
			}
		}
	}
}

// ToSlice 转换为切片
func (d *Deque[E]) ToSlice() []E {
	if d.len == 0 {
		return []E{}
	}

	result := make([]E, d.len)
	for i := 0; i < d.len; i++ {
		realIdx := (d.head + i) % cap(d.data)
		result[i] = d.data[realIdx]
	}
	return result
}

// String 返回字符串表示
func (d *Deque[E]) String() string {
	return fmt.Sprintf("Deque%v", d.ToSlice())
}

// Reverse 反转队列
func (d *Deque[E]) Reverse() {
	if d.len <= 1 {
		return
	}

	for i := 0; i < d.len/2; i++ {
		leftIdx := (d.head + i) % cap(d.data)
		rightIdx := (d.head + d.len - 1 - i) % cap(d.data)
		d.data[leftIdx], d.data[rightIdx] = d.data[rightIdx], d.data[leftIdx]
	}
}

// Clone 克隆队列
func (d *Deque[E]) Clone() Deque[E] {
	newData := make([]E, cap(d.data))
	copy(newData, d.data)
	return Deque[E]{
		data: newData,
		head: d.head,
		tail: d.tail,
		len:  d.len,
	}
}

// Shrink 收缩容量到实际使用大小
func (d *Deque[E]) Shrink() {
	if d.len == 0 {
		d.data = make([]E, 0)
		d.head = 0
		d.tail = 0
		return
	}

	if d.len == cap(d.data) {
		return
	}

	newData := make([]E, d.len)
	for i := 0; i < d.len; i++ {
		realIdx := (d.head + i) % cap(d.data)
		newData[i] = d.data[realIdx]
	}
	d.data = newData
	d.head = 0
	d.tail = d.len - 1
}

// grow 扩容
func (d *Deque[E]) grow() {
	var newCap int
	oldCap := cap(d.data)

	if oldCap == 0 {
		newCap = 8 // 初始容量
	} else if oldCap < 1024 {
		newCap = oldCap * 2
	} else {
		newCap = oldCap + oldCap/4 // 增长 25%
	}

	newData := make([]E, newCap)

	// 复制元素到新数组（按顺序）
	for i := 0; i < d.len; i++ {
		realIdx := (d.head + i) % oldCap
		newData[i] = d.data[realIdx]
	}

	d.data = newData
	d.head = 0
	if d.len > 0 {
		d.tail = d.len - 1
	} else {
		d.tail = 0
	}
}

// Contains 检查是否包含元素
func (d *Deque[E]) Contains(value E, equal func(E, E) bool) bool {
	for i := 0; i < d.len; i++ {
		realIdx := (d.head + i) % cap(d.data)
		if equal(d.data[realIdx], value) {
			return true
		}
	}
	return false
}

// Filter 过滤元素
func (d *Deque[E]) Filter(predicate func(E) bool) Deque[E] {
	result := New[E](d.len)
	for i := 0; i < d.len; i++ {
		realIdx := (d.head + i) % cap(d.data)
		if predicate(d.data[realIdx]) {
			result.PushBack(d.data[realIdx])
		}
	}
	return result
}

// Map 映射元素
func Map[E any, T any](d *Deque[E], mapper func(E) T) Deque[T] {
	result := New[T](d.len)
	for i := 0; i < d.len; i++ {
		realIdx := (d.head + i) % cap(d.data)
		result.PushBack(mapper(d.data[realIdx]))
	}
	return result
}

// Rotate 旋转队列（正数向右，负数向左）
func (d *Deque[E]) Rotate(n int) {
	if d.len <= 1 {
		return
	}

	// 标准化旋转次数
	n = n % d.len
	if n < 0 {
		n += d.len
	}

	if n == 0 {
		return
	}

	// 向右旋转 n 次 = 将后 n 个元素移到前面
	for i := 0; i < n; i++ {
		if val, ok := d.PopBack(); ok {
			d.PushFront(val)
		}
	}
}

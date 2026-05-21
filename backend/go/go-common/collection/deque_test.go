package collection

import (
	"testing"
)

// TestDequeBasicOperations 测试基本操作
func TestDequeBasicOperations(t *testing.T) {
	d := New[int](4)

	// 测试空队列
	if !d.Empty() {
		t.Error("New deque should be empty")
	}

	if d.Len() != 0 {
		t.Errorf("Expected length 0, got %d", d.Len())
	}

	// 测试 PushBack
	d.PushBack(1)
	d.PushBack(2)
	d.PushBack(3)

	if d.Len() != 3 {
		t.Errorf("Expected length 3, got %d", d.Len())
	}

	// 测试 Front 和 Back
	if front, ok := d.Front(); !ok || front != 1 {
		t.Errorf("Expected front 1, got %d", front)
	}

	if back, ok := d.Back(); !ok || back != 3 {
		t.Errorf("Expected back 3, got %d", back)
	}
}

// TestDequePushFront 测试前端插入
func TestDequePushFront(t *testing.T) {
	d := New[int](4)

	d.PushFront(3)
	d.PushFront(2)
	d.PushFront(1)

	expected := []int{1, 2, 3}
	result := d.ToSlice()

	if len(result) != len(expected) {
		t.Fatalf("Expected length %d, got %d", len(expected), len(result))
	}

	for i, v := range expected {
		if result[i] != v {
			t.Errorf("Expected result[%d] = %d, got %d", i, v, result[i])
		}
	}
}

// TestDequePushBack 测试后端插入
func TestDequePushBack(t *testing.T) {
	d := New[int](4)

	d.PushBack(1)
	d.PushBack(2)
	d.PushBack(3)

	expected := []int{1, 2, 3}
	result := d.ToSlice()

	if len(result) != len(expected) {
		t.Fatalf("Expected length %d, got %d", len(expected), len(result))
	}

	for i, v := range expected {
		if result[i] != v {
			t.Errorf("Expected result[%d] = %d, got %d", i, v, result[i])
		}
	}
}

// TestDequePopFront 测试前端弹出
func TestDequePopFront(t *testing.T) {
	d := New[int](4)
	d.PushBack(1)
	d.PushBack(2)
	d.PushBack(3)

	val, ok := d.PopFront()
	if !ok || val != 1 {
		t.Errorf("Expected PopFront to return 1, got %d", val)
	}

	val, ok = d.PopFront()
	if !ok || val != 2 {
		t.Errorf("Expected PopFront to return 2, got %d", val)
	}

	if d.Len() != 1 {
		t.Errorf("Expected length 1, got %d", d.Len())
	}

	val, ok = d.PopFront()
	if !ok || val != 3 {
		t.Errorf("Expected PopFront to return 3, got %d", val)
	}

	if !d.Empty() {
		t.Error("Deque should be empty after popping all elements")
	}

	// 测试空队列弹出
	_, ok = d.PopFront()
	if ok {
		t.Error("PopFront on empty deque should return false")
	}
}

// TestDequePopBack 测试后端弹出
func TestDequePopBack(t *testing.T) {
	d := New[int](4)
	d.PushBack(1)
	d.PushBack(2)
	d.PushBack(3)

	val, ok := d.PopBack()
	if !ok || val != 3 {
		t.Errorf("Expected PopBack to return 3, got %d", val)
	}

	val, ok = d.PopBack()
	if !ok || val != 2 {
		t.Errorf("Expected PopBack to return 2, got %d", val)
	}

	if d.Len() != 1 {
		t.Errorf("Expected length 1, got %d", d.Len())
	}
}

// TestDequeMixedOperations 测试混合操作
func TestDequeMixedOperations(t *testing.T) {
	d := New[int](4)

	d.PushBack(2)
	d.PushFront(1)
	d.PushBack(3)
	d.PushFront(0)

	expected := []int{0, 1, 2, 3}
	result := d.ToSlice()

	if len(result) != len(expected) {
		t.Fatalf("Expected length %d, got %d", len(expected), len(result))
	}

	for i, v := range expected {
		if result[i] != v {
			t.Errorf("Expected result[%d] = %d, got %d", i, v, result[i])
		}
	}

	// 弹出操作
	d.PopFront()
	d.PopBack()

	expected = []int{1, 2}
	result = d.ToSlice()

	if len(result) != len(expected) {
		t.Fatalf("Expected length %d, got %d", len(expected), len(result))
	}

	for i, v := range expected {
		if result[i] != v {
			t.Errorf("Expected result[%d] = %d, got %d", i, v, result[i])
		}
	}
}

// TestDequeAutoGrow 测试自动扩容
func TestDequeAutoGrow(t *testing.T) {
	d := New[int](2)

	// 添加超过初始容量的元素
	for i := range 10 {
		d.PushBack(i)
	}

	if d.Len() != 10 {
		t.Errorf("Expected length 10, got %d", d.Len())
	}

	if d.Cap() < 10 {
		t.Errorf("Expected capacity >= 10, got %d", d.Cap())
	}

	// 验证所有元素
	for i := range 10 {
		if val, ok := d.Get(i); !ok || val != i {
			t.Errorf("Expected Get(%d) = %d, got %d", i, i, val)
		}
	}
}

// TestDequeGetSet 测试 Get 和 Set
func TestDequeGetSet(t *testing.T) {
	d := New[int](10)

	for i := range 5 {
		d.PushBack(i)
	}

	// 测试 Get
	for i := range 5 {
		if val, ok := d.Get(i); !ok || val != i {
			t.Errorf("Expected Get(%d) = %d, got %d", i, i, val)
		}
	}

	// 测试越界 Get
	if _, ok := d.Get(10); ok {
		t.Error("Get with out of bounds index should return false")
	}

	if _, ok := d.Get(-1); ok {
		t.Error("Get with negative index should return false")
	}

	// 测试 Set
	d.Set(2, 100)
	if val, ok := d.Get(2); !ok || val != 100 {
		t.Errorf("Expected Get(2) = 100 after Set, got %d", val)
	}

	// 测试越界 Set
	if ok := d.Set(10, 200); ok {
		t.Error("Set with out of bounds index should return false")
	}
}

// TestDequeForEach 测试遍历
func TestDequeForEach(t *testing.T) {
	d := New[int](10)

	expected := []int{1, 2, 3, 4, 5}
	for _, v := range expected {
		d.PushBack(v)
	}

	result := []int{}
	d.ForEach(func(v int) {
		result = append(result, v)
	})

	if len(result) != len(expected) {
		t.Fatalf("Expected length %d, got %d", len(expected), len(result))
	}

	for i, v := range expected {
		if result[i] != v {
			t.Errorf("Expected result[%d] = %d, got %d", i, v, result[i])
		}
	}
}

// TestDequeForEachWithIndex 测试带索引的遍历
func TestDequeForEachWithIndex(t *testing.T) {
	d := New[int](10)

	for i := range 5 {
		d.PushBack(i * 10)
	}

	d.ForEachWithIndex(func(idx, val int) {
		expected := idx * 10
		if val != expected {
			t.Errorf("Expected value at index %d to be %d, got %d", idx, expected, val)
		}
	})
}

// TestDequeToSeq 测试迭代器
func TestDequeToSeq(t *testing.T) {
	d := New[int](10)

	expected := []int{1, 2, 3, 4, 5}
	for _, v := range expected {
		d.PushBack(v)
	}

	result := []int{}
	for v := range d.ToSeq() {
		result = append(result, v)
	}

	if len(result) != len(expected) {
		t.Fatalf("Expected length %d, got %d", len(expected), len(result))
	}

	for i, v := range expected {
		if result[i] != v {
			t.Errorf("Expected result[%d] = %d, got %d", i, v, result[i])
		}
	}
}

// TestDequeToSeq2 测试带索引的迭代器
func TestDequeToSeq2(t *testing.T) {
	d := New[int](10)

	expected := []int{10, 20, 30}
	for _, v := range expected {
		d.PushBack(v)
	}

	for idx, val := range d.ToSeq2() {
		if val != expected[idx] {
			t.Errorf("Expected value at index %d to be %d, got %d", idx, expected[idx], val)
		}
	}
}

// TestDequeClear 测试清空
func TestDequeClear(t *testing.T) {
	d := New[int](10)

	for i := range 5 {
		d.PushBack(i)
	}

	d.Clear()

	if !d.Empty() {
		t.Error("Deque should be empty after Clear")
	}

	if d.Len() != 0 {
		t.Errorf("Expected length 0 after Clear, got %d", d.Len())
	}

	// 验证可以继续使用
	d.PushBack(100)
	if val, ok := d.Front(); !ok || val != 100 {
		t.Error("Deque should be usable after Clear")
	}
}

// TestDequeReverse 测试反转
func TestDequeReverse(t *testing.T) {
	d := New[int](10)

	original := []int{1, 2, 3, 4, 5}
	for _, v := range original {
		d.PushBack(v)
	}

	d.Reverse()

	expected := []int{5, 4, 3, 2, 1}
	result := d.ToSlice()

	if len(result) != len(expected) {
		t.Fatalf("Expected length %d, got %d", len(expected), len(result))
	}

	for i, v := range expected {
		if result[i] != v {
			t.Errorf("Expected result[%d] = %d, got %d", i, v, result[i])
		}
	}
}

// TestDequeClone 测试克隆
func TestDequeClone(t *testing.T) {
	d := New[int](10)

	for i := range 5 {
		d.PushBack(i)
	}

	clone := d.Clone()

	// 验证内容相同
	if clone.Len() != d.Len() {
		t.Error("Clone should have same length as original")
	}

	for i := 0; i < d.Len(); i++ {
		origVal, _ := d.Get(i)
		cloneVal, _ := clone.Get(i)
		if origVal != cloneVal {
			t.Errorf("Clone value at index %d differs from original", i)
		}
	}

	// 修改原始队列，验证克隆不受影响
	d.PushBack(100)

	if clone.Len() == d.Len() {
		t.Error("Clone should be independent of original")
	}
}

// TestDequeShrink 测试收缩
func TestDequeShrink(t *testing.T) {
	d := New[int](100)

	for i := range 10 {
		d.PushBack(i)
	}

	oldCap := d.Cap()
	d.Shrink()
	newCap := d.Cap()

	if newCap >= oldCap {
		t.Errorf("Expected capacity to shrink from %d, got %d", oldCap, newCap)
	}

	if d.Len() != 10 {
		t.Errorf("Expected length 10 after shrink, got %d", d.Len())
	}

	// 验证数据完整性
	for i := range 10 {
		if val, ok := d.Get(i); !ok || val != i {
			t.Errorf("Data corrupted after shrink at index %d", i)
		}
	}
}

// TestDequeContains 测试包含检查
func TestDequeContains(t *testing.T) {
	d := New[int](10)

	d.PushBack(1)
	d.PushBack(2)
	d.PushBack(3)

	equal := func(a, b int) bool { return a == b }

	if !d.Contains(2, equal) {
		t.Error("Deque should contain 2")
	}

	if d.Contains(5, equal) {
		t.Error("Deque should not contain 5")
	}
}

// TestDequeFilter 测试过滤
func TestDequeFilter(t *testing.T) {
	d := New[int](10)

	for i := 1; i <= 10; i++ {
		d.PushBack(i)
	}

	// 过滤出偶数
	filtered := d.Filter(func(v int) bool {
		return v%2 == 0
	})

	expected := []int{2, 4, 6, 8, 10}
	result := filtered.ToSlice()

	if len(result) != len(expected) {
		t.Fatalf("Expected length %d, got %d", len(expected), len(result))
	}

	for i, v := range expected {
		if result[i] != v {
			t.Errorf("Expected result[%d] = %d, got %d", i, v, result[i])
		}
	}
}

// TestDequeMap 测试映射
func TestDequeMap(t *testing.T) {
	d := New[int](10)

	for i := 1; i <= 5; i++ {
		d.PushBack(i)
	}

	// 映射为平方
	mapped := Map(&d, func(v int) int {
		return v * v
	})

	expected := []int{1, 4, 9, 16, 25}
	result := mapped.ToSlice()

	if len(result) != len(expected) {
		t.Fatalf("Expected length %d, got %d", len(expected), len(result))
	}

	for i, v := range expected {
		if result[i] != v {
			t.Errorf("Expected result[%d] = %d, got %d", i, v, result[i])
		}
	}
}

// TestDequeRotate 测试旋转
func TestDequeRotate(t *testing.T) {
	d := New[int](10)

	for i := 1; i <= 5; i++ {
		d.PushBack(i)
	}

	// 向右旋转 2 次
	d.Rotate(2)

	expected := []int{4, 5, 1, 2, 3}
	result := d.ToSlice()

	if len(result) != len(expected) {
		t.Fatalf("Expected length %d, got %d", len(expected), len(result))
	}

	for i, v := range expected {
		if result[i] != v {
			t.Errorf("Expected result[%d] = %d after rotate, got %d", i, v, result[i])
		}
	}

	// 向左旋转 1 次
	d.Rotate(-1)

	expected = []int{5, 1, 2, 3, 4}
	result = d.ToSlice()

	for i, v := range expected {
		if result[i] != v {
			t.Errorf("Expected result[%d] = %d after rotate, got %d", i, v, result[i])
		}
	}
}

// TestDequeNewFromSlice 测试从切片创建
func TestDequeNewFromSlice(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	d := NewFromSlice(slice)

	if d.Len() != len(slice) {
		t.Errorf("Expected length %d, got %d", len(slice), d.Len())
	}

	for i, v := range slice {
		if val, ok := d.Get(i); !ok || val != v {
			t.Errorf("Expected Get(%d) = %d, got %d", i, v, val)
		}
	}
}

// TestDequeWrapAround 测试环形缓冲区边界情况
func TestDequeWrapAround(t *testing.T) {
	d := New[int](4)

	// 填满队列
	d.PushBack(1)
	d.PushBack(2)
	d.PushBack(3)
	d.PushBack(4)

	// 弹出两个
	d.PopFront()
	d.PopFront()

	// 再添加两个（会发生环绕）
	d.PushBack(5)
	d.PushBack(6)

	expected := []int{3, 4, 5, 6}
	result := d.ToSlice()

	if len(result) != len(expected) {
		t.Fatalf("Expected length %d, got %d", len(expected), len(result))
	}

	for i, v := range expected {
		if result[i] != v {
			t.Errorf("Expected result[%d] = %d, got %d", i, v, result[i])
		}
	}
}

// BenchmarkDequePushBack 性能测试：后端插入
func BenchmarkDequePushBack(b *testing.B) {
	d := New[int](0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d.PushBack(i)
	}
}

// BenchmarkDequePushFront 性能测试：前端插入
func BenchmarkDequePushFront(b *testing.B) {
	d := New[int](0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d.PushFront(i)
	}
}

// BenchmarkDequePopBack 性能测试：后端弹出
func BenchmarkDequePopBack(b *testing.B) {
	d := New[int](b.N)
	for i := 0; i < b.N; i++ {
		d.PushBack(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d.PopBack()
	}
}

// BenchmarkDequePopFront 性能测试：前端弹出
func BenchmarkDequePopFront(b *testing.B) {
	d := New[int](b.N)
	for i := 0; i < b.N; i++ {
		d.PushBack(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d.PopFront()
	}
}

// BenchmarkDequeGet 性能测试：随机访问
func BenchmarkDequeGet(b *testing.B) {
	d := New[int](1000)
	for i := range 1000 {
		d.PushBack(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d.Get(i % 1000)
	}
}

// BenchmarkDequeToSlice 性能测试：转换为切片
func BenchmarkDequeToSlice(b *testing.B) {
	d := New[int](1000)
	for i := range 1000 {
		d.PushBack(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = d.ToSlice()
	}
}

package collection

import (
	"go-common/types"
	"testing"
)

// TestBitMapBasicOperations 测试 BitMap 基本操作
func TestBitMapBasicOperations(t *testing.T) {
	t.Run("uint8", func(t *testing.T) {
		bm := BitMapNew[uint8](0)

		// 测试 Set 和 Get
		bm.Set(0, true)
		bm.Set(3, true)
		bm.Set(7, true)

		if !bm.Get(0) {
			t.Error("Expected bit 0 to be true")
		}
		if !bm.Get(3) {
			t.Error("Expected bit 3 to be true")
		}
		if !bm.Get(7) {
			t.Error("Expected bit 7 to be true")
		}
		if bm.Get(1) {
			t.Error("Expected bit 1 to be false")
		}

		// 测试 Value
		expected := uint8(0b10001001) // bits 0, 3, 7
		if bm.Value() != expected {
			t.Errorf("Expected value %d, got %d", expected, bm.Value())
		}

		// 测试 Count
		if count := bm.Count(); count != 3 {
			t.Errorf("Expected count 3, got %d", count)
		}
	})

	t.Run("uint64", func(t *testing.T) {
		bm := BitMapNew[uint64](0)

		bm.Set(0, true)
		bm.Set(31, true)
		bm.Set(63, true)

		if count := bm.Count(); count != 3 {
			t.Errorf("Expected count 3, got %d", count)
		}

		// 测试边界
		bm.Set(64, true) // 超出范围，应该被忽略
		if count := bm.Count(); count != 3 {
			t.Errorf("Expected count 3 after out of bounds set, got %d", count)
		}
	})
}

// TestBitMapSetUnset 测试位的设置和取消
func TestBitMapSetUnset(t *testing.T) {
	bm := BitMapNew[uint32](0)

	// 设置位
	bm.Set(5, true)
	if !bm.Get(5) {
		t.Error("Expected bit 5 to be true")
	}

	// 取消位
	bm.Set(5, false)
	if bm.Get(5) {
		t.Error("Expected bit 5 to be false")
	}

	// 多次设置同一位
	bm.Set(10, true)
	bm.Set(10, true)
	if count := bm.Count(); count != 1 {
		t.Errorf("Expected count 1, got %d", count)
	}
}

// TestBitMapBoundaryConditions 测试边界条件
func TestBitMapBoundaryConditions(t *testing.T) {
	bm := BitMapNew[uint16](0)

	// 负数索引
	bm.Set(-1, true)
	if bm.Get(-1) {
		t.Error("Negative index should return false")
	}

	// 超出范围
	bm.Set(16, true) // uint16 只有 16 位
	if bm.Get(16) {
		t.Error("Out of bounds index should return false")
	}

	// 最大有效索引
	bm.Set(15, true)
	if !bm.Get(15) {
		t.Error("Expected bit 15 to be true")
	}
}

// TestBytesBitMapBasicOperations 测试 BytesBitMap 基本操作
func TestBytesBitMapBasicOperations(t *testing.T) {
	bm := BytesBitMapNew(make([]byte, 2))

	// 测试 Set 和 Get
	bm.Set(0, true)
	bm.Set(7, true)
	bm.Set(8, true)
	bm.Set(15, true)

	if !bm.Get(0) || !bm.Get(7) || !bm.Get(8) || !bm.Get(15) {
		t.Error("Expected bits to be set correctly")
	}

	// 测试 Count
	if count := bm.Count(); count != 4 {
		t.Errorf("Expected count 4, got %d", count)
	}

	// 测试 Len
	if length := bm.Len(); length != 16 {
		t.Errorf("Expected length 16, got %d", length)
	}
}

// TestBytesBitMapAutoGrow 测试自动扩容
func TestBytesBitMapAutoGrow(t *testing.T) {
	bm := BytesBitMapNew([]byte{})

	// 设置超出初始容量的位
	bm.Set(100, true)

	if !bm.Get(100) {
		t.Error("Expected bit 100 to be true after auto-grow")
	}

	if len(bm.Value()) < 13 { // 100/8 + 1 = 13
		t.Errorf("Expected at least 13 bytes, got %d", len(bm.Value()))
	}

	// 验证之前未设置的位仍然是 false
	for i := range 100 {
		if bm.Get(i) {
			t.Errorf("Expected bit %d to be false", i)
		}
	}
}

// TestBytesBitMapSetUnset 测试位的设置和取消
func TestBytesBitMapSetUnset(t *testing.T) {
	bm := BytesBitMapNew(make([]byte, 4))

	indices := []int{0, 7, 8, 15, 16, 23, 24, 31}

	// 设置所有位
	for _, idx := range indices {
		bm.Set(idx, true)
	}

	// 验证所有位都已设置
	for _, idx := range indices {
		if !bm.Get(idx) {
			t.Errorf("Expected bit %d to be true", idx)
		}
	}

	if count := bm.Count(); count != len(indices) {
		t.Errorf("Expected count %d, got %d", len(indices), count)
	}

	// 取消部分位
	bm.Set(7, false)
	bm.Set(15, false)

	if bm.Get(7) || bm.Get(15) {
		t.Error("Expected bits 7 and 15 to be false")
	}

	if count := bm.Count(); count != len(indices)-2 {
		t.Errorf("Expected count %d, got %d", len(indices)-2, count)
	}
}

// TestBytesBitMapForEach 测试 ForEach 遍历
func TestBytesBitMapForEach(t *testing.T) {
	bm := BytesBitMapNew(make([]byte, 2))

	setIndices := []int{0, 3, 7, 8, 15}
	for _, idx := range setIndices {
		bm.Set(idx, true)
	}

	visited := make(map[int]bool)
	bm.ForEach(func(index int, b bool) {
		visited[index] = b
	})

	// 验证所有位都被访问
	if len(visited) != 16 {
		t.Errorf("Expected 16 bits visited, got %d", len(visited))
	}

	// 验证设置的位
	for _, idx := range setIndices {
		if !visited[idx] {
			t.Errorf("Expected bit %d to be true", idx)
		}
	}
}

// TestBytesBitMapForEachSet 测试 ForEachSet 只遍历设置的位
func TestBytesBitMapForEachSet(t *testing.T) {
	bm := BytesBitMapNew(make([]byte, 10))

	setIndices := []int{0, 5, 17, 33, 64, 79}
	for _, idx := range setIndices {
		bm.Set(idx, true)
	}

	visited := []int{}
	bm.ForEachSet(func(index int) {
		visited = append(visited, index)
	})

	if len(visited) != len(setIndices) {
		t.Errorf("Expected %d bits visited, got %d", len(setIndices), len(visited))
	}

	// 验证访问的索引
	visitedMap := make(map[int]bool)
	for _, idx := range visited {
		visitedMap[idx] = true
	}

	for _, idx := range setIndices {
		if !visitedMap[idx] {
			t.Errorf("Expected index %d to be visited", idx)
		}
	}
}

// TestToMap 测试转换为 Map
func TestToMap(t *testing.T) {
	bm := BytesBitMapNew(make([]byte, 4))

	bm.Set(0, true)
	bm.Set(5, true)
	bm.Set(10, true)
	bm.Set(20, true)

	// 转换为 map，键为索引的字符串形式
	result := ToMap(&bm, func(v int) string {
		return string(rune('A' + v))
	})

	if len(result) != 4 {
		t.Errorf("Expected map length 4, got %d", len(result))
	}

	expected := map[string]types.Unit{
		string(rune('A' + 0)):  {},
		string(rune('A' + 5)):  {},
		string(rune('A' + 10)): {},
		string(rune('A' + 20)): {},
	}

	for key := range expected {
		if _, exists := result[key]; !exists {
			t.Errorf("Expected key %s to exist in result", key)
		}
	}
}

// TestToSlice 测试转换为 Slice
func TestToSlice(t *testing.T) {
	bm := BytesBitMapNew(make([]byte, 4))

	indices := []int{0, 5, 10, 15, 20}
	for _, idx := range indices {
		bm.Set(idx, true)
	}

	// 转换为 slice，值为索引的平方
	result := ToSlice(&bm, func(v int) int {
		return v * v
	})

	if len(result) != len(indices) {
		t.Errorf("Expected slice length %d, got %d", len(indices), len(result))
	}

	expected := []int{0, 25, 100, 225, 400}
	for i, val := range expected {
		if result[i] != val {
			t.Errorf("Expected result[%d] = %d, got %d", i, val, result[i])
		}
	}
}

// TestNewBytesBitMapWithCapacity 测试预分配容量的构造函数
func TestNewBytesBitMapWithCapacity(t *testing.T) {
	bitCount := 100
	bm := NewBytesBitMapWithCapacity(bitCount)

	expectedBytes := (bitCount + 7) / 8
	if len(bm.Value()) != expectedBytes {
		t.Errorf("Expected %d bytes, got %d", expectedBytes, len(bm.Value()))
	}

	// 验证可以设置所有位
	for i := range bitCount {
		bm.Set(i, true)
	}

	if count := bm.Count(); count != bitCount {
		t.Errorf("Expected count %d, got %d", bitCount, count)
	}
}

// TestBytesBitMapClear 测试清空操作
func TestBytesBitMapClear(t *testing.T) {
	bm := BytesBitMapNew(make([]byte, 4))

	// 设置一些位
	for i := 0; i < 32; i += 3 {
		bm.Set(i, true)
	}

	if count := bm.Count(); count == 0 {
		t.Error("Expected some bits to be set")
	}

	// 清空
	bm.Clear()

	if count := bm.Count(); count != 0 {
		t.Errorf("Expected count 0 after clear, got %d", count)
	}

	// 验证所有位都是 false
	for i := range 32 {
		if bm.Get(i) {
			t.Errorf("Expected bit %d to be false after clear", i)
		}
	}
}

// TestBytesBitMapClone 测试克隆操作
func TestBytesBitMapClone(t *testing.T) {
	bm := BytesBitMapNew(make([]byte, 4))

	indices := []int{0, 5, 10, 15, 20, 25}
	for _, idx := range indices {
		bm.Set(idx, true)
	}

	// 克隆
	clone := bm.Clone()

	// 验证克隆的内容相同
	if clone.Count() != bm.Count() {
		t.Error("Clone should have same count as original")
	}

	for _, idx := range indices {
		if !clone.Get(idx) {
			t.Errorf("Expected cloned bit %d to be true", idx)
		}
	}

	// 修改原始位图，验证克隆不受影响
	bm.Set(0, false)

	if !clone.Get(0) {
		t.Error("Clone should be independent of original")
	}
}

// TestBytesBitMapBoundaryConditions 测试边界条件
func TestBytesBitMapBoundaryConditions(t *testing.T) {
	bm := BytesBitMapNew(make([]byte, 2))

	// 负数索引
	bm.Set(-1, true)
	if bm.Get(-1) {
		t.Error("Negative index should return false")
	}

	// 超出当前范围但会自动扩容
	bm.Set(100, true)
	if !bm.Get(100) {
		t.Error("Expected bit 100 to be true after auto-grow")
	}

	// 访问未初始化的位
	if bm.Get(50) {
		t.Error("Expected unset bit to be false")
	}
}

// BenchmarkBitMapCount 性能测试：BitMap Count
func BenchmarkBitMapCount(b *testing.B) {
	bm := BitMapNew[uint64](0xFFFFFFFFFFFFFFFF)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = bm.Count()
	}
}

// BenchmarkBytesBitMapCount 性能测试：BytesBitMap Count
func BenchmarkBytesBitMapCount(b *testing.B) {
	data := make([]byte, 1024)
	for i := range data {
		data[i] = 0xFF
	}
	bm := BytesBitMapNew(data)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = bm.Count()
	}
}

// BenchmarkBytesBitMapSet 性能测试：Set 操作
func BenchmarkBytesBitMapSet(b *testing.B) {
	bm := NewBytesBitMapWithCapacity(10000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bm.Set(i%10000, true)
	}
}

// BenchmarkBytesBitMapForEachSet 性能测试：ForEachSet vs ForEach
func BenchmarkBytesBitMapForEachSet(b *testing.B) {
	bm := NewBytesBitMapWithCapacity(10000)

	// 设置 10% 的位
	for i := range 1000 {
		bm.Set(i*10, true)
	}

	b.Run("ForEach", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			count := 0
			bm.ForEach(func(index int, b bool) {
				if b {
					count++
				}
			})
		}
	})

	b.Run("ForEachSet", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			count := 0
			bm.ForEachSet(func(index int) {
				count++
			})
		}
	})
}

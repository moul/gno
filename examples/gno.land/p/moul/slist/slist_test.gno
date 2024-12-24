package slist

import (
	"testing"

	"gno.land/p/demo/uassert"
)

func TestNewSkipList(t *testing.T) {
	sl := NewSkipList()
	if sl.head == nil {
		t.Error("head should not be nil")
	}
	if sl.level != 1 {
		t.Errorf("initial level should be 1, got %d", sl.level)
	}
	if sl.size != 0 {
		t.Errorf("initial size should be 0, got %d", sl.size)
	}
}

func TestInsert(t *testing.T) {
	sl := NewSkipList()

	// Test single insert
	sl.Insert(5, "five")
	if val, found := sl.Search(5); !found || val != "five" {
		t.Error("failed to insert and retrieve single value")
	}

	// Test multiple inserts
	testCases := []struct {
		key   int
		value string
	}{
		{3, "three"},
		{7, "seven"},
		{1, "one"},
		{10, "ten"},
	}

	for _, tc := range testCases {
		sl.Insert(tc.key, tc.value)
		if val, found := sl.Search(tc.key); !found || val != tc.value {
			t.Errorf("failed to insert and retrieve value for key %d", tc.key)
		}
	}

	// Test overwrite existing key
	sl.Insert(5, "FIVE")
	if val, found := sl.Search(5); !found || val != "FIVE" {
		t.Error("failed to overwrite existing key")
	}
}

func TestSearch(t *testing.T) {
	sl := NewSkipList()

	// Test search in empty list
	if _, found := sl.Search(1); found {
		t.Error("search in empty list should return false")
	}

	// Insert some values
	values := map[int]string{
		1: "one",
		3: "three",
		5: "five",
	}
	for k, v := range values {
		sl.Insert(k, v)
	}

	// Test existing keys
	for k, v := range values {
		if val, found := sl.Search(k); !found || val != v {
			t.Errorf("search failed for key %d", k)
		}
	}

	// Test non-existing key
	if _, found := sl.Search(2); found {
		t.Error("search should return false for non-existing key")
	}
}

func TestDelete(t *testing.T) {
	sl := NewSkipList()

	// Test delete from empty list
	sl.Delete(1) // Should not panic

	// Insert some values
	values := []int{1, 3, 5, 7, 9}
	for _, v := range values {
		sl.Insert(v, v)
	}

	// Test deleting existing keys
	sl.Delete(5)
	if _, found := sl.Search(5); found {
		t.Error("key 5 should have been deleted")
	}

	// Verify other keys still exist
	for _, v := range []int{1, 3, 7, 9} {
		if _, found := sl.Search(v); !found {
			t.Errorf("key %d should still exist", v)
		}
	}

	// Test deleting non-existing key
	initialSize := sl.size
	sl.Delete(100)
	if sl.size != initialSize {
		t.Error("size should not change when deleting non-existing key")
	}
}

func TestCalculateLevel(t *testing.T) {
	sl := NewSkipList()

	// Test initial level
	if level := sl.calculateLevel(); level != 1 {
		t.Errorf("initial level should be 1, got %d", level)
	}

	// Insert elements and verify level increases appropriately
	testCases := []struct {
		size          int
		expectedLevel int
	}{
		{1, 1},  // 2^0
		{2, 2},  // 2^1
		{4, 3},  // 2^2
		{8, 4},  // 2^3
		{16, 5}, // 2^4
	}

	for _, tc := range testCases {
		sl.size = tc.size
		level := sl.calculateLevel()
		if level != tc.expectedLevel {
			t.Errorf("for size %d, expected level %d, got %d",
				tc.size, tc.expectedLevel, level)
		}
	}

	// Test maximum level
	sl.size = 1 << maxLevel // Very large size
	if level := sl.calculateLevel(); level > maxLevel {
		t.Errorf("level should not exceed maxLevel(%d), got %d",
			maxLevel, level)
	}
}

func TestSize(t *testing.T) {
	sl := NewSkipList()

	// Test initial size
	if sl.size != 0 {
		t.Errorf("initial size should be 0, got %d", sl.size)
	}

	// Test size after insertions
	insertions := []int{1, 2, 3, 4, 5}
	for i, v := range insertions {
		sl.Insert(v, v)
		if sl.size != i+1 {
			t.Errorf("size should be %d after %d insertions, got %d",
				i+1, i+1, sl.size)
		}
	}

	// Test size after deletions
	sl.Delete(3)
	if sl.size != 4 {
		t.Errorf("size should be 4 after deletion, got %d", sl.size)
	}

	// Test size after deleting non-existing key
	sl.Delete(100)
	if sl.size != 4 {
		t.Errorf("size should remain 4 after deleting non-existing key, got %d",
			sl.size)
	}
}

func TestSkipList_EmptyList(t *testing.T) {
	sl := NewSkipList()

	// Test empty list operations
	_, found := sl.Search(42)
	uassert.False(t, found, "found value in empty list")

	sl.Delete(42) // Should not panic
}

func TestSkipList_SingleElement(t *testing.T) {
	sl := NewSkipList()

	// Insert and verify single element
	sl.Insert(5, "five")

	val, found := sl.Search(5)
	uassert.True(t, found, "value not found")
	uassert.Equal(t, "five", val)

	// Verify level calculation for single element
	uassert.Equal(t, 1, sl.level)
}

func TestSkipList_MultipleElements(t *testing.T) {
	sl := NewSkipList()

	// Insert elements in non-sequential order
	elements := []struct {
		key   int
		value string
	}{
		{3, "three"},
		{1, "one"},
		{4, "four"},
		{2, "two"},
		{5, "five"},
	}

	for _, e := range elements {
		sl.Insert(e.key, e.value)
	}

	// Verify all elements are found
	for _, e := range elements {
		val, found := sl.Search(e.key)
		uassert.True(t, found, "value not found for key")
		uassert.Equal(t, e.value, val)
	}

	// Verify non-existent elements
	_, found := sl.Search(0)
	uassert.False(t, found, "found non-existent value")
	_, found = sl.Search(6)
	uassert.False(t, found, "found non-existent value")
}

func TestSkipList_UpdateExisting(t *testing.T) {
	sl := NewSkipList()

	// Insert and then update
	sl.Insert(1, "one")
	sl.Insert(1, "ONE")

	val, found := sl.Search(1)
	uassert.True(t, found, "value not found")
	uassert.Equal(t, "ONE", val)
}

func TestSkipList_DeleteElements(t *testing.T) {
	sl := NewSkipList()

	// Insert elements
	sl.Insert(1, "one")
	sl.Insert(2, "two")
	sl.Insert(3, "three")

	// Delete middle element
	sl.Delete(2)
	_, found := sl.Search(2)
	uassert.False(t, found, "found deleted element")

	// Verify surrounding elements still exist
	val, found := sl.Search(1)
	uassert.True(t, found, "value not found")
	uassert.Equal(t, "one", val)

	val, found = sl.Search(3)
	uassert.True(t, found, "value not found")
	uassert.Equal(t, "three", val)
}

func TestSkipList_LevelGrowth(t *testing.T) {
	sl := NewSkipList()

	// Insert elements and verify level growth
	initialLevel := sl.level

	// Insert enough elements to force level growth
	for i := 0; i < 10; i++ {
		sl.Insert(i, i)
	}

	uassert.True(t, sl.level > initialLevel, "level did not grow with insertions")
}

func TestSkipList_LargeDataSet(t *testing.T) {
	sl := NewSkipList()

	// Insert a larger number of elements
	for i := 0; i < 100; i++ {
		sl.Insert(i, i*2)
	}

	// Verify random access
	testCases := []struct {
		key      int
		expected int
	}{
		{0, 0},
		{50, 100},
		{99, 198},
	}

	for _, tc := range testCases {
		val, found := sl.Search(tc.key)
		uassert.True(t, found, "value not found for key")
		uassert.Equal(t, tc.expected, val)
	}
}

func TestSkipList_DeleteAndReinsert(t *testing.T) {
	sl := NewSkipList()

	// Insert, delete, and reinsert elements
	sl.Insert(1, "one")
	sl.Delete(1)
	sl.Insert(1, "ONE")

	val, found := sl.Search(1)
	uassert.True(t, found, "value not found after reinsertion")
	uassert.Equal(t, "ONE", val)
}

func TestSkipList_Sprint(t *testing.T) {
	sl := NewSkipList()

	// Insert some elements
	sl.Insert(1, "one")
	sl.Insert(2, "two")
	sl.Insert(3, "three")

	// Verify Sprint output is not empty
	output := sl.Sprint()
	uassert.True(t, len(output) > 0, "Sprint output is empty")
}

func TestSkipList_FrontBack(t *testing.T) {
	sl := NewSkipList()

	// Test empty list
	_, _, exists := sl.Front()
	uassert.False(t, exists, "front should not exist in empty list")
	_, _, exists = sl.Back()
	uassert.False(t, exists, "back should not exist in empty list")

	// Test single element
	sl.Insert(1, "one")
	key, val, exists := sl.Front()
	uassert.True(t, exists, "front should exist")
	uassert.Equal(t, 1, key)
	uassert.Equal(t, "one", val)

	key, val, exists = sl.Back()
	uassert.True(t, exists, "back should exist")
	uassert.Equal(t, 1, key)
	uassert.Equal(t, "one", val)

	// Test multiple elements
	sl.Insert(2, "two")
	sl.Insert(3, "three")

	key, val, exists = sl.Front()
	uassert.True(t, exists, "front should exist")
	uassert.Equal(t, 1, key)
	uassert.Equal(t, "one", val)

	key, val, exists = sl.Back()
	uassert.True(t, exists, "back should exist")
	uassert.Equal(t, 3, key)
	uassert.Equal(t, "three", val)
}

func TestSkipList_AppendPrepend(t *testing.T) {
	sl := NewSkipList()

	// Test append to empty list
	key := sl.Append("first")
	uassert.Equal(t, 0, key)

	// Test prepend
	key = sl.Prepend("before first")
	uassert.Equal(t, -1, key)

	// Test append
	key = sl.Append("second")
	uassert.Equal(t, 1, key)

	// Verify order
	values := sl.Values()
	uassert.Equal(t, 3, len(values))
	uassert.Equal(t, "before first", values[0])
	uassert.Equal(t, "first", values[1])
	uassert.Equal(t, "second", values[2])
}

func TestSkipList_ForEach(t *testing.T) {
	sl := NewSkipList()

	// Insert test data
	expected := []string{"one", "two", "three"}
	for i, v := range expected {
		sl.Insert(i, v)
	}

	// Test ForEach
	var result []string
	sl.ForEach(func(key int, value interface{}) bool {
		result = append(result, value.(string))
		return false
	})

	uassert.Equal(t, len(expected), len(result))
	for i, v := range expected {
		uassert.Equal(t, v, result[i])
	}
}

func TestSkipList_InsertBeforeAfter(t *testing.T) {
	sl := NewSkipList()

	// Test with empty list
	_, ok := sl.InsertBefore(1, "before")
	uassert.False(t, ok, "should not insert before in empty list")

	_, ok = sl.InsertAfter(1, "after")
	uassert.False(t, ok, "should not insert after in empty list")

	// Insert initial element
	sl.Insert(10, "ten")

	// Test InsertBefore
	key, ok := sl.InsertBefore(10, "nine")
	uassert.True(t, ok, "should insert before")
	uassert.True(t, key < 10, "key should be less than 10")

	// Test InsertAfter
	key, ok = sl.InsertAfter(10, "eleven")
	uassert.True(t, ok, "should insert after")
	uassert.True(t, key > 10, "key should be greater than 10")

	// Verify order
	values := sl.Values()
	uassert.Equal(t, 3, len(values))
	uassert.Equal(t, "nine", values[0])
	uassert.Equal(t, "ten", values[1])
	uassert.Equal(t, "eleven", values[2])
}

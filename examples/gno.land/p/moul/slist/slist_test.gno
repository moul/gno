package slist

import (
	"testing"

	"gno.land/p/demo/uassert"
)

func TestNew(t *testing.T) {
	l := New()
	uassert.True(t, l.head != nil, "head should not be nil")
	uassert.Equal(t, 1, l.level, "initial level should be 1")
	uassert.Equal(t, 0, l.size, "initial size should be 0")
	uassert.Equal(t, 16, l.maxLevel, "default maxLevel should be 16")
}

func TestNewWithMaxLevel(t *testing.T) {
	l := NewWithMaxLevel(8)
	uassert.True(t, l.head != nil, "head should not be nil")
	uassert.Equal(t, 1, l.level, "initial level should be 1")
	uassert.Equal(t, 0, l.size, "initial size should be 0")
	uassert.Equal(t, 8, l.maxLevel, "maxLevel should be 8")
}

func TestInsert(t *testing.T) {
	l := New()

	// Test single insert
	l.Insert(5, "five")
	val, found := l.Search(5)
	uassert.True(t, found, "failed to find inserted value")
	uassert.Equal(t, "five", val)

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
		l.Insert(tc.key, tc.value)
		val, found := l.Search(tc.key)
		uassert.True(t, found, "failed to find inserted value")
		uassert.Equal(t, tc.value, val)
	}

	// Test overwrite
	l.Insert(5, "FIVE")
	val, found = l.Search(5)
	uassert.True(t, found, "failed to find updated value")
	uassert.Equal(t, "FIVE", val)
}

func TestSearch(t *testing.T) {
	l := New()

	// Test empty list
	_, found := l.Search(1)
	uassert.False(t, found, "found value in empty list")

	// Test with values
	values := map[int]string{
		1: "one",
		3: "three",
		5: "five",
	}

	for k, v := range values {
		l.Insert(k, v)
	}

	for k, v := range values {
		val, found := l.Search(k)
		uassert.True(t, found, "failed to find value")
		uassert.Equal(t, v, val)
	}

	// Test non-existent key
	_, found = l.Search(2)
	uassert.False(t, found, "found non-existent value")
}

func TestDelete(t *testing.T) {
	l := New()

	// Test delete on empty list
	l.Delete(1)

	// Test with values
	values := []int{1, 3, 5, 7, 9}
	for _, v := range values {
		l.Insert(v, v)
	}

	l.Delete(5)
	_, found := l.Search(5)
	uassert.False(t, found, "found deleted value")

	// Test remaining values
	for _, v := range []int{1, 3, 7, 9} {
		_, found := l.Search(v)
		uassert.True(t, found, "missing value after delete")
	}

	// Test delete non-existent
	initialSize := l.size
	l.Delete(100)
	uassert.Equal(t, initialSize, l.size, "size changed after deleting non-existent key")
}

func TestList_EmptyList(t *testing.T) {
	l := New()

	_, found := l.Search(42)
	uassert.False(t, found, "found value in empty list")

	l.Delete(42) // Should not panic
}

func TestList_FrontBack(t *testing.T) {
	l := New()

	// Test empty list
	_, _, exists := l.Front()
	uassert.False(t, exists, "front should not exist in empty list")
	_, _, exists = l.Back()
	uassert.False(t, exists, "back should not exist in empty list")

	// Test single element
	l.Insert(1, "one")
	key, val, exists := l.Front()
	uassert.True(t, exists, "front should exist")
	uassert.Equal(t, 1, key)
	uassert.Equal(t, "one", val)

	key, val, exists = l.Back()
	uassert.True(t, exists, "back should exist")
	uassert.Equal(t, 1, key)
	uassert.Equal(t, "one", val)

	// Test multiple elements
	l.Insert(2, "two")
	l.Insert(3, "three")

	key, val, exists = l.Front()
	uassert.True(t, exists, "front should exist")
	uassert.Equal(t, 1, key)
	uassert.Equal(t, "one", val)

	key, val, exists = l.Back()
	uassert.True(t, exists, "back should exist")
	uassert.Equal(t, 3, key)
	uassert.Equal(t, "three", val)
}

func TestList_AppendPrepend(t *testing.T) {
	l := New()

	// Test append to empty list
	key := l.Append("first")
	uassert.Equal(t, 0, key)

	// Test prepend
	key = l.Prepend("before first")
	uassert.True(t, key < 0, "prepend key should be negative")

	// Test append
	key = l.Append("second")
	uassert.True(t, key > 0, "append key should be positive")

	// Verify order
	values := l.Values()
	uassert.Equal(t, 3, len(values))
	uassert.Equal(t, "before first", values[0])
	uassert.Equal(t, "first", values[1])
	uassert.Equal(t, "second", values[2])
}

func TestList_ForEach(t *testing.T) {
	l := New()

	// Insert test data
	expected := []string{"one", "two", "three"}
	for i, v := range expected {
		l.Insert(i, v)
	}

	// Test ForEach
	var result []string
	l.ForEach(func(key int, value interface{}) bool {
		result = append(result, value.(string))
		return false
	})

	uassert.Equal(t, len(expected), len(result))
	for i, v := range expected {
		uassert.Equal(t, v, result[i])
	}
}

func TestList_InsertBeforeAfter(t *testing.T) {
	l := New()

	// Test with empty list
	_, ok := l.InsertBefore(1, "before")
	uassert.False(t, ok, "should not insert before in empty list")

	_, ok = l.InsertAfter(1, "after")
	uassert.False(t, ok, "should not insert after in empty list")

	// Insert initial element
	l.Insert(10, "ten")

	// Test InsertBefore
	key, ok := l.InsertBefore(10, "nine")
	uassert.True(t, ok, "should insert before")
	uassert.True(t, key < 10, "key should be less than 10")

	// Test InsertAfter
	key, ok = l.InsertAfter(10, "eleven")
	uassert.True(t, ok, "should insert after")
	uassert.True(t, key > 10, "key should be greater than 10")

	// Verify order
	values := l.Values()
	uassert.Equal(t, 3, len(values))
	uassert.Equal(t, "nine", values[0])
	uassert.Equal(t, "ten", values[1])
	uassert.Equal(t, "eleven", values[2])
}

func TestMerge(t *testing.T) {
	// Create first list with even numbers
	l1 := New()
	l1.Insert(2, "two")
	l1.Insert(4, "four")
	l1.Insert(6, "six")

	// Create second list with odd numbers
	l2 := New()
	l2.Insert(1, "one")
	l2.Insert(3, "three")
	l2.Insert(5, "five")

	// Merge lists
	result := Merge(l1, l2)

	// Verify size
	uassert.Equal(t, l1.size+l2.size, result.size, "merged list size should be sum of input sizes")

	// Verify order and values
	expected := []string{"one", "two", "three", "four", "five", "six"}
	values := result.Values()
	uassert.Equal(t, len(expected), len(values), "merged list should have correct number of elements")

	for i, v := range expected {
		uassert.Equal(t, v, values[i], "merged list should maintain correct order")
	}

	// Test merging with empty list
	empty := New()
	result = Merge(empty, l1)
	uassert.Equal(t, l1.size, result.size, "merging with empty list should return same size")

	// Test merging two empty lists
	empty2 := New()
	result = Merge(empty, empty2)
	uassert.Equal(t, 0, result.size, "merging empty lists should return empty list")

	// Test merging with nil
	result = Merge(nil, l1)
	uassert.Equal(t, l1.size, result.size, "merging with nil should return other list")
}

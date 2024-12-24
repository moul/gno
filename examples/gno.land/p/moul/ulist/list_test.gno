package ulist

import (
	"testing"
)

func TestList(t *testing.T) {
	t.Run("large index access", func(t *testing.T) {
		list := New()

		// Add 1000 elements
		for i := 0; i < 1000; i++ {
			list.Append(i)
		}

		// Test random access
		indices := []int{0, 42, 100, 500, 999}
		for _, idx := range indices {
			val, ok := list.Get(idx)
			if !ok || val.(int) != idx {
				t.Errorf("Expected %d at index %d, got %v", idx, idx, val)
			}
		}
	})

	t.Run("prepend performance", func(t *testing.T) {
		list := New()

		// Prepend 1000 elements
		for i := 999; i >= 0; i-- {
			list.Prepend(i)
		}

		// Verify order
		for i := 0; i < 1000; i++ {
			val, ok := list.Get(i)
			if !ok || val.(int) != i {
				t.Errorf("Expected %d at index %d, got %v", i, i, val)
			}
		}
	})

	t.Run("insert operations", func(t *testing.T) {
		list := New()

		// Insert at various positions
		list.Append(0)    // [0]
		list.Append(2)    // [0,2]
		list.Insert(1, 1) // [0,1,2]

		expected := []int{0, 1, 2}
		for i, exp := range expected {
			val, ok := list.Get(i)
			if !ok || val.(int) != exp {
				t.Errorf("Expected %d at index %d, got %v", exp, i, val)
			}
		}
	})

	t.Run("size tracking", func(t *testing.T) {
		list := New()

		if list.Size() != 0 {
			t.Errorf("Expected empty list size 0, got %d", list.Size())
		}

		list.Append(1)
		list.Append(2)
		list.Prepend(0)

		if list.Size() != 3 {
			t.Errorf("Expected list size 3, got %d", list.Size())
		}

		list.Insert(1, 0.5)

		if list.Size() != 4 {
			t.Errorf("Expected list size 4, got %d", list.Size())
		}
	})
}

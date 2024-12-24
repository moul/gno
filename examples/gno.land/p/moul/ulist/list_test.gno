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

func TestDetailedStats(t *testing.T) {
	cases := []struct {
		name string
		fn   func() *List
	}{
		{
			name: "Append 100th item",
			fn: func() *List {
				l := New()
				for i := 0; i < 99; i++ {
					l.Append(i)
				}
				ResetStats()
				l.Append(99)
				return l
			},
		},
		{
			name: "Prepend 100th item",
			fn: func() *List {
				l := New()
				for i := 0; i < 99; i++ {
					l.Append(i)
				}
				ResetStats()
				l.Prepend(99)
				return l
			},
		},
		{
			name: "Insert 100th item in middle",
			fn: func() *List {
				l := New()
				for i := 0; i < 99; i++ {
					l.Append(i)
				}
				ResetStats()
				l.Insert(49, 99)
				return l
			},
		},
		{
			name: "Get middle item from list of 100",
			fn: func() *List {
				l := New()
				for i := 0; i < 100; i++ {
					l.Append(i)
				}
				ResetStats()
				l.Get(49)
				return l
			},
		},
		{
			name: "Append 100 items",
			fn: func() *List {
				l := New()
				ResetStats()
				for i := 0; i < 100; i++ {
					l.Append(i)
				}
				return l
			},
		},
		{
			name: "Prepend 100 items",
			fn: func() *List {
				l := New()
				ResetStats()
				for i := 0; i < 100; i++ {
					l.Prepend(i)
				}
				return l
			},
		},
		{
			name: "Sequential read of 100 items",
			fn: func() *List {
				l := New()
				for i := 0; i < 100; i++ {
					l.Append(i)
				}
				ResetStats()
				for i := 0; i < 100; i++ {
					l.Get(i)
				}
				return l
			},
		},
		{
			name: "Random access of 20 items in list of 100",
			fn: func() *List {
				l := New()
				for i := 0; i < 100; i++ {
					l.Append(i)
				}
				ResetStats()
				// Access some scattered indices
				indices := []int{0, 5, 10, 20, 30, 40, 50, 60, 70, 80, 90, 95, 99}
				for _, idx := range indices {
					l.Get(idx)
				}
				return l
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ResetStats()
			tc.fn()
			stats := GetStats()
			t.Logf("Operation: %s", tc.name)
			t.Logf("- Reads: %d", stats.reads)
			t.Logf("- Writes: %d", stats.writes)
			t.Logf("- Balances: %d", stats.balances)
			t.Logf("- Rotations: %d", stats.rotations)
			t.Logf("- Height calculations: %d", stats.heightCalcs)
			t.Logf("- Size calculations: %d", stats.sizeCalcs)
			t.Logf("---")
		})
	}
}

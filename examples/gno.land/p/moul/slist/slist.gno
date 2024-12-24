// Package slist implements a skip list data structure with deterministic level selection.
//
// This implementation avoids using math/rand for performance reasons, instead using
// a size-based approach for level determination. The structure is based on the list size,
// making it more efficient in the Gno environment.
//
// Performance Characteristics (compared to AVL trees):
// - Search: O(log n) - Similar to AVL trees
// - Insert: O(log n) - Similar to AVL trees, but with different trade-offs:
//   - Skip lists require more memory due to multiple levels of pointers
//   - In Gno specifically, skip lists may be more expensive for insertions
//     as they modify slices (which creates new copies)
//   - However, skip lists don't require complex rotation operations
//
// - Delete: O(log n) - Similar considerations as Insert
//
// Implementation Notes:
//   - Uses a leveled structure where each level acts as a "fast lane" for searching
//   - The number of levels grows logarithmically with the size of the list
//   - Level selection is based on size thresholds (powers of 2) rather than random selection
//   - In Gno, operations that modify slices are expensive as they create new copies
//     This implementation accepts this trade-off for simplicity
//
// Features:
// - Standard operations: Insert, Search, Delete
// - List-like operations: Front, Back, Len, Clear
// - Sequential operations: Append, Prepend
// - Relative insertions: InsertBefore, InsertAfter
// - Iteration: ForEach, Values, Keys
//
// Use Cases:
// - When a simpler implementation is preferred over AVL trees
// - When the data structure needs to maintain sorted order
// - When rotation-free operations are desired
// - When list-like sequential operations are needed alongside ordered access
package slist

import (
	"strings"

	"gno.land/p/demo/ufmt"
)

const (
	maxLevel = 16 // Maximum levels in the skip list
)

// node represents a single node in the skip list
type node struct {
	key   int
	value interface{}
	next  []*node // Slices for level pointers
}

// SkipList represents the entire skip list
type SkipList struct {
	head  *node
	level int
	size  int
}

// NewSkipList creates and initializes a new deterministic SkipList
func NewSkipList() *SkipList {
	head := &node{
		next: make([]*node, maxLevel),
	}
	return &SkipList{
		head:  head,
		level: 1,
		size:  0,
	}
}

// calculateLevel deterministically calculates the level of a new node
func (sl *SkipList) calculateLevel() int {
	level := 1
	threshold := 2
	for sl.size >= threshold && level < maxLevel {
		level++
		threshold <<= 1 // Double the threshold at each level (powers of 2)
	}
	return level
}

// Insert adds a new key-value pair to the skip list
func (sl *SkipList) Insert(key int, value interface{}) {
	update := make([]*node, maxLevel)
	current := sl.head

	// Find the position for the new key
	for i := sl.level - 1; i >= 0; i-- {
		for current.next[i] != nil && current.next[i].key < key {
			current = current.next[i]
		}
		update[i] = current
	}

	// Check if the key already exists
	if current.next[0] != nil && current.next[0].key == key {
		current.next[0].value = value
		return
	}

	// Determine the level of the new node
	level := sl.calculateLevel()
	if level > sl.level {
		for i := sl.level; i < level; i++ {
			update[i] = sl.head
		}
		sl.level = level
	}

	// Create a new node
	newNode := &node{
		key:   key,
		value: value,
		next:  make([]*node, level),
	}

	// Update pointers
	for i := 0; i < level; i++ {
		newNode.next[i] = update[i].next[i]
		update[i].next[i] = newNode
	}

	sl.size++
}

// Search finds a value by its key
func (sl *SkipList) Search(key int) (interface{}, bool) {
	current := sl.head
	for i := sl.level - 1; i >= 0; i-- {
		for current.next[i] != nil && current.next[i].key < key {
			current = current.next[i]
		}
	}
	current = current.next[0]
	if current != nil && current.key == key {
		return current.value, true
	}
	return nil, false
}

// Delete removes a key-value pair from the skip list
func (sl *SkipList) Delete(key int) {
	update := make([]*node, maxLevel)
	current := sl.head

	// Find the position of the node to delete
	for i := sl.level - 1; i >= 0; i-- {
		for current.next[i] != nil && current.next[i].key < key {
			current = current.next[i]
		}
		update[i] = current
	}

	target := current.next[0]
	if target == nil || target.key != key {
		return // Key not found
	}

	// Update pointers
	for i := 0; i < sl.level; i++ {
		if update[i].next[i] != target {
			break
		}
		update[i].next[i] = target.next[i]
	}

	// Reduce the level if necessary
	for sl.level > 1 && sl.head.next[sl.level-1] == nil {
		sl.level--
	}

	sl.size--
}

// Print displays the contents of the skip list
func (sl *SkipList) Sprint() string {
	var sb strings.Builder
	for i := sl.level - 1; i >= 0; i-- {
		current := sl.head.next[i]
		sb.WriteString(ufmt.Sprintf("Level %d: ", i+1))
		for current != nil {
			sb.WriteString(ufmt.Sprintf("%d ", current.key))
			current = current.next[i]
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

// Front returns the first element of the list.
func (sl *SkipList) Front() (key int, value interface{}, exists bool) {
	if sl.size == 0 {
		return 0, nil, false
	}
	current := sl.head.next[0]
	return current.key, current.value, true
}

// Back returns the last element of the list.
func (sl *SkipList) Back() (key int, value interface{}, exists bool) {
	if sl.size == 0 {
		return 0, nil, false
	}
	current := sl.head
	for i := sl.level - 1; i >= 0; i-- {
		for current.next[i] != nil {
			current = current.next[i]
		}
	}
	return current.key, current.value, true
}

// Len returns the number of elements in the list.
func (sl *SkipList) Len() int {
	return sl.size
}

// Clear removes all elements from the list.
func (sl *SkipList) Clear() {
	sl.head.next = make([]*node, maxLevel)
	sl.level = 1
	sl.size = 0
}

// Append adds a new element to the end of the list.
// It automatically assigns a key that is greater than the current maximum key.
func (sl *SkipList) Append(value interface{}) int {
	var newKey int
	if sl.size == 0 {
		newKey = 0
	} else {
		lastKey, _, _ := sl.Back()
		newKey = lastKey + 1
	}
	sl.Insert(newKey, value)
	return newKey
}

// Prepend adds a new element to the beginning of the list.
// It automatically assigns a key that is less than the current minimum key.
func (sl *SkipList) Prepend(value interface{}) int {
	var newKey int
	if sl.size == 0 {
		newKey = 0
	} else {
		firstKey, _, _ := sl.Front()
		newKey = firstKey - 1
	}
	sl.Insert(newKey, value)
	return newKey
}

// ForEach iterates through all elements in the list in order.
func (sl *SkipList) ForEach(fn func(key int, value interface{}) bool) {
	current := sl.head.next[0]
	for current != nil {
		if fn(current.key, current.value) {
			break
		}
		current = current.next[0]
	}
}

// Values returns a slice of all values in the list, in order.
func (sl *SkipList) Values() []interface{} {
	values := make([]interface{}, 0, sl.size)
	sl.ForEach(func(key int, value interface{}) bool {
		values = append(values, value)
		return false
	})
	return values
}

// Keys returns a slice of all keys in the list, in order.
func (sl *SkipList) Keys() []int {
	keys := make([]int, 0, sl.size)
	sl.ForEach(func(key int, value interface{}) bool {
		keys = append(keys, key)
		return false
	})
	return keys
}

// InsertAfter inserts a new element after the element with the given key.
// Returns false if the key is not found.
func (sl *SkipList) InsertAfter(key int, value interface{}) (newKey int, ok bool) {
	current := sl.head
	for i := sl.level - 1; i >= 0; i-- {
		for current.next[i] != nil && current.next[i].key < key {
			current = current.next[i]
		}
	}
	current = current.next[0]
	if current == nil || current.key != key {
		return 0, false
	}

	// Calculate the new key
	var nextKey int
	if current.next[0] == nil {
		nextKey = key + 1
	} else {
		nextKey = (key + current.next[0].key) / 2
		if nextKey == key { // If there's no space between keys
			nextKey = key + 1
		}
	}

	sl.Insert(nextKey, value)
	return nextKey, true
}

// InsertBefore inserts a new element before the element with the given key.
// Returns false if the key is not found.
func (sl *SkipList) InsertBefore(key int, value interface{}) (newKey int, ok bool) {
	current := sl.head
	var prevNode *node

	// Find the node before the target key
	for i := sl.level - 1; i >= 0; i-- {
		for current.next[i] != nil && current.next[i].key < key {
			current = current.next[i]
		}
		if current.next[i] != nil && current.next[i].key == key {
			prevNode = current
		}
	}

	// Check if key exists
	if current.next[0] == nil || current.next[0].key != key {
		return 0, false
	}

	// Calculate the new key
	var prevKey int
	if prevNode == sl.head {
		prevKey = key - 1
	} else {
		prevKey = (prevNode.key + key) / 2
		if prevKey == key { // If there's no space between keys
			prevKey = key - 1
		}
	}

	sl.Insert(prevKey, value)
	return prevKey, true
}

// Has returns true if the key exists in the list.
func (sl *SkipList) Has(key int) bool {
	_, found := sl.Search(key)
	return found
}

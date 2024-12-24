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

// List represents a deterministic skip list.
type List struct {
	head     *node
	level    int
	size     int
	maxLevel int
}

// New creates and initializes a new List with default maximum level of 16.
func New() *List {
	return NewWithMaxLevel(16)
}

// NewWithMaxLevel creates and initializes a new List with specified maximum level.
func NewWithMaxLevel(maxLevel int) *List {
	head := &node{
		next: make([]*node, maxLevel),
	}
	return &List{
		head:     head,
		level:    1,
		size:     0,
		maxLevel: maxLevel,
	}
}

// node represents a single node in the skip list.
type node struct {
	key   int
	value interface{}
	next  []*node // Slices for level pointers
}

// Core Operations

// Insert adds a new key-value pair to the list.
func (l *List) Insert(key int, value interface{}) {
	update := make([]*node, l.maxLevel)
	current := l.head

	// Find the position for the new key
	for i := l.level - 1; i >= 0; i-- {
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
	level := l.calculateLevel()
	if level > l.level {
		for i := l.level; i < level; i++ {
			update[i] = l.head
		}
		l.level = level
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

	l.size++
}

// Search finds a value by its key.
func (l *List) Search(key int) (interface{}, bool) {
	current := l.head
	for i := l.level - 1; i >= 0; i-- {
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

// Delete removes a key-value pair from the list.
func (l *List) Delete(key int) {
	update := make([]*node, l.maxLevel)
	current := l.head

	// Find the position of the node to delete
	for i := l.level - 1; i >= 0; i-- {
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
	for i := 0; i < l.level; i++ {
		if update[i].next[i] != target {
			break
		}
		update[i].next[i] = target.next[i]
	}

	// Reduce the level if necessary
	for l.level > 1 && l.head.next[l.level-1] == nil {
		l.level--
	}

	l.size--
}

// List-like Operations

// Front returns the first element of the list.
func (l *List) Front() (key int, value interface{}, exists bool) {
	if l.size == 0 {
		return 0, nil, false
	}
	current := l.head.next[0]
	return current.key, current.value, true
}

// Back returns the last element of the list.
func (l *List) Back() (key int, value interface{}, exists bool) {
	if l.size == 0 {
		return 0, nil, false
	}
	current := l.head
	for i := l.level - 1; i >= 0; i-- {
		for current.next[i] != nil {
			current = current.next[i]
		}
	}
	return current.key, current.value, true
}

// Len returns the number of elements in the list.
func (l *List) Len() int {
	return l.size
}

// Clear removes all elements from the list.
func (l *List) Clear() {
	l.head.next = make([]*node, l.maxLevel)
	l.level = 1
	l.size = 0
}

// Helper Methods

// calculateLevel deterministically calculates the level of a new node.
func (l *List) calculateLevel() int {
	level := 1
	threshold := 2
	for l.size >= threshold && level < l.maxLevel {
		level++
		threshold <<= 1 // Double the threshold at each level (powers of 2)
	}
	return level
}

// Has returns true if the key exists in the list.
func (l *List) Has(key int) bool {
	_, found := l.Search(key)
	return found
}

// Sprint displays the contents of the list.
func (l *List) Sprint() string {
	var sb strings.Builder
	for i := l.level - 1; i >= 0; i-- {
		current := l.head.next[i]
		sb.WriteString(ufmt.Sprintf("Level %d: ", i+1))
		for current != nil {
			sb.WriteString(ufmt.Sprintf("%d ", current.key))
			current = current.next[i]
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

// Extended Operations

// Append adds a new element to the end of the list.
func (l *List) Append(value interface{}) int {
	var newKey int
	if l.size == 0 {
		newKey = 0
	} else {
		lastKey, _, _ := l.Back()
		newKey = lastKey + 1
	}
	l.Insert(newKey, value)
	return newKey
}

// Prepend adds a new element to the beginning of the list.
func (l *List) Prepend(value interface{}) int {
	var newKey int
	if l.size == 0 {
		newKey = 0
	} else {
		firstKey, _, _ := l.Front()
		newKey = firstKey - 1
	}
	l.Insert(newKey, value)
	return newKey
}

// ForEach iterates through all elements in the list in order.
func (l *List) ForEach(fn func(key int, value interface{}) bool) {
	current := l.head.next[0]
	for current != nil {
		if fn(current.key, current.value) {
			break
		}
		current = current.next[0]
	}
}

// Values returns a slice of all values in the list, in order.
func (l *List) Values() []interface{} {
	values := make([]interface{}, 0, l.size)
	l.ForEach(func(key int, value interface{}) bool {
		values = append(values, value)
		return false
	})
	return values
}

// Keys returns a slice of all keys in the list, in order.
func (l *List) Keys() []int {
	keys := make([]int, 0, l.size)
	l.ForEach(func(key int, value interface{}) bool {
		keys = append(keys, key)
		return false
	})
	return keys
}

// InsertAfter inserts a new element after the element with the given key.
func (l *List) InsertAfter(key int, value interface{}) (newKey int, ok bool) {
	current := l.head
	for i := l.level - 1; i >= 0; i-- {
		for current.next[i] != nil && current.next[i].key < key {
			current = current.next[i]
		}
	}
	current = current.next[0]
	if current == nil || current.key != key {
		return 0, false
	}

	var nextKey int
	if current.next[0] == nil {
		nextKey = key + 1
	} else {
		nextKey = (key + current.next[0].key) / 2
		if nextKey == key {
			nextKey = key + 1
		}
	}

	l.Insert(nextKey, value)
	return nextKey, true
}

// InsertBefore inserts a new element before the element with the given key.
func (l *List) InsertBefore(key int, value interface{}) (newKey int, ok bool) {
	current := l.head
	var prevNode *node

	for i := l.level - 1; i >= 0; i-- {
		for current.next[i] != nil && current.next[i].key < key {
			current = current.next[i]
		}
		if current.next[i] != nil && current.next[i].key == key {
			prevNode = current
		}
	}

	if current.next[0] == nil || current.next[0].key != key {
		return 0, false
	}

	var prevKey int
	if prevNode == l.head {
		prevKey = key - 1
	} else {
		prevKey = (prevNode.key + key) / 2
		if prevKey == key {
			prevKey = key - 1
		}
	}

	l.Insert(prevKey, value)
	return prevKey, true
}

// Merge efficiently combines two lists into a new one.
// Both input lists remain unchanged.
func Merge(l1, l2 *List) *List {
	if l1 == nil || l1.size == 0 {
		return l2
	}
	if l2 == nil || l2.size == 0 {
		return l1
	}

	// Create a new list with max level from either list
	maxLevel := l1.maxLevel
	if l2.maxLevel > maxLevel {
		maxLevel = l2.maxLevel
	}
	result := NewWithMaxLevel(maxLevel)

	// Traverse both lists simultaneously
	curr1 := l1.head.next[0]
	curr2 := l2.head.next[0]

	// Pre-allocate update array for better performance
	update := make([]*node, maxLevel)
	current := result.head

	// Merge nodes from both lists
	for curr1 != nil && curr2 != nil {
		var nextNode *node
		if curr1.key < curr2.key {
			nextNode = &node{
				key:   curr1.key,
				value: curr1.value,
				next:  make([]*node, result.calculateLevel()),
			}
			curr1 = curr1.next[0]
		} else if curr2.key < curr1.key {
			nextNode = &node{
				key:   curr2.key,
				value: curr2.value,
				next:  make([]*node, result.calculateLevel()),
			}
			curr2 = curr2.next[0]
		} else {
			// In case of equal keys, take from l1 (arbitrary choice)
			nextNode = &node{
				key:   curr1.key,
				value: curr1.value,
				next:  make([]*node, result.calculateLevel()),
			}
			curr1 = curr1.next[0]
			curr2 = curr2.next[0]
		}

		// Update all levels
		for i := 0; i < len(nextNode.next); i++ {
			nextNode.next[i] = update[i]
			update[i] = nextNode
		}

		if len(nextNode.next) > result.level {
			result.level = len(nextNode.next)
		}

		result.size++
		current = nextNode
	}

	// Add remaining nodes from l1
	for curr1 != nil {
		nextNode := &node{
			key:   curr1.key,
			value: curr1.value,
			next:  make([]*node, result.calculateLevel()),
		}

		// Update all levels
		for i := 0; i < len(nextNode.next); i++ {
			nextNode.next[i] = update[i]
			update[i] = nextNode
		}

		if len(nextNode.next) > result.level {
			result.level = len(nextNode.next)
		}

		result.size++
		current = nextNode
		curr1 = curr1.next[0]
	}

	// Add remaining nodes from l2
	for curr2 != nil {
		nextNode := &node{
			key:   curr2.key,
			value: curr2.value,
			next:  make([]*node, result.calculateLevel()),
		}

		// Update all levels
		for i := 0; i < len(nextNode.next); i++ {
			nextNode.next[i] = update[i]
			update[i] = nextNode
		}

		if len(nextNode.next) > result.level {
			result.level = len(nextNode.next)
		}

		result.size++
		current = nextNode
		curr2 = curr2.next[0]
	}

	// Fix head pointers
	for i := 0; i < result.level; i++ {
		result.head.next[i] = update[i]
	}

	return result
}

package index

import (
	"gno.land/p/demo/avl"
)

// KeyExtractor is a function that extracts an index key from a value
type KeyExtractor func(value interface{}) string

// IndexedTree wraps an AVL tree and maintains secondary indexes
type IndexedTree struct {
	primary *avl.Tree
	indexes map[string]*Index
}

// Index represents a secondary index
type Index struct {
	tree    *avl.Tree
	extract KeyExtractor
}

// NewIndexedTree creates a new indexed tree wrapper
func NewIndexedTree() *IndexedTree {
	return &IndexedTree{
		primary: avl.NewTree(),
		indexes: make(map[string]*Index),
	}
}

// AddIndex creates a new secondary index with the given name and key extractor
func (it *IndexedTree) AddIndex(name string, extract KeyExtractor) {
	it.indexes[name] = &Index{
		tree:    avl.NewTree(),
		extract: extract,
	}
}

// Set adds or updates a value in the primary tree and all indexes
func (it *IndexedTree) Set(key string, value interface{}) bool {
	// First, if this is an update, we need to remove old index entries
	if oldValue, exists := it.primary.Get(key); exists {
		it.removeFromIndexes(key, oldValue)
	}

	// Update primary tree
	updated := it.primary.Set(key, value)

	// Update all indexes
	it.addToIndexes(key, value)

	return updated
}

// Remove removes a value from the primary tree and all indexes
func (it *IndexedTree) Remove(key string) (interface{}, bool) {
	// Get the value first so we can remove it from indexes
	value, exists := it.primary.Get(key)
	if !exists {
		return nil, false
	}

	// Remove from indexes first
	it.removeFromIndexes(key, value)

	// Remove from primary tree
	return it.primary.Remove(key)
}

// GetByIndex retrieves values from a secondary index
func (it *IndexedTree) GetByIndex(indexName string, indexKey string) []interface{} {
	index, exists := it.indexes[indexName]
	if !exists {
		return nil
	}

	var results []interface{}
	index.tree.Iterate(indexKey, indexKey, func(key string, value interface{}) bool {
		primaryKey := value.(string) // The index stores primary keys as values
		if primaryValue, exists := it.primary.Get(primaryKey); exists {
			results = append(results, primaryValue)
		}
		return false // Continue iteration
	})

	return results
}

// Internal helper methods

func (it *IndexedTree) addToIndexes(primaryKey string, value interface{}) {
	for _, index := range it.indexes {
		indexKey := index.extract(value)
		index.tree.Set(indexKey, primaryKey)
	}
}

func (it *IndexedTree) removeFromIndexes(primaryKey string, value interface{}) {
	for _, index := range it.indexes {
		indexKey := index.extract(value)
		index.tree.Remove(indexKey)
	}
}

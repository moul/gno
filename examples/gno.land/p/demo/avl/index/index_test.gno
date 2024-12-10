package index

import (
	"strconv"
	"testing"
)

type Person struct {
	ID   string
	Name string
	Age  int
}

func TestIndexedTree(t *testing.T) {
	// Create a new indexed tree
	it := NewIndexedTree()

	// Add secondary indexes
	it.AddIndex("name", func(value interface{}) string {
		return value.(*Person).Name
	})

	it.AddIndex("age", func(value interface{}) string {
		return strconv.Itoa(value.(*Person).Age)
	})

	// Add some test data
	person1 := &Person{ID: "1", Name: "Alice", Age: 30}
	person2 := &Person{ID: "2", Name: "Bob", Age: 30}

	it.Set(person1.ID, person1)
	it.Set(person2.ID, person2)

	// Query by name index
	aliceResults := it.GetByIndex("name", "Alice")
	println(aliceResults)
	if len(aliceResults) != 1 {
		t.Error("Expected 1 result for Alice")
	}

	// Query by age index
	age30Results := it.GetByIndex("age", "30")
	if len(age30Results) != 2 {
		t.Error("Expected 2 results for age 30")
	}
}

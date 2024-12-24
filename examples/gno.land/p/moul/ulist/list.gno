package ulist

// Node represents a node in the balanced list tree
type Node struct {
	value   interface{}
	size_   int  // total number of nodes in this subtree
	height_ int8 // height of this subtree; using int8 limits the list to ~2^128 elements (2^height)
	left    *Node
	right   *Node
}

// List represents a balanced tree-based list
type List struct {
	root *Node
}

// New creates a new empty list
func New() *List {
	return &List{}
}

// Size returns the number of elements in the list
func (l *List) Size() int {
	if l.root == nil {
		return 0
	}
	return l.root.size()
}

// Get returns the element at the specified index
func (l *List) Get(index int) (interface{}, bool) {
	if index < 0 || index >= l.Size() {
		return nil, false
	}
	return l.root.getByIndex(index)
}

// Prepend adds an element to the beginning of the list
func (l *List) Prepend(value interface{}) {
	newNode := &Node{
		value:   value,
		size_:   1,
		height_: 1,
	}
	if l.root == nil {
		l.root = newNode
		return
	}
	l.root = insertAt(l.root, 0, newNode)
}

// Append adds an element to the end of the list
func (l *List) Append(value interface{}) {
	newNode := &Node{
		value:   value,
		size_:   1,
		height_: 1,
	}
	if l.root == nil {
		l.root = newNode
		return
	}
	l.root = insertAt(l.root, l.Size(), newNode)
}

// Insert adds an element at the specified index
func (l *List) Insert(index int, value interface{}) bool {
	if index < 0 || index > l.Size() {
		return false
	}
	newNode := &Node{
		value:   value,
		size_:   1,
		height_: 1,
	}
	l.root = insertAt(l.root, index, newNode)
	return true
}

// Internal helper functions

func (n *Node) getByIndex(index int) (interface{}, bool) {
	leftSize := n.left.size()

	if index == leftSize {
		return n.value, true
	}
	if index < leftSize {
		return n.left.getByIndex(index)
	}
	return n.right.getByIndex(index - leftSize - 1)
}

func insertAt(node *Node, index int, newNode *Node) *Node {
	if node == nil {
		return newNode
	}

	leftSize := node.left.size()
	if index <= leftSize {
		node.left = insertAt(node.left, index, newNode)
	} else {
		node.right = insertAt(node.right, index-leftSize-1, newNode)
	}

	node.size_ = node.left.size() + node.right.size() + 1
	node.height_ = maxInt8(node.left.height(), node.right.height()) + 1

	return balance(node)
}

// Balance functions

func balance(node *Node) *Node {
	if node == nil {
		return nil
	}

	balanceFactor := node.getBalance()

	// Left heavy
	if balanceFactor > 1 {
		if node.left.getBalance() < 0 {
			node.left = rotateLeft(node.left)
		}
		return rotateRight(node)
	}

	// Right heavy
	if balanceFactor < -1 {
		if node.right.getBalance() > 0 {
			node.right = rotateRight(node.right)
		}
		return rotateLeft(node)
	}

	return node
}

func rotateLeft(node *Node) *Node {
	rightNode := node.right
	node.right = rightNode.left
	rightNode.left = node

	node.size_ = node.left.size() + node.right.size() + 1
	node.height_ = maxInt8(node.left.height(), node.right.height()) + 1

	rightNode.size_ = rightNode.left.size() + rightNode.right.size() + 1
	rightNode.height_ = maxInt8(rightNode.left.height(), rightNode.right.height()) + 1

	return rightNode
}

func rotateRight(node *Node) *Node {
	leftNode := node.left
	node.left = leftNode.right
	leftNode.right = node

	node.size_ = node.left.size() + node.right.size() + 1
	node.height_ = maxInt8(node.left.height(), node.right.height()) + 1

	leftNode.size_ = leftNode.left.size() + leftNode.right.size() + 1
	leftNode.height_ = maxInt8(leftNode.left.height(), leftNode.right.height()) + 1

	return leftNode
}

// Convert helper functions to methods
func (n *Node) size() int {
	if n == nil {
		return 0
	}
	return n.size_
}

func (n *Node) height() int8 {
	if n == nil {
		return 0
	}
	return n.height_
}

func (n *Node) getBalance() int {
	if n == nil {
		return 0
	}
	return int(n.left.height()) - int(n.right.height())
}

func maxInt8(a, b int8) int8 {
	if a > b {
		return a
	}
	return b
}

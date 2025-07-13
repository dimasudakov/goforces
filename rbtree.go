package main

import (
	"cmp"
	"fmt"
)

// original implementation: https://github.com/emirpasic/gods/tree/master

type color bool

const (
	black, red color = true, false
)

type RBComparator[T any] func(x, y T) int

// RBTree holds elements of the red-black tree
type RBTree[K comparable, V any] struct {
	Root       *RBNode[K, V]
	size       int
	Comparator RBComparator[K]
}

// RBNode is a single element within the tree
type RBNode[K comparable, V any] struct {
	Key    K
	Value  V
	color  color
	Left   *RBNode[K, V]
	Right  *RBNode[K, V]
	Parent *RBNode[K, V]
}

// NewRBT instantiates a red-black tree with the built-in comparator for K
func NewRBT[K cmp.Ordered, V any]() *RBTree[K, V] {
	return &RBTree[K, V]{Comparator: cmp.Compare[K]}
}

// NewRBTWithComp instantiates a red-black tree with the custom comparator.
func NewRBTWithComp[K comparable, V any](comparator RBComparator[K]) *RBTree[K, V] {
	return &RBTree[K, V]{Comparator: comparator}
}

// Put inserts node into the tree.
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (tree *RBTree[K, V]) Put(key K, value V) {
	var insertedNode *RBNode[K, V]
	if tree.Root == nil {
		// Assert key is of comparator's type for initial tree
		tree.Comparator(key, key)
		tree.Root = &RBNode[K, V]{Key: key, Value: value, color: red}
		insertedNode = tree.Root
	} else {
		node := tree.Root
		loop := true
		for loop {
			compare := tree.Comparator(key, node.Key)
			switch {
			case compare == 0:
				node.Key = key
				node.Value = value
				return
			case compare < 0:
				if node.Left == nil {
					node.Left = &RBNode[K, V]{Key: key, Value: value, color: red}
					insertedNode = node.Left
					loop = false
				} else {
					node = node.Left
				}
			case compare > 0:
				if node.Right == nil {
					node.Right = &RBNode[K, V]{Key: key, Value: value, color: red}
					insertedNode = node.Right
					loop = false
				} else {
					node = node.Right
				}
			}
		}
		insertedNode.Parent = node
	}
	tree.insertCase1(insertedNode)
	tree.size++
}

// Get searches the node in the tree by key and returns its value or nil if key is not found in tree.
// Second return parameter is true if key was found, otherwise false.
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (tree *RBTree[K, V]) Get(key K) (value V, found bool) {
	node := tree.lookup(key)
	if node != nil {
		return node.Value, true
	}
	return value, false
}

// GetNode searches the node in the tree by key and returns its node or nil if key is not found in tree.
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (tree *RBTree[K, V]) GetNode(key K) *RBNode[K, V] {
	return tree.lookup(key)
}

// Remove remove the node from the tree by key.
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (tree *RBTree[K, V]) Remove(key K) {
	var child *RBNode[K, V]
	node := tree.lookup(key)
	if node == nil {
		return
	}
	if node.Left != nil && node.Right != nil {
		pred := node.Left.maximumNode()
		node.Key = pred.Key
		node.Value = pred.Value
		node = pred
	}
	if node.Left == nil || node.Right == nil {
		if node.Right == nil {
			child = node.Left
		} else {
			child = node.Right
		}
		if node.color == black {
			node.color = nodeColor(child)
			tree.deleteCase1(node)
		}
		tree.replaceNode(node, child)
		if node.Parent == nil && child != nil {
			child.color = black
		}
	}
	tree.size--
}

// Empty returns true if tree does not contain any nodes
func (tree *RBTree[K, V]) Empty() bool {
	return tree.size == 0
}

// Size returns number of nodes in the tree.
func (tree *RBTree[K, V]) Size() int {
	return tree.size
}

// Size returns the number of elements stored in the subtree.
// Computed dynamically on each call, i.e. the subtree is traversed to count the number of the nodes.
func (node *RBNode[K, V]) Size() int {
	if node == nil {
		return 0
	}
	size := 1
	if node.Left != nil {
		size += node.Left.Size()
	}
	if node.Right != nil {
		size += node.Right.Size()
	}
	return size
}

// Keys returns all keys in-order
func (tree *RBTree[K, V]) Keys() []K {
	keys := make([]K, tree.size)
	it := tree.Iterator()
	for i := 0; it.Next(); i++ {
		keys[i] = it.Key()
	}
	return keys
}

// Values returns all values in-order based on the key.
func (tree *RBTree[K, V]) Values() []V {
	values := make([]V, tree.size)
	it := tree.Iterator()
	for i := 0; it.Next(); i++ {
		values[i] = it.Value()
	}
	return values
}

// Left returns the left-most (min) node or nil if tree is empty.
func (tree *RBTree[K, V]) Left() *RBNode[K, V] {
	var parent *RBNode[K, V]
	current := tree.Root
	for current != nil {
		parent = current
		current = current.Left
	}
	return parent
}

// Right returns the right-most (max) node or nil if tree is empty.
func (tree *RBTree[K, V]) Right() *RBNode[K, V] {
	var parent *RBNode[K, V]
	current := tree.Root
	for current != nil {
		parent = current
		current = current.Right
	}
	return parent
}

// Floor Finds floor node of the input key, return the floor node or nil if no floor is found.
// Second return parameter is true if floor was found, otherwise false.
//
// Floor node is defined as the largest node that is smaller than or equal to the given node.
// A floor node may not be found, either because the tree is empty, or because
// all nodes in the tree are larger than the given node.
//
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (tree *RBTree[K, V]) Floor(key K) (floor *RBNode[K, V], found bool) {
	found = false
	node := tree.Root
	for node != nil {
		compare := tree.Comparator(key, node.Key)
		switch {
		case compare == 0:
			return node, true
		case compare < 0:
			node = node.Left
		case compare > 0:
			floor, found = node, true
			node = node.Right
		}
	}
	if found {
		return floor, true
	}
	return nil, false
}

// Ceiling finds ceiling node of the input key, return the ceiling node or nil if no ceiling is found.
// Second return parameter is true if ceiling was found, otherwise false.
//
// Ceiling node is defined as the smallest node that is larger than or equal to the given node.
// A ceiling node may not be found, either because the tree is empty, or because
// all nodes in the tree are smaller than the given node.
//
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (tree *RBTree[K, V]) Ceiling(key K) (ceiling *RBNode[K, V], found bool) {
	found = false
	node := tree.Root
	for node != nil {
		compare := tree.Comparator(key, node.Key)
		switch {
		case compare == 0:
			return node, true
		case compare < 0:
			ceiling, found = node, true
			node = node.Left
		case compare > 0:
			node = node.Right
		}
	}
	if found {
		return ceiling, true
	}
	return nil, false
}

// Clear removes all nodes from the tree.
func (tree *RBTree[K, V]) Clear() {
	tree.Root = nil
	tree.size = 0
}

// String returns a string representation of container
func (tree *RBTree[K, V]) String() string {
	str := "RedBlackTree\n"
	if !tree.Empty() {
		output(tree.Root, "", true, &str)
	}
	return str
}

func (node *RBNode[K, V]) String() string {
	return fmt.Sprintf("%v", node.Key)
}

func output[K comparable, V any](node *RBNode[K, V], prefix string, isTail bool, str *string) {
	if node.Right != nil {
		newPrefix := prefix
		if isTail {
			newPrefix += "│   "
		} else {
			newPrefix += "    "
		}
		output(node.Right, newPrefix, false, str)
	}
	*str += prefix
	if isTail {
		*str += "└── "
	} else {
		*str += "┌── "
	}
	*str += node.String() + "\n"
	if node.Left != nil {
		newPrefix := prefix
		if isTail {
			newPrefix += "    "
		} else {
			newPrefix += "│   "
		}
		output(node.Left, newPrefix, true, str)
	}
}

func (tree *RBTree[K, V]) lookup(key K) *RBNode[K, V] {
	node := tree.Root
	for node != nil {
		compare := tree.Comparator(key, node.Key)
		switch {
		case compare == 0:
			return node
		case compare < 0:
			node = node.Left
		case compare > 0:
			node = node.Right
		}
	}
	return nil
}

func (node *RBNode[K, V]) grandparent() *RBNode[K, V] {
	if node != nil && node.Parent != nil {
		return node.Parent.Parent
	}
	return nil
}

func (node *RBNode[K, V]) uncle() *RBNode[K, V] {
	if node == nil || node.Parent == nil || node.Parent.Parent == nil {
		return nil
	}
	return node.Parent.sibling()
}

func (node *RBNode[K, V]) sibling() *RBNode[K, V] {
	if node == nil || node.Parent == nil {
		return nil
	}
	if node == node.Parent.Left {
		return node.Parent.Right
	}
	return node.Parent.Left
}

func (tree *RBTree[K, V]) rotateLeft(node *RBNode[K, V]) {
	right := node.Right
	tree.replaceNode(node, right)
	node.Right = right.Left
	if right.Left != nil {
		right.Left.Parent = node
	}
	right.Left = node
	node.Parent = right
}

func (tree *RBTree[K, V]) rotateRight(node *RBNode[K, V]) {
	left := node.Left
	tree.replaceNode(node, left)
	node.Left = left.Right
	if left.Right != nil {
		left.Right.Parent = node
	}
	left.Right = node
	node.Parent = left
}

func (tree *RBTree[K, V]) replaceNode(old *RBNode[K, V], new *RBNode[K, V]) {
	if old.Parent == nil {
		tree.Root = new
	} else {
		if old == old.Parent.Left {
			old.Parent.Left = new
		} else {
			old.Parent.Right = new
		}
	}
	if new != nil {
		new.Parent = old.Parent
	}
}

func (tree *RBTree[K, V]) insertCase1(node *RBNode[K, V]) {
	if node.Parent == nil {
		node.color = black
	} else {
		tree.insertCase2(node)
	}
}

func (tree *RBTree[K, V]) insertCase2(node *RBNode[K, V]) {
	if nodeColor(node.Parent) == black {
		return
	}
	tree.insertCase3(node)
}

func (tree *RBTree[K, V]) insertCase3(node *RBNode[K, V]) {
	uncle := node.uncle()
	if nodeColor(uncle) == red {
		node.Parent.color = black
		uncle.color = black
		node.grandparent().color = red
		tree.insertCase1(node.grandparent())
	} else {
		tree.insertCase4(node)
	}
}

func (tree *RBTree[K, V]) insertCase4(node *RBNode[K, V]) {
	grandparent := node.grandparent()
	if node == node.Parent.Right && node.Parent == grandparent.Left {
		tree.rotateLeft(node.Parent)
		node = node.Left
	} else if node == node.Parent.Left && node.Parent == grandparent.Right {
		tree.rotateRight(node.Parent)
		node = node.Right
	}
	tree.insertCase5(node)
}

func (tree *RBTree[K, V]) insertCase5(node *RBNode[K, V]) {
	node.Parent.color = black
	grandparent := node.grandparent()
	grandparent.color = red
	if node == node.Parent.Left && node.Parent == grandparent.Left {
		tree.rotateRight(grandparent)
	} else if node == node.Parent.Right && node.Parent == grandparent.Right {
		tree.rotateLeft(grandparent)
	}
}

func (node *RBNode[K, V]) maximumNode() *RBNode[K, V] {
	if node == nil {
		return nil
	}
	for node.Right != nil {
		node = node.Right
	}
	return node
}

func (tree *RBTree[K, V]) deleteCase1(node *RBNode[K, V]) {
	if node.Parent == nil {
		return
	}
	tree.deleteCase2(node)
}

func (tree *RBTree[K, V]) deleteCase2(node *RBNode[K, V]) {
	sibling := node.sibling()
	if nodeColor(sibling) == red {
		node.Parent.color = red
		sibling.color = black
		if node == node.Parent.Left {
			tree.rotateLeft(node.Parent)
		} else {
			tree.rotateRight(node.Parent)
		}
	}
	tree.deleteCase3(node)
}

func (tree *RBTree[K, V]) deleteCase3(node *RBNode[K, V]) {
	sibling := node.sibling()
	if nodeColor(node.Parent) == black &&
		nodeColor(sibling) == black &&
		nodeColor(sibling.Left) == black &&
		nodeColor(sibling.Right) == black {
		sibling.color = red
		tree.deleteCase1(node.Parent)
	} else {
		tree.deleteCase4(node)
	}
}

func (tree *RBTree[K, V]) deleteCase4(node *RBNode[K, V]) {
	sibling := node.sibling()
	if nodeColor(node.Parent) == red &&
		nodeColor(sibling) == black &&
		nodeColor(sibling.Left) == black &&
		nodeColor(sibling.Right) == black {
		sibling.color = red
		node.Parent.color = black
	} else {
		tree.deleteCase5(node)
	}
}

func (tree *RBTree[K, V]) deleteCase5(node *RBNode[K, V]) {
	sibling := node.sibling()
	if node == node.Parent.Left &&
		nodeColor(sibling) == black &&
		nodeColor(sibling.Left) == red &&
		nodeColor(sibling.Right) == black {
		sibling.color = red
		sibling.Left.color = black
		tree.rotateRight(sibling)
	} else if node == node.Parent.Right &&
		nodeColor(sibling) == black &&
		nodeColor(sibling.Right) == red &&
		nodeColor(sibling.Left) == black {
		sibling.color = red
		sibling.Right.color = black
		tree.rotateLeft(sibling)
	}
	tree.deleteCase6(node)
}

func (tree *RBTree[K, V]) deleteCase6(node *RBNode[K, V]) {
	sibling := node.sibling()
	sibling.color = nodeColor(node.Parent)
	node.Parent.color = black
	if node == node.Parent.Left && nodeColor(sibling.Right) == red {
		sibling.Right.color = black
		tree.rotateLeft(node.Parent)
	} else if nodeColor(sibling.Left) == red {
		sibling.Left.color = black
		tree.rotateRight(node.Parent)
	}
}

func nodeColor[K comparable, V any](node *RBNode[K, V]) color {
	if node == nil {
		return black
	}
	return node.color
}

// RBTreeIterator holding the iterator's state
type RBTreeIterator[K comparable, V any] struct {
	tree     *RBTree[K, V]
	node     *RBNode[K, V]
	position position
}

type position byte

const (
	begin, between, end position = 0, 1, 2
)

// RBTreeIterator returns a stateful iterator whose elements are key/value pairs.
func (tree *RBTree[K, V]) Iterator() *RBTreeIterator[K, V] {
	return &RBTreeIterator[K, V]{tree: tree, node: nil, position: begin}
}

// IteratorAt returns a stateful iterator whose elements are key/value pairs that is initialised at a particular node.
func (tree *RBTree[K, V]) IteratorAt(node *RBNode[K, V]) *RBTreeIterator[K, V] {
	return &RBTreeIterator[K, V]{tree: tree, node: node, position: between}
}

// Next moves the iterator to the next element and returns true if there was a next element in the container.
// If Next() returns true, then next element's key and value can be retrieved by Key() and Value().
// If Next() was called for the first time, then it will point the iterator to the first element if it exists.
// Modifies the state of the iterator.
func (iterator *RBTreeIterator[K, V]) Next() bool {
	if iterator.position == end {
		goto end
	}
	if iterator.position == begin {
		left := iterator.tree.Left()
		if left == nil {
			goto end
		}
		iterator.node = left
		goto between
	}
	if iterator.node.Right != nil {
		iterator.node = iterator.node.Right
		for iterator.node.Left != nil {
			iterator.node = iterator.node.Left
		}
		goto between
	}
	for iterator.node.Parent != nil {
		node := iterator.node
		iterator.node = iterator.node.Parent
		if node == iterator.node.Left {
			goto between
		}
	}

end:
	iterator.node = nil
	iterator.position = end
	return false

between:
	iterator.position = between
	return true
}

// Prev moves the iterator to the previous element and returns true if there was a previous element in the container.
// If Prev() returns true, then previous element's key and value can be retrieved by Key() and Value().
// Modifies the state of the iterator.
func (iterator *RBTreeIterator[K, V]) Prev() bool {
	if iterator.position == begin {
		goto begin
	}
	if iterator.position == end {
		right := iterator.tree.Right()
		if right == nil {
			goto begin
		}
		iterator.node = right
		goto between
	}
	if iterator.node.Left != nil {
		iterator.node = iterator.node.Left
		for iterator.node.Right != nil {
			iterator.node = iterator.node.Right
		}
		goto between
	}
	for iterator.node.Parent != nil {
		node := iterator.node
		iterator.node = iterator.node.Parent
		if node == iterator.node.Right {
			goto between
		}
	}

begin:
	iterator.node = nil
	iterator.position = begin
	return false

between:
	iterator.position = between
	return true
}

// Value returns the current element's value.
// Does not modify the state of the iterator.
func (iterator *RBTreeIterator[K, V]) Value() V {
	return iterator.node.Value
}

// Key returns the current element's key.
// Does not modify the state of the iterator.
func (iterator *RBTreeIterator[K, V]) Key() K {
	return iterator.node.Key
}

// RBNode returns the current element's node.
// Does not modify the state of the iterator.
func (iterator *RBTreeIterator[K, V]) Node() *RBNode[K, V] {
	return iterator.node
}

// Begin resets the iterator to its initial state (one-before-first)
// Call Next() to fetch the first element if any.
func (iterator *RBTreeIterator[K, V]) Begin() {
	iterator.node = nil
	iterator.position = begin
}

// End moves the iterator past the last element (one-past-the-end).
// Call Prev() to fetch the last element if any.
func (iterator *RBTreeIterator[K, V]) End() {
	iterator.node = nil
	iterator.position = end
}

// First moves the iterator to the first element and returns true if there was a first element in the container.
// If First() returns true, then first element's key and value can be retrieved by Key() and Value().
// Modifies the state of the iterator
func (iterator *RBTreeIterator[K, V]) First() bool {
	iterator.Begin()
	return iterator.Next()
}

// Last moves the iterator to the last element and returns true if there was a last element in the container.
// If Last() returns true, then last element's key and value can be retrieved by Key() and Value().
// Modifies the state of the iterator.
func (iterator *RBTreeIterator[K, V]) Last() bool {
	iterator.End()
	return iterator.Prev()
}

// NextTo moves the iterator to the next element from current position that satisfies the condition given by the
// passed function, and returns true if there was a next element in the container.
// If NextTo() returns true, then next element's key and value can be retrieved by Key() and Value().
// Modifies the state of the iterator.
func (iterator *RBTreeIterator[K, V]) NextTo(f func(key K, value V) bool) bool {
	for iterator.Next() {
		key, value := iterator.Key(), iterator.Value()
		if f(key, value) {
			return true
		}
	}
	return false
}

// PrevTo moves the iterator to the previous element from current position that satisfies the condition given by the
// passed function, and returns true if there was a next element in the container.
// If PrevTo() returns true, then next element's key and value can be retrieved by Key() and Value().
// Modifies the state of the iterator.
func (iterator *RBTreeIterator[K, V]) PrevTo(f func(key K, value V) bool) bool {
	for iterator.Prev() {
		key, value := iterator.Key(), iterator.Value()
		if f(key, value) {
			return true
		}
	}
	return false
}

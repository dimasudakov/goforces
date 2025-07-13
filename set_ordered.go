package main

import (
	"cmp"
	"fmt"
	"reflect"
	"strings"
)

// original implementation: https://github.com/emirpasic/gods/tree/master

// OrderedSet holds elements in a red-black tree
type OrderedSet[T comparable] struct {
	tree *RBTree[T, struct{}]
}

var itemExists = struct{}{}

func NewOrderedSet[T cmp.Ordered](values ...T) *OrderedSet[T] {
	return NewOrderedSetWith[T](cmp.Compare[T], values...)
}

// NewOrderedSetWith instantiates a new empty set with the custom comparator.
func NewOrderedSetWith[T comparable](comparator RBComparator[T], values ...T) *OrderedSet[T] {
	set := &OrderedSet[T]{tree: NewRBTWithComp[T, struct{}](comparator)}
	if len(values) > 0 {
		set.Add(values...)
	}
	return set
}

// Add adds the items (one or more) to the set.
func (set *OrderedSet[T]) Add(items ...T) {
	for _, item := range items {
		set.tree.Put(item, itemExists)
	}
}

// Remove removes the items (one or more) from the set.
func (set *OrderedSet[T]) Remove(items ...T) {
	for _, item := range items {
		set.tree.Remove(item)
	}
}

// Contains checks weather items (one or more) are present in the set.
// All items have to be present in the set for the method to return true.
// Returns true if no arguments are passed at all, i.e. set is always superset of empty set.
func (set *OrderedSet[T]) Contains(items ...T) bool {
	for _, item := range items {
		if _, contains := set.tree.Get(item); !contains {
			return false
		}
	}
	return true
}

// Empty returns true if set does not contain any elements.
func (set *OrderedSet[T]) Empty() bool {
	return set.tree.Size() == 0
}

// Size returns number of elements within the set.
func (set *OrderedSet[T]) Size() int {
	return set.tree.Size()
}

// Clear clears all values in the set.
func (set *OrderedSet[T]) Clear() {
	set.tree.Clear()
}

// Values returns all items in the set.
func (set *OrderedSet[T]) Values() []T {
	return set.tree.Keys()
}

// String returns a string representation of container
func (set *OrderedSet[T]) String() string {
	str := "TreeSet\n"
	items := []string{}
	for _, v := range set.tree.Keys() {
		items = append(items, fmt.Sprintf("%v", v))
	}
	str += strings.Join(items, ", ")
	return str
}

// Intersection returns the intersection between two sets.
// The new set consists of all elements that are both in "set" and "another".
// The two sets should have the same comparators, otherwise the result is empty set.
// Ref: https://en.wikipedia.org/wiki/Intersection_(set_theory)
func (set *OrderedSet[T]) Intersection(another *OrderedSet[T]) *OrderedSet[T] {
	result := NewOrderedSetWith(set.tree.Comparator)

	setComparator := reflect.ValueOf(set.tree.Comparator)
	anotherComparator := reflect.ValueOf(another.tree.Comparator)
	if setComparator.Pointer() != anotherComparator.Pointer() {
		return result
	}

	// Iterate over smaller set (optimization)
	if set.Size() <= another.Size() {
		for it := set.Iterator(); it.Next(); {
			if another.Contains(it.Value()) {
				result.Add(it.Value())
			}
		}
	} else {
		for it := another.Iterator(); it.Next(); {
			if set.Contains(it.Value()) {
				result.Add(it.Value())
			}
		}
	}

	return result
}

// Union returns the union of two sets.
// The new set consists of all elements that are in "set" or "another" (possibly both).
// The two sets should have the same comparators, otherwise the result is empty set.
// Ref: https://en.wikipedia.org/wiki/Union_(set_theory)
func (set *OrderedSet[T]) Union(another *OrderedSet[T]) *OrderedSet[T] {
	result := NewOrderedSetWith(set.tree.Comparator)

	setComparator := reflect.ValueOf(set.tree.Comparator)
	anotherComparator := reflect.ValueOf(another.tree.Comparator)
	if setComparator.Pointer() != anotherComparator.Pointer() {
		return result
	}

	for it := set.Iterator(); it.Next(); {
		result.Add(it.Value())
	}
	for it := another.Iterator(); it.Next(); {
		result.Add(it.Value())
	}

	return result
}

// Difference returns the difference between two sets.
// The two sets should have the same comparators, otherwise the result is empty set.
// The new set consists of all elements that are in "set" but not in "another".
// Ref: https://proofwiki.org/wiki/Definition:Set_Difference
func (set *OrderedSet[T]) Difference(another *OrderedSet[T]) *OrderedSet[T] {
	result := NewOrderedSetWith(set.tree.Comparator)

	setComparator := reflect.ValueOf(set.tree.Comparator)
	anotherComparator := reflect.ValueOf(another.tree.Comparator)
	if setComparator.Pointer() != anotherComparator.Pointer() {
		return result
	}

	for it := set.Iterator(); it.Next(); {
		if !another.Contains(it.Value()) {
			result.Add(it.Value())
		}
	}

	return result
}

// OrderedSetIterator returns a stateful iterator whose values can be fetched by an index.
type OrderedSetIterator[T comparable] struct {
	index    int
	iterator *RBTreeIterator[T, struct{}]
	tree     *RBTree[T, struct{}]
}

// Iterator holding the iterator's state
func (set *OrderedSet[T]) Iterator() OrderedSetIterator[T] {
	return OrderedSetIterator[T]{index: -1, iterator: set.tree.Iterator(), tree: set.tree}
}

// Next moves the iterator to the next element and returns true if there was a next element in the container.
// If Next() returns true, then next element's index and value can be retrieved by Index() and Value().
// If Next() was called for the first time, then it will point the iterator to the first element if it exists.
// Modifies the state of the iterator.
func (iterator *OrderedSetIterator[T]) Next() bool {
	if iterator.index < iterator.tree.Size() {
		iterator.index++
	}
	return iterator.iterator.Next()
}

// Prev moves the iterator to the previous element and returns true if there was a previous element in the container.
// If Prev() returns true, then previous element's index and value can be retrieved by Index() and Value().
// Modifies the state of the iterator.
func (iterator *OrderedSetIterator[T]) Prev() bool {
	if iterator.index >= 0 {
		iterator.index--
	}
	return iterator.iterator.Prev()
}

// Value returns the current element's value.
// Does not modify the state of the iterator.
func (iterator *OrderedSetIterator[T]) Value() T {
	return iterator.iterator.Key()
}

// Index returns the current element's index.
// Does not modify the state of the iterator.
func (iterator *OrderedSetIterator[T]) Index() int {
	return iterator.index
}

// Begin resets the iterator to its initial state (one-before-first)
// Call Next() to fetch the first element if any.
func (iterator *OrderedSetIterator[T]) Begin() {
	iterator.index = -1
	iterator.iterator.Begin()
}

// End moves the iterator past the last element (one-past-the-end).
// Call Prev() to fetch the last element if any.
func (iterator *OrderedSetIterator[T]) End() {
	iterator.index = iterator.tree.Size()
	iterator.iterator.End()
}

// First moves the iterator to the first element and returns true if there was a first element in the container.
// If First() returns true, then first element's index and value can be retrieved by Index() and Value().
// Modifies the state of the iterator.
func (iterator *OrderedSetIterator[T]) First() bool {
	iterator.Begin()
	return iterator.Next()
}

// Last moves the iterator to the last element and returns true if there was a last element in the container.
// If Last() returns true, then last element's index and value can be retrieved by Index() and Value().
// Modifies the state of the iterator.
func (iterator *OrderedSetIterator[T]) Last() bool {
	iterator.End()
	return iterator.Prev()
}

// NextTo moves the iterator to the next element from current position that satisfies the condition given by the
// passed function, and returns true if there was a next element in the container.
// If NextTo() returns true, then next element's index and value can be retrieved by Index() and Value().
// Modifies the state of the iterator.
func (iterator *OrderedSetIterator[T]) NextTo(f func(index int, value T) bool) bool {
	for iterator.Next() {
		index, value := iterator.Index(), iterator.Value()
		if f(index, value) {
			return true
		}
	}
	return false
}

// PrevTo moves the iterator to the previous element from current position that satisfies the condition given by the
// passed function, and returns true if there was a next element in the container.
// If PrevTo() returns true, then next element's index and value can be retrieved by Index() and Value().
// Modifies the state of the iterator.
func (iterator *OrderedSetIterator[T]) PrevTo(f func(index int, value T) bool) bool {
	for iterator.Prev() {
		index, value := iterator.Index(), iterator.Value()
		if f(index, value) {
			return true
		}
	}
	return false
}

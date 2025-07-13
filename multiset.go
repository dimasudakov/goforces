package main

type Multiset[T comparable] struct {
	data map[T]int
	size int
}

func NewMultiset[T comparable]() *Multiset[T] {
	return &Multiset[T]{
		data: make(map[T]int),
	}
}

func (m *Multiset[T]) Add(values ...T) {
	for _, v := range values {
		m.data[v]++
		m.size++
	}
}

func (m *Multiset[T]) Remove(val T) {
	if count, ok := m.data[val]; ok {
		if count > 1 {
			m.data[val]--
		} else {
			delete(m.data, val)
		}
		m.size--
	}
}

func (m *Multiset[T]) RemoveAll(val T) {
	if count, ok := m.data[val]; ok {
		delete(m.data, val)
		m.size -= count
	}
}

func (m *Multiset[T]) Has(val T) bool {
	_, ok := m.data[val]
	return ok
}

func (m *Multiset[T]) Count(val T) int {
	return m.data[val]
}

func (m *Multiset[T]) Len() int {
	return m.size
}

func (m *Multiset[T]) UniqueValues() []T {
	res := make([]T, 0, len(m.data))
	for val := range m.data {
		res = append(res, val)
	}
	return res
}

func (m *Multiset[T]) AllValues() []T {
	res := make([]T, 0, m.size)
	for val, count := range m.data {
		for i := 0; i < count; i++ {
			res = append(res, val)
		}
	}
	return res
}

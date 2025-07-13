package main

type Set[T comparable] struct {
	data map[T]struct{}
}

func NewSet[T comparable](values ...T) *Set[T] {
	s := &Set[T]{data: make(map[T]struct{}, len(values))}
	s.Add(values...)
	return s
}

func NewSetN[T comparable](n int) *Set[T] {
	return &Set[T]{data: make(map[T]struct{}, n)}
}

func (s *Set[T]) Add(val ...T) {
	for _, v := range val {
		s.data[v] = struct{}{}
	}
}

func (s *Set[T]) Remove(val T) {
	delete(s.data, val)
}

func (s *Set[T]) Contains(val T) bool {
	_, ok := s.data[val]
	return ok
}

func (s *Set[T]) Len() int {
	return len(s.data)
}

func (s *Set[T]) Empty() bool {
	return s.Len() == 0
}

func (s *Set[T]) Clear() {
	s.data = make(map[T]struct{})
}

func (s *Set[T]) Values() []T {
	res := make([]T, 0, len(s.data))
	for k := range s.data {
		res = append(res, k)
	}
	return res
}

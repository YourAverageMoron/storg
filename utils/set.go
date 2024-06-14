package utils

type Set[T interface{ comparable }] struct {
	list map[T]struct{}
}

func (s *Set[T]) Has(v T) bool {
	_, ok := s.list[v]
	return ok
}

func (s *Set[T]) Add(v T) {
	s.list[v] = struct{}{}
}

func (s *Set[T]) Remove(v T) {
	delete(s.list, v)
}

func (s *Set[T]) Clear() {
	s.list = make(map[T]struct{})
}

func (s *Set[T]) Size() int {
	return len(s.list)
}

func NewSet[T interface{ comparable }]() *Set[T] {
	s := &Set[T]{}
	s.list = make(map[T]struct{})
	return s
}

func (s *Set[T]) AddMulti(list ...T) {
	for _, v := range list {
		s.Add(v)
	}
}

func (s *Set[T]) Iterate() []T {
	v := make([]T, 0, len(s.list))
	for key := range s.list {
		v = append(v, key)
	}
    return v
}

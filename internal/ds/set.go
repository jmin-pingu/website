package ds

type Set[T comparable] map[T]bool

func (s Set[T]) Add(v T) {
	(s)[v] = true
}

func (s Set[T]) Remove(v T) {
	delete(s, v)
}

func (s Set[T]) Contains(v T) bool {
	_, ok := (s)[v]
	return ok
}

func (s Set[T]) IsEmpty() bool {
	return len(s) == 0
}

func (s1 Set[T]) Union(s2 Set[T]) Set[T] {
	s := make(Set[T], 0)
	for k, _ := range s1 {
		s[k] = true
	}
	for k, _ := range s2 {
		s[k] = true
	}
	return s
}

func (s1 Set[T]) Intersection(s2 Set[T]) Set[T] {
	s := make(Set[T], 0)
	for k, _ := range s1 {
		if s2.Contains(k) {
			s[k] = true
		}
	}
	return s
}

func (s1 Set[T]) SubsetOf(s2 Set[T]) bool {
	for k, _ := range s1 {
		if !s2.Contains(k) {
			return false
		}
	}
	return true
}

func (s1 Set[T]) ToList() []T {
	var s []T
	for k, _ := range s1 {
		s = append(s, k)
	}
	return s
}

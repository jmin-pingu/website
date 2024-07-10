package ds

type Set[T comparable] map[T]bool

// Method: Add
func (s Set[T]) Add(v T) {
	(s)[v] = true
}

// Method: Remove
func (s Set[T]) Remove(v T) {
	delete(s, v)
}

// Contains:
func (s Set[T]) Contains(v T) bool {
	_, ok := (s)[v]
	return ok
}

// IsEmpty:
func (s Set[T]) IsEmpty() bool {
	return len(s) == 0
}

// Union:
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

// Intersection:
func (s1 Set[T]) Intersection(s2 Set[T]) Set[T] {
	s := make(Set[T], 0)
	for k, _ := range s1 {
		if s2.Contains(k) {
			s[k] = true
		}
	}
	return s
}

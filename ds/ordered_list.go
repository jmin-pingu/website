package ds

type OrderedList[T comparable] []T

func (l *OrderedList[T]) Add(v T) {
	if !l.Contains(v) {
		*l = append(*l, v)
	}
}

func (l *OrderedList[T]) index(v T) int {
	for i, entry := range *l {
		if entry == v {
			return i
		}
	}
	return -1
}

func (l *OrderedList[T]) Remove(v T) {
	i := l.index(v)
	if i != -1 {
		*l = append((*l)[:i], (*l)[i+1:]...)
	}
}

func (l *OrderedList[T]) Contains(v T) bool {
	i := l.index(v)
	if i == -1 {
		return false
	} else {
		return true
	}
}

package ds

import "fmt"

type StrictDict[K comparable, V any] struct {
	Categories []K
	Values     [][]V
}

func NewStrictDict[K comparable, V any](categories []K) (StrictDict[K, V], error) {
	seen := make(map[K]bool)
	for _, c := range categories {
		if seen[c] {
			return StrictDict[K, V]{}, fmt.Errorf("categories %v are not unique", categories)
		}
		seen[c] = true
		fmt.Println(c)
	}
	return StrictDict[K, V]{Categories: categories, Values: make([][]V, len(categories))}, nil
}

func (sd *StrictDict[K, V]) Append(key K, value V) error {
	for i, c := range sd.Categories {
		if c == key {
			sd.Values[i] = append(sd.Values[i], value)
			return nil
		}
	}
	return fmt.Errorf("the key %v does not exist", key)
}

func (sd *StrictDict[K, V]) Prepend(key K, value V) error {
	for i, c := range sd.Categories {
		if c == key {
			sd.Values[i] = append(sd.Values[i], value)
			return nil
		}
	}
	return fmt.Errorf("the key %v does not exist", key)
}

func (sd *StrictDict[K, V]) Get(key K) ([]V, error) {
	for i, c := range sd.Categories {
		if c == key {
			return sd.Values[i], nil
		}
	}
	return nil, fmt.Errorf("the key %v does not exist", key)
}

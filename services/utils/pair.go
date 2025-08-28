package utils

type Pair[T interface{}] struct {
	Key   T
	Value T
}

func (p Pair[T]) Get() (T, T) {
	return p.Key, p.Value
}

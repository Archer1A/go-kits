package set

type set[T comparable] struct {
	m map[T]struct{}
}

func New[T comparable]() *set[T] {
	return &set[T]{
		m: make(map[T]struct{}),
	}
}

func (set *set[T]) Add(key T) {
	set.m[key] = struct{}{}
}

func (set *set[T]) Delete(key T) {
	delete(set.m, key)
}

func (set *set[T]) Contain(key T) bool {
	_, ok := set.m[key]
	return ok
}

func (set *set[T]) Clear() {
	set.m = make(map[T]struct{})
}

func (set *set[T]) Len() int {
	return len(set.m)
}

func (set *set[T]) Values() []T {
	arr := make([]T, 0, len(set.m))
	for k, _ := range set.m {
		arr = append(arr, k)
	}
	return arr
}

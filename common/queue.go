package common

type Queue[T any] struct {
	data []T
}

func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{}
}
func (q *Queue[T]) Range(fn func(v T)) {
	for _, v := range q.data {
		fn(v)
	}
}
func (q *Queue[T]) PushBack(v T) {
	q.data = append(q.data, v)
}

func (q *Queue[T]) PopBack() (v T) {
	if q.Size() == 0 {
		return v
	}
	v = q.data[len(q.data)-1]
	q.data = q.data[:len(q.data)-1]
	return v
}

func (q *Queue[T]) Index(i int) T {
	return q.data[i]
}

func (q *Queue[T]) Front() (v T) {
	if q.Size() == 0 {
		return
	}
	v = q.data[0]
	q.data = q.data[1:]
	return v
}

func (q *Queue[T]) PushFront(v T) {
	tmp := make([]T, len(q.data)+1)
	tmp[0] = v
	copy(tmp[1:], q.data)
}

func (q *Queue[T]) Size() int {
	return len(q.data)
}

func (q *Queue[T]) FindIndex(v T, eq func(v1, v2 T) bool) (i int) {
	for i, v1 := range q.data {
		if eq(v1, v) {
			return i
		}
	}
	return -1
}

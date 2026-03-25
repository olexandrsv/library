package trace

type TraceIterator[T any] interface {
	Next() (T, bool)
}

type traceIterator[T any] struct {
	items    []T
	position int
}

func NewIterator[T any](items []T) TraceIterator[T] {
	return &traceIterator[T]{
		items: items,
	}
}

func (i *traceIterator[T]) Next() (T, bool) {
	var item T
	var hasNext bool
	if i.position+1 <= len(i.items)-1 {
		hasNext = true
	}
	if i.position <= len(i.items)-1 {
		item = i.items[i.position]
		i.position++
	}
	return item, hasNext
}

func IteratorToSlice[T any](iterator TraceIterator[T]) []T {
	if iterator == nil {
		return nil
	}
	var items []T
	for {
		item, ok := iterator.Next()
		items = append(items, item)

		if !ok {
			return items
		}
	}
}

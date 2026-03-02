package request

import "library/errors"

type ValueRequest[T any] struct {
	value          T
	defaultValue   *T
	errorContainer errors.ErrorContainer
}

func NewValueRequest[T any](c errors.ErrorContainer, value T) *ValueRequest[T] {
	return &ValueRequest[T]{
		value:          value,
		errorContainer: c,
	}
}

func (r *ValueRequest[T]) DefaultValue(t T) *ValueRequest[T] {
	r.defaultValue = &t
	return r
}

func (r *ValueRequest[T]) Do() T {
	var nilValue T
	if r.errorContainer.FatalErr() != nil {
		return nilValue
	}

	containerEmpty := r.errorContainer.Error() == nil
	if !containerEmpty && r.defaultValue != nil{
		return *r.defaultValue
	}
	if !containerEmpty {
		return nilValue
	}

	return r.value
}

type Builder[T any] struct {
	c          errors.ErrorContainer
	value      T
}

func NewRequestBuilder[T any](c errors.ErrorContainer) *Builder[T] {
	return &Builder[T]{
		c: c,
	}
}

func (b *Builder[T]) WithValue(value T) *Builder[T] {
	b.value = value
	return b
}

func (b *Builder[T]) Create() *ValueRequest[T] {
	return &ValueRequest[T]{
		errorContainer: b.c,
		value:          b.value,
	}
}

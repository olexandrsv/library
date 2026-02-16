package request

import "library/errors"

type ValueRequest[T any] struct {
	value          T
	defaultValue   *T
	errorContainer errors.ErrorContainer
	fatalErr       error
	err            error
	isCritical     bool
}

func NewValueRequest[T any](c errors.ErrorContainer, fatalErr error, err error, value T) *ValueRequest[T] {
	return &ValueRequest[T]{
		value:          value,
		errorContainer: c,
		fatalErr:       fatalErr,
		err:            err,
	}
}

func (r *ValueRequest[T]) DefaultValue(t T) *ValueRequest[T] {
	r.defaultValue = &t
	return r
}

func (r *ValueRequest[T]) OrFatalError() *ValueRequest[T] {
	r.isCritical = true
	return r
}

func (r *ValueRequest[T]) Do() T {
	var nilValue T
	if r.errorContainer.FatalErr() != nil {
		return nilValue
	}
	if r.fatalErr != nil {
		r.errorContainer.AddFatalError(r.fatalErr)
		return nilValue
	}

	if r.err != nil {
		if r.isCritical {
			r.errorContainer.AddFatalError(r.err)
			return nilValue
		}

		if r.defaultValue != nil {
			return *r.defaultValue
		}

		r.errorContainer.AddError(r.err)
		return nilValue
	}

	return r.value
}

type Builder[T any] struct {
	c          errors.ErrorContainer
	fatalError error
	err        error
	value      T
}

func NewRequestBuilder[T any](c errors.ErrorContainer) *Builder[T] {
	return &Builder[T]{
		c: c,
	}
}

func (b *Builder[T]) WithFatalError(err error) *Builder[T] {
	b.fatalError = err
	return b
}

func (b *Builder[T]) WithError(err error) *Builder[T] {
	b.err = err
	return b
}

func (b *Builder[T]) WithValue(value T) *Builder[T] {
	b.value = value
	return b
}

func (b *Builder[T]) Create() *ValueRequest[T] {
	return &ValueRequest[T]{
		errorContainer: b.c,
		value:          b.value,
		err:            b.err,
		fatalErr:       b.fatalError,
	}
}

package types

type Option[T any] struct {
	value *T
}

func Some[T any](value T) Option[T] {
	return Option[T]{
		value: &value,
	}
}

func None[T any]() Option[T] {
	return Option[T]{
		value: nil,
	}
}

func (o *Option[T]) IsSome() bool {
	return o.value != nil
}

func (o *Option[T]) IsNone() bool {
	return o.value == nil
}

func (o *Option[T]) Unwrap() T {
	if o.IsNone() {
		panic("called `Option.Get()` on a `None` value")
	}

	return *o.value
}

func (o *Option[T]) Unchecked() *T {
	return o.value
}

func (o *Option[T]) Take() Option[T] {
	if o.IsNone() {
		return Option[T]{
			value: nil,
		}
	}

	ret := Option[T]{
		value: o.value,
	}

	o.value = nil

	return ret
}

func (o *Option[T]) Replace(value T) Option[T] {
	tmp := *o

	o.value = &value

	return tmp
}

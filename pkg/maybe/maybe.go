package maybe

type MayBe[T any] struct {
	val T
	err error
}

func From[T any](v T, err error) MayBe[T] {
	return MayBe[T]{
		val: v,
		err: err,
	}
}

func FromVal[T any](v T) MayBe[T] {
	return From(v, nil)
}

func Map[I, J any](m MayBe[I], f func(v I) (J, error)) MayBe[J] {
	if m.err != nil {
		return MayBe[J]{err: m.err}
	}

	return From[J](f(m.val))
}

func (m MayBe[T]) Then(f func(v T) (T, error)) MayBe[T] {
	if m.err != nil {
		return m
	}

	return From(f(m.val))
}

func (m MayBe[T]) Catch(f func(error) error) MayBe[T] {
	if m.err != nil {
		return From(m.val, f(m.err))
	}

	return m
}

func (m MayBe[T]) Do(f func() error) MayBe[T] {
	if err := f(); err != nil {
		return From(m.val, err)
	}

	return m
}

func (m MayBe[T]) Unpack() (T, error) {
	return m.val, m.err
}

func (m MayBe[T]) Val() T {
	return m.val
}

func (m MayBe[T]) Err() error {
	return m.err
}

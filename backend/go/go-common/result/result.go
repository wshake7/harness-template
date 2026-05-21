package result

type Result[T any] struct {
	value T
	err   error
}

func Ok[T any](value T) Result[T] {
	return Result[T]{value: value}
}

func Err[T any](err error) Result[T] {
	return Result[T]{err: err}
}

func (r Result[T]) IsOk() bool {
	return r.err == nil
}

func (r Result[T]) IsErr() bool {
	return r.err != nil
}

func (r Result[T]) Get() T {
	if r.err != nil {
		panic(r.err)
	}
	return r.value
}

func (r Result[T]) GetOr(defaultValue T) T {
	if r.err != nil {
		return defaultValue
	}
	return r.value
}

func (r Result[T]) Map(f func(T) T) Result[T] {
	if r.IsOk() {
		return Ok(f(r.value))
	}
	return r
}

func (r Result[T]) MapErr(f func(error) error) Result[T] {
	if r.IsErr() {
		return Err[T](f(r.err))
	}
	return r
}

func (r Result[T]) Expect() {
	if r.IsErr() {
		panic(r.err)
	}
}

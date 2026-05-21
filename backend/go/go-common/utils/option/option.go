package option

import (
	"errors"
	"fmt"
)

type Opt[T any] struct {
	V T
	B bool
}

func OptOf[T any](t T, b bool) Opt[T] {
	return Opt[T]{t, b}
}

func OptOfEmpty[T any]() Opt[T] {
	return Opt[T]{}
}

func Some[T any](t T) Opt[T] {
	return Opt[T]{t, true}
}

func None[T any]() Opt[T] {
	return Opt[T]{}
}

func (o Opt[T]) Unravel() (T, bool) {
	return o.V, o.B
}

func (o Opt[T]) IsSome() bool {
	return o.B
}

func (o Opt[T]) IsNone() bool {
	return !o.B
}

func (o Opt[T]) Expect() T {
	if o.IsSome() {
		return o.V
	}
	panic("option is none")
}

func (o Opt[T]) ExpectErr(msg error) T {
	if o.IsSome() {
		return o.V
	}
	panic(msg)
}

func (o Opt[T]) ExpectMsg(msg string) T {
	if o.IsSome() {
		return o.V
	}
	panic(msg)
}

func (o Opt[T]) Get() T {
	if o.IsSome() {
		return o.V
	}
	panic("option is none")
}

func (o Opt[T]) GetOr(t T) T {
	if o.IsSome() {
		return o.V
	}
	return t
}

func (o Opt[T]) GetOrDefault() T {
	if o.IsSome() {
		return o.V
	}
	return *new(T)
}

func (o Opt[T]) GetElse(fn func() T) T {
	if o.IsSome() {
		return o.V
	}
	return fn()
}

func (o Opt[T]) GetOrElse(fn0 func() T, fn1 func(t T) T) T {
	if o.IsSome() {
		return fn1(o.V)
	}
	return fn0()
}

func (o Opt[T]) Map(fn func(t T)) {
	if o.IsSome() {
		fn(o.V)
	}
}

func (o Opt[T]) MapOrElse(fn0 func(t T), fn1 func()) {
	if o.IsSome() {
		fn0(o.V)
	} else {
		fn1()
	}
}

func (o Opt[T]) String() string {
	if o.IsSome() {
		return fmt.Sprintf("some(%v)", o.V)
	}
	return "none"
}

// NzOpt Non zero option
type NzOpt[T comparable] struct {
	V T
}

func NzOptOf[T comparable](t T) NzOpt[T] {
	return NzOpt[T]{t}
}

func NzOptOfEmpty[T comparable]() NzOpt[T] {
	return NzOpt[T]{}
}

func (o NzOpt[T]) D() (T, bool) {
	return o.V, o.V != *new(T)
}

func (o NzOpt[T]) IsSome() bool {
	return o.V != *new(T)
}

func (o NzOpt[T]) IsNone() bool {
	return o.V == *new(T)
}

func (o NzOpt[T]) Expect() T {
	if o.IsSome() {
		return o.V
	}
	panic("option is none")
}

func (o NzOpt[T]) ExpectErr(err error) T {
	if o.IsSome() {
		return o.V
	}
	panic(err)
}

func (o NzOpt[T]) ExpectString(msg string) T {
	if o.IsSome() {
		return o.V
	}
	panic(msg)
}

func (o NzOpt[T]) ToOpt() Opt[T] {
	if o.IsSome() {
		return Some(o.V)
	}
	return None[T]()
}

// Get 获取值 如果为none 则会panic
func (o NzOpt[T]) Get() T {
	if o.IsSome() {
		return o.V
	}
	panic(errors.New("option is none"))
}

func (o NzOpt[T]) GetOr(t T) T {
	if o.IsSome() {
		return o.V
	}
	return t
}

func (o NzOpt[T]) GetElse(fn func() T) T {
	if o.IsSome() {
		return o.V
	}
	return fn()
}

func (o NzOpt[T]) MapOrElse(fn0 func(t T), fn1 func()) {
	if o.IsSome() {
		fn0(o.V)
	} else {
		fn1()
	}
}

func (o NzOpt[T]) Map(fn func(t T)) {
	if o.IsSome() {
		fn(o.V)
	}
}

func (o NzOpt[T]) String() string {
	if o.IsSome() {
		return fmt.Sprintf("some(%v)", o.V)
	}
	return "none"
}

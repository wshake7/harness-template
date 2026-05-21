package coroutine

import "go-common/utils/function"

// Launch 禁止在程序中直接使用go关键字，若需要用这个代替
func Launch(fn func()) {
	go func() {
		function.RecFn(func() {
			fn()
		})
	}()
}

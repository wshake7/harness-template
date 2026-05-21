package function

import (
	"go.uber.org/zap"
)

func RecObjFn[T any](fn func() T) T {
	defer func() {
		if err := recover(); err != nil {
			zap.S().Errorf("RecObjFn failed: %v", err)
		}
	}()
	return fn()
}

func RecFn(fn func()) {
	defer func() {
		if err := recover(); err != nil {
			zap.S().Errorf("RecFn failed: %v", err)
		}
	}()
	fn()
}

package delay_task

import (
	"go-common/collection"
	"go-common/utils/function"
	"time"
)

var (
	tasks = collection.SyncMap[uint32, *time.Timer]{}
	idInc = uint32(0)
)

func Add(delayTime time.Duration, fn func()) uint32 {
	idInc++
	fn0 := func() {
		function.RecFn(func() {
			fn()
			tasks.LoadAndDelete(idInc)
		})
	}
	tasks.Store(idInc, time.AfterFunc(delayTime, fn0))
	return idInc
}

func Stop(id uint32) {
	value, loaded := tasks.LoadAndDelete(id)
	if loaded {
		value.Stop()
	}
}

func Reset(id uint32, delayTime time.Duration) {
	value, ok := tasks.Load(id)
	if ok {
		value.Reset(delayTime)
	}
}

package lifecycle

import "sync"

var (
	shutdownOnce sync.Once
	shutdownDone = make(chan struct{})
)

func BeginShutdown() {
	shutdownOnce.Do(func() {
		close(shutdownDone)
	})
}

func ShutdownDone() <-chan struct{} {
	return shutdownDone
}

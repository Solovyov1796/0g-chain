package cmd

import (
	"runtime"
	"sync"
)

func safeStartGoroutine(do func()) {
	defer func() {
		if r := recover(); r != nil {
			println("Caught panic in goroutine: ", r)

			buf := make([]byte, 4*1024)
			n := runtime.Stack(buf, true)
			println("panic stack trace:\n%s\n", buf[:n])
		}
	}()

	do()
}

func StartDeamon(do func()) {
	var wg sync.WaitGroup

	for {
		go func() {
			defer func() {
				wg.Done()
			}()
			safeStartGoroutine(do)
		}()
		wg.Add(1)
		wg.Wait()
	}
}

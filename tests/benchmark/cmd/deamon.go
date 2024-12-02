package cmd

import (
	"sync"
)

func safeStartGoroutine(do func()) {
	defer func() {
		if r := recover(); r != nil {
			println("Caught panic in goroutine: ", r)
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

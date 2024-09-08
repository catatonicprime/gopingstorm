package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		begin_capture()
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		Display()
	}()
	wg.Wait()
	//TODO: Start UI coroutine
	//TODO: Start WebUI coroutine
	fmt.Printf("goping complete\n")
}

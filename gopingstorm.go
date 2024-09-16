package main

import (
	"fmt"
	"sync"
)

func main() {
	options, err := parseArgs()
	if err != nil {
		fmt.Printf("%s\n", err)
	}
	fmt.Printf("%+v\n", options)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		begin_capture(options)
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

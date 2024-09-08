package main

import (
	"fmt"
	"time"
)

func Display() {
	for {
		fmt.Printf("Length of arpCache: %d\n", len(arpCache))
		time.Sleep(2 * time.Second)
	}
}

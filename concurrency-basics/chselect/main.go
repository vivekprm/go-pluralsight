package main

import (
	"fmt"
	"time"
)

func main() {
	ch1 := make(chan string)
	ch2 := make(chan string)

	go func() {
		for {
			ch1 <- "every 500ms"
			time.Sleep(500 * time.Millisecond)
		}
	}()
	go func() {
		for {
			ch2 <- "every two second"
			time.Sleep(2 * time.Second)
		}
	}()

	for {
		fmt.Println(<-ch1)
		fmt.Println(<-ch2)
	}
}

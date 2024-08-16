package main

import (
	"fmt"
	"sync"
)

var wg sync.WaitGroup

func main() {
	done := make(chan interface{})
	defer close(done)

	cows := make(chan interface{}, 100)
	pigs := make(chan interface{}, 100)

	// Producers
	go func() {
		for {
			select {
			case <-done:
				return
			case cows <- "moo":
			}
		}
	}()
	go func() {
		for {
			select {
			case <-done:
				return
			case pigs <- "oink":
			}
		}
	}()

	// consumers
	wg.Add(1)
	go consumeCows(done, cows)
	wg.Add(1)
	go consumePigs(done, pigs)

	wg.Wait()
}

func consumeCows(done, cows <-chan interface{}) {
	defer wg.Done()
	for {
		select {
		case <-done:
			return
		case cow, ok := <-cows:
			if !ok {
				fmt.Println("Channel closed")
				return
			}
			// Complex logic
			fmt.Println(cow)
		}
	}
}

func consumePigs(done, pigs <-chan interface{}) {
	defer wg.Done()
	for {
		select {
		case <-done:
			return
		case pig, ok := <-pigs:
			if !ok {
				fmt.Println("Channel closed")
				return
			}
			// Complex logic
			fmt.Println(pig)
		}
	}
}

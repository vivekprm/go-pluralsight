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
	for val := range orDone(done, cows) {
		fmt.Println(val)
	}
}

func consumePigs(done, pigs <-chan interface{}) {
	defer wg.Done()
	for val := range orDone(done, pigs) {
		fmt.Println(val)
	}
}

func orDone(done, c <-chan interface{}) <-chan interface{} {
	relayStream := make(chan interface{})
	go func() {
		defer close(relayStream)
		for {
			select {
			case <-done:
				return
			case val, ok := <-c:
				if !ok {
					return
				}
				// Put on relayStream
				select {
				case relayStream <- val:
				// Reason we need to put this case here is, if we don't have this and we
				// put a value on releayStream e.g. and caller of orDone never read that
				// value, this goroutine will be blocked because relayStream is unbuffered
				// channel.
				case <-done:
					return
				}
			}
		}
	}()
	return relayStream
}

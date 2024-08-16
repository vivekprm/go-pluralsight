package main

import "fmt"

func main() {
	evenCh := make(chan int)
	oddCh := make(chan int)

	go getNumber(evenCh, oddCh)

	for {
		fmt.Println(<-evenCh)
		fmt.Println(<-oddCh)
	}
}

func getNumber(evenCh, oddCh chan int) {
	i := 0
	for {
		if i%2 == 0 {
			evenCh <- i
		} else {
			oddCh <- i
		}
		i++
	}
}

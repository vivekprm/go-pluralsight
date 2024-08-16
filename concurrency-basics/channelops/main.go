package main

import "fmt"

func main() {
	c := make(chan string, 2)
	// Send on channel
	c <- "hello"
	c <- "world"
	// Receive from channel
	msg := <-c
	fmt.Println(msg)

	msg = <-c
	fmt.Println(msg)
}

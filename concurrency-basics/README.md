https://www.youtube.com/watch?v=LvgVSSpwND8

Parallelism is running multiple things at the same time on multi-core processor.

Concurrency is about breaking up a program into independently executing tasks that could potentially run at the same time and still getting the right result at the end.

So a concurrent program is one that can be parallelized.

Let's conside below program:
```go
package main

import (
	"fmt"
	"time"
)

func main() {
	count("sheep")
	count("fish")
}

func count(thing string) {
	for i := 1; true; i++ {
		fmt.Println(i, thing)
		time.Sleep(time.Millisecond * 500)
	}
}
```

It's not a concurrent function. It's going to execute first count function and wait for it to finish before it executes second count function call. But the first count function never finishes, so it's gonna count sheep forever and never get to the fish.

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	go count("sheep")
	count("fish")
}

func count(thing string) {
	for i := 1; true; i++ {
		fmt.Println(i, thing)
		time.Sleep(time.Millisecond * 500)
	}
}
```

However if we change first call to a goroutine. It won't wait for it to finish before moving to the next line. If we change the second call also in the separate goroutine, nothing is printed. As there is no way to keep main goroutine in waiting before other goroutines finish.

```go
func main() {
	go count("sheep")
	go count("fish")
}

func count(thing string) {
	for i := 1; true; i++ {
		fmt.Println(i, thing)
		time.Sleep(time.Millisecond * 500)
	}
}
```

we can add scanln call to fix that it's blocking call, so it prevents main function from exiting until we press Enter.

```go
func main() {
	go count("sheep")
	go count("fish")
	fmt.Scanln()
}

func count(thing string) {
	for i := 1; true; i++ {
		fmt.Println(i, thing)
		time.Sleep(time.Millisecond * 500)
	}
}
```

However, it's not very practical solution as it requires manual user input. What we can do instead is to use WaitGroup.

```go
func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		count("sheep")
		wg.Done()
	}()
	wg.Wait()
}

func count(thing string) {
	for i := 1; i < 5; i++ {
		fmt.Println(i, thing)
		time.Sleep(time.Millisecond * 500)
	}
}
```

What we need next is Channel. **Channels** are a way for goroutines to communicate with each other. So what if instead of printing in count we want to return a value to main goroutine. We can use channel as below:

```go
 func main() {
	c := make(chan string)
	go func() {
		count("sheep", c)
	}()
	msg := <-c
	fmt.Println(msg)
}

func count(thing string, c chan string) {
	for i := 1; i < 5; i++ {
		c <- thing
		time.Sleep(time.Millisecond * 500)
	}
}
```

In this case main go routine will wait until msg receive anything on the channel send or receive on a channel is blocking. So blocking nature of channels allow us to use them to synchronize go routines.

Imagine we have two independent goroutines.

In goroutine1 we receive something from the channel at some point and in goroutine2 we send somthing on the channel. Before sending or receving on the channel both goroutines execute independently. But once goroutine1 reaches the line where it's receiving from the channel it stops. And similary goroutine2 stops where we send on the channel and then they will be able to communicate.

Above code received just one message but if we want to receive all of them we can loop over channel.

```go
func main() {
	c := make(chan string)
	go func() {
		count("sheep", c)
	}()

	for {
		msg := <-c
		fmt.Println(msg)
	}
}

func count(thing string, c chan string) {
	for i := 1; i < 5; i++ {
		c <- thing
		time.Sleep(time.Millisecond * 500)
	}
    close(c)
}
```

We get the output but at the end we get a deadlock because count function is finished but main goroutine has no way to know so it is waiting to receive on the channel but nothing else is ever gonna send a message on the channel. Go was able to detect this problem at runtime not at compile time. To solve this we can close the channel. We never close the channel from receiving end. Because as a receiver we are not sure whether sender is finished or not. If you close the channel prematurely as a receiver and then the sender tries to send on the closed channel it will panic.

Wne we recieve from a channel , we can receive a second value which tells whether channel is open or closed. If it's not open we can break out of the for loop.

```go
func main() {
	c := make(chan string)
	go func() {
		count("sheep", c)
	}()

	for {
		msg, open := <-c
		if !open {
			break
		}
		fmt.Println(msg)
	}
}
```

There is a cleaner way to it by iterating over range of channel.
```go
func main() {
	c := make(chan string)
	go func() {
		count("sheep", c)
	}()

	for msg := range c {
		fmt.Println(msg)
	}
}
```

So here we don't need to manually check if channel is closed.

We looked at how channel operations are blocking. Let's look at something simple.
```go
func main() {
	c := make(chan string)
	// Send on channel
	c <- "hello"
	// Receive from channel
	msg := <-c
	fmt.Println(msg)
}
```

We might think that it will work. But it's going to deadlock again as sending on channel is blocking and there is no other goroutine to receive on the channel and the code never progresses to receive line.  

So alternatively we can receive on other goroutine or we can use a buffer channel. So we can give a buffer size and it won't block until the buffer is full.
```go
func main() {
	c := make(chan string, 2)
	// Send on channel
	c <- "hello"
    c <- " world"
	// Receive from channel
	msg := <-c
	fmt.Println(msg)
    msg = <-c
	fmt.Println(msg)
}
```

# Select statement
If we have two goroutines and two channels as below:
```go
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
```

We see we get one and the other, eventhough the first go routine is ready to send. It is because we are gonna block eachtime waiting for the slower one. So everytime we try to receive from channel 2 we are gonna wait to 2 seconds. So it's really slowing down the first goroutine.

So to fix that instead of receiving on in an infinite for loop in main goroutine we can use select statement, which allows us to receive from whichever channel is ready.


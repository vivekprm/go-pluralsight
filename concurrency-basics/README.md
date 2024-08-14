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
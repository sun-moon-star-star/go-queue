# go_queue
go_queue

```golang
package main

import (
	"fmt"

	"github.com/sun-moon-star-star/go_queue"
)

func main() {
	queue := go_queue.New()

	queue.Append("sun-moon-star-star")

	len := queue.Len()

	element := queue.Remove()
	queue.Done()

	elementString := element.(string)

	fmt.Println(len, elementString)
}
```

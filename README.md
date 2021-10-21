# go-cache

`go get -u github.com/1makarov/go-cache`

```go
package main

import (
	"fmt"
	"log"
	"time"
	"github.com/1makarov/go-cache"
)

func main() {
	c := cache.New()
	defer c.Close()

	if err := c.SetWithExpire("userId", 42, time.Second*5); err != nil {
		log.Fatal(err)
	}

	userId, err := c.Get("userId")
	if err != nil { // err == nil
		log.Fatal(err)
	}
	fmt.Println(userId) // Output: 42

	time.Sleep(time.Second * 6) // прошло 5 секунд 

	userId, err = c.Get("userId")
	if err != nil { // err != nil
		log.Fatal(err) // сработает этот код
	}
}
```

## golang in-memory cache

```go
package main

import (
	"fmt"
	"github.com/1makarov/go-cache/cache"
)

func main() {
	c := cache.New()

	c.Set("userId", 42)
	userId := c.Get("userId")

	fmt.Println(userId)

	c.Delete("userId")
	userId = c.Get("userId")

	fmt.Println(userId)
}
```
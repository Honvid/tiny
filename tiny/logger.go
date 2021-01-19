package tiny

import (
	"fmt"
	"time"
)

func Logger() Handler {
	return func(c *Context) {
		t := time.Now()
		fmt.Printf("[%d] %s in %v\n", c.StatusCode, c.Request.RequestURI, time.Since(t))
		c.Next()
		fmt.Printf("[%d] %s in %v\n", c.StatusCode, c.Request.RequestURI, time.Since(t))
	}
}

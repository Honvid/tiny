package tiny

import (
	"fmt"
	"time"
)

func Logger() Handler {
	return func(c *Context) {
		// Start timer
		t := time.Now()
		fmt.Printf("[%d] %s in %v\n", c.StatusCode, c.Request.RequestURI, time.Since(t))
		// Process request
		c.Next()
		// Calculate resolution time
		fmt.Printf("[%d] %s in %v\n", c.StatusCode, c.Request.RequestURI, time.Since(t))
	}
}

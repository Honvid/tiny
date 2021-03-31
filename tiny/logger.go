package tiny

import (
	"log"
	"time"
)

func Logger() Handler {
	return func(c *Context) {
		t := time.Now()
		log.Printf("REQUEST_IN [%d] %s\n", c.StatusCode, c.Request.RequestURI)
		c.Next()
		log.Printf("REQUEST_OUT [%d] %s COST %v\n", c.StatusCode, c.Request.RequestURI, time.Since(t))
	}
}

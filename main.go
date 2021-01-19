package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	"tiny/tiny"
)

func onlyForV2() tiny.Handler {
	return func(c *tiny.Context) {
		// Start timer
		t := time.Now()
		// if a server error occurred
		c.Fail(500, "Internal Server Error")
		// Calculate resolution time
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Request.RequestURI, time.Since(t))
	}
}

func main() {

	r := tiny.Default()

	r.GET("/", func(c *tiny.Context) {
		c.HTML(http.StatusOK, "<h1>Hello world!</h1>")
	})

	r.POST("/hello", func(c *tiny.Context) {
		c.JSON(http.StatusOK, tiny.H{
			"username": c.Post("username"),
			"password": c.Post("password"),
		})
	})

	r.GET("/pat/:name", func(c *tiny.Context) {
		c.JSON(http.StatusOK, tiny.H{"name": c.Segment("name")})
	})

	v1 := r.Group("/v1")
	{
		v1.GET("/", func(c *tiny.Context) {
			c.HTML(http.StatusOK, "<h1>Hello Tiny</h1>")
		})

		v1.GET("/hello", func(c *tiny.Context) {
			// expect /hello?name=geektutu
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
		})
	}
	v2 := r.Group("/v2")
	v2.Use(onlyForV2())
	{
		v2.GET("/hello/:name", func(c *tiny.Context) {
			// expect /hello/geektutu
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Segment("name"), c.Path)
		})
		v2.POST("/login", func(c *tiny.Context) {
			c.JSON(http.StatusOK, tiny.H{
				"username": c.Post("username"),
				"password": c.Post("password"),
			})
		})

	}

	r.Static("/assets", "./tiny")

	r.GET("/panic", func(c *tiny.Context) {
		//fmt.Println("3")
		//panic("eror is error")
		names := []string{"geektutu"}
		c.String(http.StatusOK, names[100])
	})

	r.GET("/re1/{id:\\d+}", func(c *tiny.Context) {
		id := c.Segment("id")
		c.String(http.StatusOK, "re1 id: %s", id)
	})

	r.GET("/re2/{id:[a-z]+}", func(c *tiny.Context) {
		id := c.Segment("id")
		c.String(http.StatusOK, "re2 id: %s", id)
	})

	r.GET("/re3/{year:[12][0-9]{3}}/{month:[1-9]{2}}/{day:[1-9]{2}}/{hour:(12|[3-9])}", func(c *tiny.Context) {
		year := c.Segment("year")
		month := c.Segment("month")
		day := c.Segment("day")
		hour := c.Segment("hour")
		c.String(http.StatusOK, "re3 year: %s, month: %s, day: %s, hour: %s", year, month, day, hour)
	})
	
	srv := &http.Server{
		Addr:    ":9999",
		Handler: r,
	}

	go func() {
		log.Println("Server Start @", ":9999")
		err := srv.ListenAndServe()
		fmt.Println(err)
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server Start Error: %s\n", err)
		}
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown Error:", err)
	}
	log.Println("Server Shutdown")
}

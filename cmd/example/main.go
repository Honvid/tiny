package main

import (
	"context"
	"fmt"
	"honvid/cmd/example/model"
	"honvid/pkg/orm"
	"honvid/pkg/tiny"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func onlyForV2() tiny.Handler {
	return func(c *tiny.Context) {
		t := time.Now()
		c.Fail(500, "Internal Server Error")
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Request.RequestURI, time.Since(t))
	}
}

func main() {

	r := tiny.Default()

	r.GET("/", func(c *tiny.Context) {
		c.HTML(http.StatusOK, "<h1>Hello world!</h1>")
	})

	engine, _ := orm.New("mysql", "root:123456@tcp(127.0.0.1:3306)/demo")
	//defer engine.Close()
	r.GET("/mysql", func(c *tiny.Context) {
		s := engine.NewSession().Model(&model.User{})
		_ = s.DropTable()
		_ = s.CreateTable()
		if !s.HasTable(s.RefTable().Name) {
			c.HTML(http.StatusOK, "failed to create table")
		}
		//_, _ = s.Raw("DROP TABLE IF EXISTS User;").Exec()
		//_, _ = s.Raw("CREATE TABLE User(Name text);").Exec()
		//_, _ = s.Raw("CREATE TABLE User(Name text);").Exec()
		//result, _ := s.Raw("INSERT INTO User(`Name`) values (?), (?)", "Tom", "Sam").Exec()
		//count, _ := result.RowsAffected()
		//var TestDial, _ = dialect.GetDialect("mysql")
		//sa := schema.Parse(&model.User{}, TestDial)
		//if sa.Name != "User" || len(sa.Fields) != 2 {
		//	c.HTML(http.StatusOK, "failed to parse User struct")
		//}
		//c.HTML(http.StatusOK, fmt.Sprint("D--", sa.GetField("Name").Tag))
		//if sa.GetField("Name").Tag != "PRIMARY KEY" {
		//	c.HTML(http.StatusOK, fmt.Sprint(sa.GetField("Name").Type))
		//	//c.HTML(http.StatusOK, "failed to parse primary key")
		//}
		//result := s.Raw("SELECT * FROM User WHERE name = ?", "Tom").QueryRow()
		//c.HTML(http.StatusOK, fmt.Sprint(result))

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
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
		})
	}
	v2 := r.Group("/v2")
	v2.Use(onlyForV2())
	{
		v2.GET("/hello/:name", func(c *tiny.Context) {
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
		Addr:    ":9099",
		Handler: r,
	}

	go func() {
		log.Println("Server Start @", ":9099")
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

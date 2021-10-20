package main

import (
	"fmt"
	"gee"
	"log"
	"net/http"
	"time"
)



func main() {
	r := gee.Default()
	r.GET("/", func(c *gee.Context) {
		c.String(http.StatusOK, "Hello Geektutu\n")
	})
	// index out of range for testing Recovery()
	r.GET("/panic", func(c *gee.Context) {
		names := []string{"geektutu"}
		c.String(http.StatusOK, names[100])
	})

	r.Run(":9999")
}

type student struct {
	Name string
	Age  int8
}

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func onlyForV2() gee.HandlerFunc {
	return func(c *gee.Context) {
		// Start timer
		t := time.Now()
		// if a server error occurred
		c.Fail(500, "Internal Server Error")
		// Calculate resolution time
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}


//func main() {
//	r := gee.New()
//	r.Use(gee.Logger()) // global midlleware
//	r.GET("/", func(c *gee.Context) {
//		c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
//	})
//
//	v2 := r.Group("/v2")
//	v2.Use(onlyForV2()) // v2 group middleware
//	{
//		v2.GET("/hello/:name", func(c *gee.Context) {
//			// expect /hello/geektutu
//			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
//		})
//	}
//
//	r.Run(":9999")
//}




//func main() {
//	r := gee.New()
//	r.GET("/index", func(c *gee.Context) {
//		c.HTML(http.StatusOK, "<h1>Index Page</h1>")
//	})
//	v1 := r.Group("/v1")
//	{
//		v1.GET("/", func(c *gee.Context) {
//			c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
//		})
//
//		v1.GET("/hello", func(c *gee.Context) {
//			// expect /hello?name=geektutu
//			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
//		})
//	}
//	v2 := r.Group("/v2")
//	{
//		v2.GET("/hello/:name", func(c *gee.Context) {
//			// expect /hello/geektutu
//			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
//		})
//		v2.POST("/login", func(c *gee.Context) {
//			c.JSON(http.StatusOK, gee.H{
//				"username": c.PostForm("username"),
//				"password": c.PostForm("password"),
//			})
//		})
//
//	}
//
//	r.Run(":9999")
//}

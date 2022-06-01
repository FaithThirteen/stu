package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	for true {
		requestWithNoClose()
		time.Sleep(time.Microsecond * 100)
	}
}

var client2  = http.Client{}
func requestWithNoClose() {


	_, err := client2.Get("https://www.baidu.com")
	if err != nil {
		fmt.Printf("error occurred while fetching page, error: %s", err.Error())
	}
	fmt.Println("ok")
}
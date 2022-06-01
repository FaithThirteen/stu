package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	for true {
		requestWithClose()
		time.Sleep(time.Microsecond * 10)
	}
}

var client1  = http.Client{}
func requestWithClose() {
	resp, err := client1.Get("https://www.baidu.com")
	if err != nil {
		fmt.Printf("error occurred while fetching page, error: %s", err.Error())
		return
	}
	defer resp.Body.Close()
	fmt.Println("ok")
}
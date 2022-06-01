package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	totalSentRequests := &sync.WaitGroup{}
	allRequestsBackCh := make(chan struct{})
	multiplexCh := make(chan struct {
		result string
		retry  int
	})
	go func() {
		//所有请求完成之后会close掉allRequestsBackCh
		totalSentRequests.Wait()
		close(allRequestsBackCh)
	}()
	for i := 1; i <= 10; i++ {
		totalSentRequests.Add(1)
		go func() {
			// 标记已经执行完
			defer totalSentRequests.Done()
			// 模拟耗时操作
			time.Sleep(500 * time.Microsecond)
			// 模拟处理成功
			if random.Intn(500)%2 == 0 {
				multiplexCh <- struct {
					result string
					retry  int
				}{"finsh success", i}
			}
			// 处理失败不关心，当然，也可以加入一个错误的channel中进一步处理
		}()
	}
	select {
	case <-multiplexCh:
		fmt.Println("finish success")
	case <-allRequestsBackCh:
		// 到这里，说明全部的 goroutine 都执行完毕，但是都请求失败了
		fmt.Println("all req finish，but all fail")
	}
}
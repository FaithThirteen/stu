package main

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"
)

func retryHedged(req *http.Request, maxRetries int, timeout time.Duration,
	backoffStrategy BackoffStrategy) (*http.Response, error) {
	var (
		originalBody []byte
		err          error
	)
	if req != nil && req.Body != nil {
		originalBody, err = copyBody(req.Body)
	}
	if err != nil {
		return nil, err
	}

	AttemptLimit := maxRetries
	if AttemptLimit <= 0 {
		AttemptLimit = 1
	}

	client := http.Client{
		Timeout: timeout,
	}

	// 每次请求copy新的request
	copyRequest := func() (request *http.Request) {
		request = req.Clone(req.Context())
		if request.Body != nil {
			resetBody(request, originalBody)
		}
		return
	}

	multiplexCh := make(chan struct {
		resp  *http.Response
		err   error
		retry int
	})

	totalSentRequests := &sync.WaitGroup{}
	allRequestsBackCh := make(chan struct{})
	go func() {
		totalSentRequests.Wait()
		close(allRequestsBackCh)
	}()
	var resp *http.Response
	for i := 1; i <= AttemptLimit; i++ {
		totalSentRequests.Add(1)
		go func() {
			// 标记已经执行完
			defer totalSentRequests.Done()
			req = copyRequest()
			resp, err = client.Do(req)
			if err != nil {
				fmt.Printf("error sending the first time: %v\n", err)
			}
			// 重试 500 以上的错误码
			if err == nil && resp.StatusCode < 500 {
				multiplexCh <- struct {
					resp  *http.Response
					err   error
					retry int
				}{resp: resp, err: err, retry: i}
				return
			}
			// 如果正在重试，那么释放fd
			if resp != nil {
				resp.Body.Close()
			}
			// 重置body
			if req.Body != nil {
				resetBody(req, originalBody)
			}
			time.Sleep(backoffStrategy(i) + 1*time.Microsecond)
		}()
	}

	select {
	case res := <-multiplexCh:
		return res.resp, res.err
	case <-allRequestsBackCh:
		// 到这里，说明全部的 goroutine 都执行完毕，但是都请求失败了
		return nil, errors.New("all req finish，but all fail")
	}
}
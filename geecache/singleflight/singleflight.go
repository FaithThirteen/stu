package singleflight


import "sync"

type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

type Group struct {
	mu sync.Mutex       // protects m
	m  map[string]*call
}

// Do 并发请求条件相同时无论 Do 被调用多少次，函数 fn 都只会被调用一次
// https://pkg.go.dev/golang.org/x/sync/singleflight 标准库的singleflight实现
func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	// 加锁防止进入时其他请求添加key
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	// map中含有说明请求进行中，需要解锁等待执行完成然后取上个请求填入的值
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}

	// 执行中，添加到 g.m，表明 key 已经有对应的请求在处理
	// 通过sync.WaitGroup实现让其他请求等待返回
	c := new(call)
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

	// 获取值与错误
	c.val, c.err = fn()
	c.wg.Done()

	// 获取完成删除key，更新 g.m
	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()

	return c.val, c.err
}

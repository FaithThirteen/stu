package lru

import "container/list"

type Cache struct {
	// 最大可使用内存
	maxByte int64
	// 当前使用内存
	curByte int64
	// 真实存储值的列表
	ll *list.List
	// 保存key与链表节点的映射关系
	cache     map[string]*list.Element
	OnEvicted func(key string, value Value)
}

type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}

func New(maxByte int64, onEvicted func(key string, value Value)) *Cache {
	return &Cache{
		maxByte:   maxByte,
		curByte:   0,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// Get 查找
func (c *Cache) Get(key string) (value Value, ok bool) {

	if e, ok := c.cache[key]; ok {
		c.ll.MoveToFront(e)
		v := e.Value.(*entry)
		return v.value, true
	}

	return nil, false
}

// RemoveOldest 缓存淘汰
func (c *Cache) RemoveOldest() {

	e := c.ll.Back()
	if e != nil {
		kv := e.Value.(*entry)
		delete(c.cache, kv.key)
		c.ll.Remove(e)
		c.curByte -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// Add 添加缓存
func (c *Cache) Add(key string, value Value) {

	// 如果存在，更新缓存
	if e, ok := c.cache[key]; ok {
		c.ll.MoveToFront(e)
		kv := e.Value.(*entry)
		e.Value = kv.value
		c.curByte += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	}else {
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.curByte += int64(len(key)) + int64(value.Len())
	}

	// 判断是否需要移除更新缓存
	for c.maxByte != 0 && c.maxByte < c.curByte {
		c.RemoveOldest()
	}

}

func (c *Cache) Len() int {
	return c.ll.Len()
}
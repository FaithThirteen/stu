package test

import (
	"container/list"
	"errors"
)

type lru struct {
	// 最大容量
	maxCap int

	// 是否开启惰性删除，即取值时删除
	isLazyDel bool

	// 数据列表
	dlist *list.List

	// 列表映射关系
	cacheMap map[string]*list.Element
}

type entry struct {
	key   string
	value interface{}
}

//type lruconf struct {
//	maxCap         int
//	expirationTime int
//	isLazyDel      bool
//}

// DefaultLru 创建默认的lru
func DefaultLru() *lru {
	return &lru{
		maxCap:    4,
		dlist:     &list.List{},
		cacheMap:  make(map[string]*list.Element),
		isLazyDel: true,
	}
}

func NewLruCache() *lru {
	return &lru{
		maxCap:    20,
		dlist:     &list.List{},
		cacheMap:  make(map[string]*list.Element),
		isLazyDel: true,
	}
}

func (l *lru) Set(key string, value interface{}) error {
	if l.dlist == nil {
		return errors.New("lru结构体未初始化")
	}

	// 元素存在，放到队首
	if e, ok := l.cacheMap[key]; ok {
		l.dlist.MoveToFront(e)
		// 更新值
		e.Value.(*entry).value = value
		return nil
	}

	newEle := l.dlist.PushFront(&entry{key: key, value: value})
	l.cacheMap[key] = newEle

	// 超出缓存限制
	if l.dlist.Len() > l.maxCap {
		lastEle := l.dlist.Back()
		if lastEle == nil {
			return nil
		}
		delete(l.cacheMap, lastEle.Value.(*entry).key)
		l.dlist.Remove(lastEle)
	}

	return nil
}

func (l *lru) Get(key string) (interface{}, error) {

	if l.cacheMap == nil {
		return nil, errors.New("lru结构体未初始化")
	}

	if e, ok := l.cacheMap[key]; ok {
		l.dlist.PushFront(e)

		return e.Value.(*entry).value, nil
	}

	return nil, nil
}

func (l *lru) Remove(key string) bool {

	if l.cacheMap == nil {
		return false
	}

	if e, ok := l.cacheMap[key]; ok {
		delete(l.cacheMap, e.Value.(*entry).key)
		l.dlist.Remove(e)
		return true
	}

	return false
}

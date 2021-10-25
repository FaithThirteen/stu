package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// Hash 允许用于替换成自定义的 Hash 函数，也方便测试时替换，默认为 crc32.ChecksumIEEE 算法
type Hash func(data []byte) uint32

// Map 一致性哈希结构，包含虚拟节点倍数，哈希环，虚拟与真实节点映射表
type Map struct {
	hash     Hash           // 自定义的hash函数，默认为ChecksumIEEE
	replicas int            // 虚拟节点倍数
	keys     []int          // 哈希环
	hashMap  map[int]string // 键是虚拟节点的值，值是真实节点的名称
}

// New 创建Map实例
func New(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// Add 添加真实节点
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		// 对每个真实节点进行虚拟节点创建
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key))) // 生成虚拟节点的hash值
			m.keys = append(m.keys, hash)                      // 插入hash环
			m.hashMap[hash] = key                              // 映射虚拟节点和真实节点
		}
	}
	// 排序，便于后续的查找
	sort.Ints(m.keys)
}

// Get 获取key的的节点，search获取一个(0,n]范围的idx，用来确保key获取顺时针遇到的第一个节点，
// 当idx==len(m.keys)应该选择m.keys[0]
func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}

	hash := int(m.hash([]byte(key)))
	// 二分查找最近节点，hash大于最大值时会返回n
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})

	// 返回虚拟节点映射的真实节点
	return m.hashMap[m.keys[idx%len(m.keys)]]
}

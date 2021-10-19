package gee

type node struct {
	pattern  string  // 一个完整的URL，没有则为空
	part     string  // 通过 '/' 分割的节点的路由，比如/abc/123，abc和123就是2个part
	children []*node // 子节点
	isWild   bool    // 是否模糊匹配，比如:filename或*filename这样的node就为true
}

func (n *node) matchChild(part string) *node {

	for _, child := range n.children {
		if child.part == part || n.isWild {
			return child
		}
	}
	return nil
}

// @params pattern string "完整的URL"
// @params parts []string " '/'分割后的URL数组"
// @params height int "parts的下标，用于parts取值"
func (n *node) insert(pattern string, parts []string, height int) {
	// 递归终止条件
	// 如果已经匹配完了，那么将pattern赋值给该node，表示它是一个完整的url
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	// 取出路由，查询是否包含此节点
	part := parts[height]
	child := n.matchChild(part)

	// 未找到节点，初始化
	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	n.insert(pattern, parts, height+1)
}

// 这个函数跟matchChild有点像，但它是返回所有匹配的子节点，原因是它的场景是用以查找
// 它必须返回所有可能的子节点来进行遍历查找
// @params part string "路由子节点"
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

func (n *node) search(parts []string, height int) *node {
	if len(parts) == height {
		if n.pattern == "" {
			return nil
		}
	}

	part := parts[height]
	// 获取所有可能的子路径
	children := n.matchChildren(part)

	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}
	return nil
}

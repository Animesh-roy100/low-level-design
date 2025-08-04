package main

// lruNode for LRU doubly linked list.
type lruNode struct {
	key  string
	prev *lruNode
	next *lruNode
}

// LRUPolicy implements LRU eviction.
type LRUPolicy struct {
	head  *lruNode
	tail  *lruNode
	nodes map[string]*lruNode
}

// NewLRUPolicy creates a new LRU policy.
func NewLRUPolicy() *LRUPolicy {
	head := &lruNode{}
	tail := &lruNode{}
	head.next = tail
	tail.prev = head
	return &LRUPolicy{
		head:  head,
		tail:  tail,
		nodes: make(map[string]*lruNode),
	}
}

func (p *LRUPolicy) moveToFront(node *lruNode) {
	node.prev.next = node.next
	node.next.prev = node.prev
	node.next = p.head.next
	node.prev = p.head
	p.head.next.prev = node
	p.head.next = node
}

func (p *LRUPolicy) Access(key string) {
	if node, ok := p.nodes[key]; ok {
		p.moveToFront(node)
	}
}

func (p *LRUPolicy) Add(key string) {
	node := &lruNode{key: key}
	p.nodes[key] = node
	node.next = p.head.next
	node.prev = p.head
	p.head.next.prev = node
	p.head.next = node
}

func (p *LRUPolicy) Evict() string {
	if p.tail.prev == p.head {
		return ""
	}
	node := p.tail.prev
	node.prev.next = p.tail
	p.tail.prev = node.prev
	delete(p.nodes, node.key)
	return node.key
}

func (p *LRUPolicy) Remove(key string) {
	if node, ok := p.nodes[key]; ok {
		node.prev.next = node.next
		node.next.prev = node.prev
		delete(p.nodes, key)
	}
}

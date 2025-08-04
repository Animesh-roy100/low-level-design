package main

// lfuNode for LFU doubly linked list.
type lfuNode struct {
	key  string
	freq int
	prev *lfuNode
	next *lfuNode
}

// freqList is a DLL for a specific frequency.
type freqList struct {
	head *lfuNode
	tail *lfuNode
}

// LFUPolicy implements LFU eviction.
type LFUPolicy struct {
	minFreq int
	freqMap map[int]*freqList
	nodeMap map[string]*lfuNode
}

// NewLFUPolicy creates a new LFU policy.
func NewLFUPolicy() *LFUPolicy {
	return &LFUPolicy{
		minFreq: 0,
		freqMap: make(map[int]*freqList),
		nodeMap: make(map[string]*lfuNode),
	}
}

func (p *LFUPolicy) getList(freq int) *freqList {
	if _, ok := p.freqMap[freq]; !ok {
		head := &lfuNode{}
		tail := &lfuNode{}
		head.next = tail
		tail.prev = head
		p.freqMap[freq] = &freqList{head: head, tail: tail}
	}
	return p.freqMap[freq]
}

func (p *LFUPolicy) moveToNewFreq(node *lfuNode, newFreq int) {
	// Remove from old list
	node.prev.next = node.next
	node.next.prev = node.prev
	oldList := p.freqMap[node.freq]
	if oldList.head.next == oldList.tail {
		delete(p.freqMap, node.freq)
	}
	node.freq = newFreq
	newList := p.getList(newFreq)
	// Add to front
	node.next = newList.head.next
	node.prev = newList.head
	newList.head.next.prev = node
	newList.head.next = node
}

func (p *LFUPolicy) Access(key string) {
	node, ok := p.nodeMap[key]
	if !ok {
		return
	}
	oldFreq := node.freq
	p.moveToNewFreq(node, node.freq+1)
	// Update minFreq if old freq was min and now empty
	if oldFreq == p.minFreq {
		if _, exists := p.freqMap[oldFreq]; !exists {
			p.minFreq++
		}
	}
}

func (p *LFUPolicy) Add(key string) {
	node := &lfuNode{key: key, freq: 1}
	p.nodeMap[key] = node
	list := p.getList(1)
	node.next = list.head.next
	node.prev = list.head
	list.head.next.prev = node
	list.head.next = node
	p.minFreq = 1
}

func (p *LFUPolicy) Evict() string {
	if p.minFreq == 0 {
		return ""
	}
	list, ok := p.freqMap[p.minFreq]
	if !ok || list.head.next == list.tail {
		return ""
	}
	node := list.tail.prev
	node.prev.next = node.next
	node.next.prev = node.prev
	if list.head.next == list.tail {
		delete(p.freqMap, p.minFreq)
		// Find new minFreq
		p.minFreq = 0
		for f := range p.freqMap {
			if p.minFreq == 0 || f < p.minFreq {
				p.minFreq = f
			}
		}
	}
	delete(p.nodeMap, node.key)
	return node.key
}

func (p *LFUPolicy) Remove(key string) {
	node, ok := p.nodeMap[key]
	if !ok {
		return
	}
	node.prev.next = node.next
	node.next.prev = node.prev
	list := p.freqMap[node.freq]
	if list.head.next == list.tail {
		delete(p.freqMap, node.freq)
		if node.freq == p.minFreq {
			// Find new minFreq
			p.minFreq = 0
			for f := range p.freqMap {
				if p.minFreq == 0 || f < p.minFreq {
					p.minFreq = f
				}
			}
		}
	}
	delete(p.nodeMap, key)
}

package main

import (
	"fmt"
	"time"
)

func main() {
	// LRU example
	lruCache := NewCache(2, NewLRUPolicy())
	lruCache.Set("key1", "value1", 0)
	lruCache.Set("key2", "value2", time.Second*10) // with TTL
	val, found := lruCache.Get("key1")             // Access moves to recent
	fmt.Println(val, found)                        // "value1" true

	// LFU example
	lfuCache := NewCache(2, NewLFUPolicy())
	lfuCache.Set("keyA", "valueA", 0)
	lfuCache.Get("keyA") // Increase frequency
	lfuCache.Set("keyB", "valueB", 0)
	lfuCache.Set("keyC", "valueC", 0) // Evicts least frequent (keyB)
}

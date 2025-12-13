package main

/*

Objects of a superclass should be replaceable with objects of its subclasses without affecting the correctness of the program.

This means a subclass should be able to be substituted for its parent class without breaking the application

Liskov Substitution Principle states that subtypes must be completely substitutable for their base types without changing program behavior or correctness.

Before - Violates LSP

type Rectangle struct {
	Width, Height float64
}

func (r *Rectangle) SetWidth(w float64) {
	r.Width = w
}

func (r *Rectangle) SetHeight(h float64) {
	r.Height = h
}

func (r *Rectangle) Area() float64 {
	return r.Width * r.Height
}

type Square struct {
	Rectangle
}

func (s *Square) SetWidth(w float64) {
	s.Width = w
	s.Height = w // Violates LSP: changing width also changes height
}

func (s *Square) SetHeight(h float64) {
	s.Height = h
	s.Width = h // Violates LSP: changing height also changes width
}

func main() {
	var r *Rectangle = &Square{}
	r.SetWidth(5)
	r.SetHeight(4)

	fmt.Println(Area(r)) // Expected 20, gets 16
}

*/

type Cache interface {
	Get(key string) string
}

type InMemoryCache struct {
	store map[string]string
}

func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{store: make(map[string]string)}
}

func (c *InMemoryCache) Get(key string) string {
	return c.store[key]
}

type RedisCache struct {
	store map[string]string
}

func NewRedisCache() *RedisCache {
	return &RedisCache{store: make(map[string]string)}
}

func (c *RedisCache) Get(key string) string {
	return c.store[key]
}

// func main() {
// 	var cache Cache
// 	cache = NewInMemoryCache()
// 	cache.Get("example_key")

// 	cache = NewRedisCache()
// 	cache.Get("example_key")
// }

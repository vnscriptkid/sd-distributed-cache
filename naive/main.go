package main

import (
	"fmt"
	"sync"
)

type Cache struct {
	data map[string]interface{}
	mu   sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{
		data: make(map[string]interface{}),
	}
}

func (c *Cache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, exists := c.data[key]
	return value, exists
}

func main() {
	cache := NewCache()

	// string
	cache.Set("example", "value")
	if val, found := cache.Get("example"); found {
		fmt.Printf("%T\n", val) // Output: string
		fmt.Println(val)        // Output: value
	}

	fmt.Println()

	// int
	cache.Set("test", 123)
	if val, found := cache.Get("test"); found {
		// print type of val
		fmt.Printf("%T\n", val) // Output: int
		fmt.Println(val)        // Output: 123
	}
}

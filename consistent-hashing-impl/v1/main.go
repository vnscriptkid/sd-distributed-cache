package main

import (
	"fmt"
	"hash/crc32"
	"sort"
	"sync"
)

// ConsistentHash struct to hold the hash ring and other related information
type ConsistentHash struct {
	sync.RWMutex
	keys    []int          // sorted hash ring
	hashMap map[int]string // mapping from hash to node
}

// NewConsistentHash creates a new ConsistentHash object
func NewConsistentHash() *ConsistentHash {
	return &ConsistentHash{
		hashMap: make(map[int]string),
	}
}

// Hash function to compute the hash of a key
func (ch *ConsistentHash) Hash(key string) int {
	return int(crc32.ChecksumIEEE([]byte(key)))
}

// Add a node to the hash ring
func (ch *ConsistentHash) Add(nodes ...string) {
	ch.Lock()
	defer ch.Unlock()
	for _, node := range nodes {
		hash := ch.Hash(node)
		ch.keys = append(ch.keys, hash)
		ch.hashMap[hash] = node
	}
	// Sort the keys to make the binary search work
	sort.Ints(ch.keys)
}

// Remove a node from the hash ring
func (ch *ConsistentHash) Remove(node string) {
	ch.Lock()
	defer ch.Unlock()
	hash := ch.Hash(node)
	for i, key := range ch.keys {
		if key == hash {
			ch.keys = append(ch.keys[:i], ch.keys[i+1:]...)
			break
		}
	}
	delete(ch.hashMap, hash)
}

// Get retrieves the closest node in the hash ring for the given key
func (ch *ConsistentHash) Get(key string) string {
	ch.RLock()
	defer ch.RUnlock()
	if len(ch.keys) == 0 {
		return ""
	}
	hash := ch.Hash(key)
	// Binary search for the appropriate node
	// idx is the index of the first node whose hash is greater than or equal to the key's hash
	idx := sort.Search(len(ch.keys), func(i int) bool {
		return ch.keys[i] >= hash
	})
	// If idx is equal to len(keys), wrap around to the first node
	if idx == len(ch.keys) {
		idx = 0
	}
	return ch.hashMap[ch.keys[idx]]
}

func main() {
	// Create a new consistent hash ring
	ch := NewConsistentHash()

	// Add nodes to the ring
	ch.Add("NodeA", "NodeB", "NodeC")

	// Test with some keys
	keys := []string{"my-key1", "my-key2", "my-key3", "my-key4"}

	for _, key := range keys {
		node := ch.Get(key)
		fmt.Printf("Key %s is assigned to node %s\n", key, node)
	}

	// Remove a node and test again
	ch.Remove("NodeB")
	fmt.Println("\nAfter removing NodeB:")
	for _, key := range keys {
		node := ch.Get(key)
		fmt.Printf("Key %s is assigned to node %s\n", key, node)
	}
}

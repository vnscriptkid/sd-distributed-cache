package main

import (
	"fmt"
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
)

// ConsistentHash struct to hold the hash ring and other related information
type ConsistentHash struct {
	sync.RWMutex
	vNodesCount int            // number of virtual nodes (vNodesCount) per node
	keys        []int          // sorted hash ring
	hashToNode  map[int]string // mapping from hash to node
}

// NewConsistentHash creates a new ConsistentHash object
func NewConsistentHash(vNodesCount int) *ConsistentHash {
	return &ConsistentHash{
		vNodesCount: vNodesCount,
		hashToNode:  make(map[int]string),
	}
}

// Hash function to compute the hash of a key
// Values range from 0 to 2^32-1
func (ch *ConsistentHash) Hash(key string) int {
	return int(crc32.ChecksumIEEE([]byte(key)))
}

// Add a node to the hash ring
func (ch *ConsistentHash) Add(nodes ...string) {
	ch.Lock()
	defer ch.Unlock()
	for _, node := range nodes {
		for i := 0; i < ch.vNodesCount; i++ {
			hash := ch.Hash(node + strconv.Itoa(i))
			ch.keys = append(ch.keys, hash)
			ch.hashToNode[hash] = node
		}
	}
	sort.Ints(ch.keys)

	// Calculate and print the ranges for each node
	ch.printNodeRanges()
}

type nodeRange struct {
	from      int
	to        int
	isBoundry bool
}

func (nr nodeRange) countValues() int64 {
	if nr.isBoundry {
		return int64(1<<32-nr.from) + int64(nr.to) + 1
	}

	return int64(nr.to - nr.from + 1)
}

func (ch *ConsistentHash) printNodeRanges() {
	// printNodeRanges calculates and prints the ranges of hash values for each node
	// [  5 10, 20, 30, 40, 50, 60  ]

	// hash(nodeA1)=20 -> [11, 20]
	// hash(nodeA2)=50 -> [41, 50]
	// nodeA: [11, 20] [41, 50]
	// totalValues: 20-11+1 + 50-41+1 = 10 + 10 = 20
	// percentage: 20/60 = 33.33%
	nodeToNodeRanges := make(map[string][]nodeRange)
	for i, key := range ch.keys {
		node := ch.hashToNode[key]
		if _, ok := nodeToNodeRanges[node]; !ok {
			nodeToNodeRanges[node] = []nodeRange{}
		}

		if i == 0 {
			nodeToNodeRanges[node] = append(nodeToNodeRanges[node], nodeRange{from: ch.keys[len(ch.keys)-1], to: key, isBoundry: true})
		} else {
			nodeToNodeRanges[node] = append(nodeToNodeRanges[node], nodeRange{from: ch.keys[i-1], to: key})
		}
	}

	// Print the ranges for each node
	for node, ranges := range nodeToNodeRanges {
		values := int64(0)
		for _, r := range ranges {
			values += r.countValues()
		}
		percentage := float64(values) / float64(1<<32) * 100
		fmt.Printf("Node %s accounts for %.2f%% of the hash ring\n", node, percentage)
	}
}

// Get retrieves the closest node in the hash ring for the given key
func (ch *ConsistentHash) Get(key string) string {
	ch.RLock()
	defer ch.RUnlock()
	if len(ch.keys) == 0 {
		return ""
	}
	hash := ch.Hash(key)
	// Binary search for the appropriate replica
	idx := sort.Search(len(ch.keys), func(i int) bool {
		return ch.keys[i] >= hash
	})
	// If idx is equal to len(keys), wrap around to the first node
	if idx == len(ch.keys) {
		idx = 0
	}
	return ch.hashToNode[ch.keys[idx]]
}

func main() {
	// Create a new consistent hash ring with 3 vNodesCount per node
	ch := NewConsistentHash(100)

	// Add nodes to the ring
	ch.Add("NodeA", "NodeB", "NodeC")

	// Test with some keys
	keys := []string{"my-key1", "my-key2", "my-key3", "my-key4"}

	for _, key := range keys {
		node := ch.Get(key)
		fmt.Printf("Key %s is assigned to node %s\n", key, node)
	}
}

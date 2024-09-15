package main

import (
	"fmt"
)

// Node represents an element in the linked list
type Node struct {
	key   string
	value interface{}
	next  *Node
}

// LinkedList represents a linked list
type LinkedList struct {
	head *Node
}

// HashTable represents the hash table structure
type HashTable struct {
	table []*LinkedList
	size  int
}

// NewHashTable initializes a new hash table
func NewHashTable(size int) *HashTable {
	table := make([]*LinkedList, size)
	for i := range table {
		table[i] = &LinkedList{}
	}
	return &HashTable{
		table: table,
		size:  size,
	}
}

// hash function maps a key to an index
func (ht *HashTable) hash(key string) int {
	hash := 0
	for _, char := range key {
		hash += int(char)
	}
	return hash % ht.size
}

// Put inserts a key-value pair into the hash table
func (ht *HashTable) Put(key string, value interface{}) {
	index := ht.hash(key)
	list := ht.table[index]

	// Check if the key already exists, and update the value if it does
	for node := list.head; node != nil; node = node.next {
		if node.key == key {
			node.value = value
			return
		}
	}

	// Insert new node at the beginning of the linked list
	newNode := &Node{key: key, value: value, next: list.head}
	list.head = newNode
}

// Get retrieves the value for a given key from the hash table
func (ht *HashTable) Get(key string) (interface{}, bool) {
	index := ht.hash(key)
	list := ht.table[index]

	for node := list.head; node != nil; node = node.next {
		if node.key == key {
			return node.value, true
		}
	}

	return nil, false
}

// Remove deletes a key-value pair from the hash table
func (ht *HashTable) Remove(key string) {
	index := ht.hash(key)
	list := ht.table[index]

	var prev *Node
	for node := list.head; node != nil; node = node.next {
		if node.key == key {
			if prev == nil {
				// Removing the head of the list
				list.head = node.next
			} else {
				// Removing a node from the middle or end
				prev.next = node.next
			}
			return
		}
		prev = node
	}
}

func main() {
	ht := NewHashTable(10)

	ht.Put("name", "John")
	ht.Put("age", 30)
	ht.Put("city", "New York")
	ht.Put("age", 31) // Update the value for key "age"

	name, found := ht.Get("name")
	if found {
		fmt.Println("Name:", name)
	}

	age, found := ht.Get("age")
	if found {
		fmt.Println("Age:", age)
	}

	ht.Remove("city")
	_, found = ht.Get("city")
	if !found {
		fmt.Println("City not found")
	}
}

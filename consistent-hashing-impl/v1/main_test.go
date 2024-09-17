package main

import (
	"testing"
)

func TestConsistentHash(t *testing.T) {
	ch := NewConsistentHash()
	ch.Add("NodeA", "NodeB", "NodeC")

	tests := []struct {
		key      string
		expected string
	}{
		{"my-key1", "NodeB"}, // Expected node for my-key1
		{"my-key2", "NodeC"}, // Expected node for my-key2
		{"my-key3", "NodeB"}, // Expected node for my-key3
		{"my-key4", "NodeB"}, // Expected node for my-key4
	}

	for _, test := range tests {
		node := ch.Get(test.key)
		if node != test.expected {
			t.Errorf("Key %s: expected %s, got %s", test.key, test.expected, node)
		}
	}

	// Test removing a node
	ch.Remove("NodeB")
	if node := ch.Get("my-key2"); node == "NodeB" {
		t.Error("Expected NodeB to be removed, but it was still returned")
	}
}

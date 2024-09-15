package main

import (
	"testing"
)

func TestHashTable(t *testing.T) {
	ht := NewHashTable(10)

	// Test Put and Get
	ht.Put("name", "John")
	if value, found := ht.Get("name"); !found || value != "John" {
		t.Errorf("Expected 'John', got %v", value)
	}

	ht.Put("age", 30)
	if value, found := ht.Get("age"); !found || value != 30 {
		t.Errorf("Expected 30, got %v", value)
	}

	// Test Update
	ht.Put("age", 31)
	if value, found := ht.Get("age"); !found || value != 31 {
		t.Errorf("Expected 31, got %v", value)
	}

	// Test Remove
	ht.Remove("name")
	if _, found := ht.Get("name"); found {
		t.Error("Expected 'name' to be removed")
	}

	// Test non-existent key
	if _, found := ht.Get("city"); found {
		t.Error("Expected 'city' to not be found")
	}
}

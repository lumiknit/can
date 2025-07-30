package jdb

import (
	"fmt"
	"os"
	"testing"
)

func TestJDB_BasicOperations(t *testing.T) {
	filename := "test_basic.jdb"
	defer os.Remove(filename)

	db, err := New(filename)
	if err != nil {
		t.Fatalf("Failed to create JDB: %v", err)
	}
	defer db.Close()

	// Test Set and Get
	err = db.Set("key1", "value1")
	if err != nil {
		t.Fatalf("Failed to set key1: %v", err)
	}

	value, exists := db.Get("key1")
	if !exists {
		t.Fatal("Key1 should exist")
	}
	if value != "value1" {
		t.Fatalf("Expected 'value1', got %v", value)
	}

	// Test Del
	err = db.Del("key1")
	if err != nil {
		t.Fatalf("Failed to delete key1: %v", err)
	}

	_, exists = db.Get("key1")
	if exists {
		t.Fatal("Key1 should not exist after deletion")
	}
}

func TestJDB_Persistence(t *testing.T) {
	filename := "test_persistence.jdb"
	defer os.Remove(filename)

	// Create and populate database
	db1, err := New(filename)
	if err != nil {
		t.Fatalf("Failed to create JDB: %v", err)
	}

	err = db1.Set("persistent_key", "persistent_value")
	if err != nil {
		t.Fatalf("Failed to set persistent_key: %v", err)
	}

	db1.Close()

	// Reopen and check persistence
	db2, err := New(filename)
	if err != nil {
		t.Fatalf("Failed to reopen JDB: %v", err)
	}
	defer db2.Close()

	value, exists := db2.Get("persistent_key")
	if !exists {
		t.Fatal("Persistent key should exist after reopening")
	}
	if value != "persistent_value" {
		t.Fatalf("Expected 'persistent_value', got %v", value)
	}
}

func TestJDB_ComplexValues(t *testing.T) {
	filename := "test_complex.jdb"
	defer os.Remove(filename)

	db, err := New(filename)
	if err != nil {
		t.Fatalf("Failed to create JDB: %v", err)
	}
	defer db.Close()

	// Test with map
	complexValue := map[string]any{
		"name":   "test",
		"count":  42,
		"active": true,
	}

	err = db.Set("complex", complexValue)
	if err != nil {
		t.Fatalf("Failed to set complex value: %v", err)
	}

	value, exists := db.Get("complex")
	if !exists {
		t.Fatal("Complex key should exist")
	}

	// Note: JSON unmarshaling converts numbers to float64
	if retrieved, ok := value.(map[string]any); ok {
		if retrieved["name"] != "test" {
			t.Fatalf("Expected name 'test', got %v", retrieved["name"])
		}
		if retrieved["count"] != float64(42) {
			t.Fatalf("Expected count 42, got %v", retrieved["count"])
		}
		if retrieved["active"] != true {
			t.Fatalf("Expected active true, got %v", retrieved["active"])
		}
	} else {
		t.Fatalf("Retrieved value is not a map: %T", value)
	}
}

func TestJDB_Compaction(t *testing.T) {
	filename := "test_compaction.jdb"
	defer os.Remove(filename)

	db, err := New(filename)
	if err != nil {
		t.Fatalf("Failed to create JDB: %v", err)
	}
	defer db.Close()

	// Add many entries to trigger compaction
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("key%d", i)
		value := fmt.Sprintf("value%d", i)
		err = db.Set(key, value)
		if err != nil {
			t.Fatalf("Failed to set %s: %v", key, err)
		}
	}

	// Verify all data is still accessible after potential compaction
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("key%d", i)
		expectedValue := fmt.Sprintf("value%d", i)

		value, exists := db.Get(key)
		if !exists {
			t.Fatalf("Key %s should exist", key)
		}
		if value != expectedValue {
			t.Fatalf("Expected %s, got %v", expectedValue, value)
		}
	}
}

func TestJDB_ConcurrentAccess(t *testing.T) {
	filename := "test_concurrent.jdb"
	defer os.Remove(filename)

	db, err := New(filename)
	if err != nil {
		t.Fatalf("Failed to create JDB: %v", err)
	}
	defer db.Close()

	// Run concurrent operations
	done := make(chan bool, 3)

	// Writer 1
	go func() {
		for i := 0; i < 50; i++ {
			key := fmt.Sprintf("writer1_%d", i)
			db.Set(key, i)
		}
		done <- true
	}()

	// Writer 2
	go func() {
		for i := 0; i < 50; i++ {
			key := fmt.Sprintf("writer2_%d", i)
			db.Set(key, i*2)
		}
		done <- true
	}()

	// Reader
	go func() {
		for i := 0; i < 100; i++ {
			key := fmt.Sprintf("writer1_%d", i%50)
			db.Get(key)
		}
		done <- true
	}()

	// Wait for all goroutines to complete
	for i := 0; i < 3; i++ {
		<-done
	}

	// Verify some data
	value, exists := db.Get("writer1_0")
	if exists && value != float64(0) {
		t.Fatalf("Expected 0, got %v", value)
	}

	value, exists = db.Get("writer2_1")
	if exists && value != float64(2) {
		t.Fatalf("Expected 2, got %v", value)
	}
}

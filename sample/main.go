package main

import (
	fastmap "github.com/billowdev/fastmap/hashmap"
)

func main() {
	// Example usage
	hashMap := fastmap.NewHashMap[string, int]()

	// Add some values
	hashMap.Put("one", 1)
	hashMap.Put("two", 2)
	hashMap.Put("three", 3)

	// Use callback function to print all entries
	hashMap.ForEach(func(key string, value int) {
		println(key, ":", value)
	})

	// Get a value
	if value, exists := hashMap.Get("two"); exists {
		println("Value for 'two':", value)
	}

	// Remove a value
	hashMap.Remove("one")
	println("Size after removal:", hashMap.Size())
}

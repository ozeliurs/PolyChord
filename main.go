package main

import (
	"fmt"
	"time"
)

func main() {
	network := NewNetwork()

	// Create the first node
	node1 := NewNode(1, network)
	go node1.RunStabilizer(5 * time.Second)

	// Create additional nodes and join the network
	for i := 2; i <= 5; i++ {
		node := NewNode(i, network)
		err := node.Join(node1)
		if err != nil {
			fmt.Printf("Error joining node %d: %v\n", i, err)
			continue
		}
		go node.RunStabilizer(5 * time.Second)
	}

	// Insert and retrieve data
	node1.Put("foo", "bar")
	time.Sleep(2 * time.Second) // Allow time for stabilization

	value, found := node1.Get("foo")
	if found {
		fmt.Printf("Found value: %s\n", value)
	} else {
		fmt.Println("Key not found")
	}

	// Simulate more keys and values
	node1.Put("baz", "qux")
	node1.Put("apple", "banana")

	// Retrieve the values from the network
	time.Sleep(2 * time.Second) // Allow time for stabilization
	value, found = node1.Get("baz")
	if found {
		fmt.Printf("Found value for 'baz': %s\n", value)
	} else {
		fmt.Println("Key 'baz' not found")
	}

	value, found = node1.Get("apple")
	if found {
		fmt.Printf("Found value for 'apple': %s\n", value)
	} else {
		fmt.Println("Key 'apple' not found")
	}
}

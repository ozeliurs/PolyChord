package main

import (
	"fmt"
	"time"
)

func main() {
	network := NewNetwork()
	go network.StartPeriodicNetworkInfoDump()

	// Create the first node with a specified ID
	node1 := NewNodeWithRandomID(network)

	numberOfNodes := 3

	// Create additional nodes with random IDs and join the network
	for i := 0; i < numberOfNodes; i++ {
		node := NewNodeWithRandomID(network)
		if node == nil {
			continue
		}
		err := node.Join(node1)
		if err != nil {
			fmt.Printf("Error joining node with ID %d: %v\n", node.ID, err)
			continue
		}
		// time.Sleep(100 * time.Millisecond)
	}

	time.Sleep(2 * time.Second)

	// Insert some data into the network
	node1.Put("foo", "bar")
	node1.Put("baz", "qux")
	node1.Put("apple", "banana")

	time.Sleep(2 * time.Second) // Allow time for stabilization

	jsonInfo, err := network.PrintNetworkInfoJSON()
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println(jsonInfo)
	}

	NewNodeWithRandomID(network)

	node1.Put("luigi", "liquori")
}

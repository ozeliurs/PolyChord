package main

import (
	"fmt"
	"time"
)

func main() {
	stabilizerInterval := 100 * time.Millisecond
	network := NewNetwork()

	saveInterval := 10 * time.Millisecond
	filename := "network_states.json"
	SaveNetworkState(network, saveInterval, filename)

	// Create the first node with a specified ID
	node1 := NewNodeWithRandomID(network)
	go node1.RunStabilizer(stabilizerInterval)

	// Create additional nodes with random IDs and join the network
	for i := 2; i <= 5; i++ {
		node := NewNodeWithRandomID(network)
		if node == nil {
			continue
		}
		err := node.Join(node1)
		if err != nil {
			fmt.Printf("Error joining node with ID %d: %v\n", node.ID, err)
			continue
		}
		go node.RunStabilizer(stabilizerInterval)
	}

	time.Sleep(2 * time.Second)

	// Insert some data into the network
	node1.Put("foo", "bar")
	node1.Put("baz", "qux")
	node1.Put("apple", "banana")

	time.Sleep(2 * time.Second) // Allow time for stabilization

	// Visualize the network and the hash tables on each node
	VisualizeNetwork(network)

	// Run the network for a specified duration
	runDuration := 30 * time.Second
	fmt.Printf("Running network for %s...\n", runDuration)
	time.Sleep(runDuration)

	fmt.Println("Program completed.")
}

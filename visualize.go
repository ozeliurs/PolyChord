package main

import (
	"fmt"
)

// Visualize the entire network and each node's data
func VisualizeNetwork(network *Network) {
	fmt.Println("---- Chord DHT Network Visualization ----")
	fmt.Printf("Number of nodes: %d\n\n", len(network.Nodes))

	// Iterate through all nodes in the network
	for _, node := range network.Nodes {
		fmt.Printf("Node ID: %d\n", node.ID)

		// Display the node's key-value pairs
		fmt.Println("Stored Key-Value Pairs:")
		if len(node.Data) > 0 {
			for key, value := range node.Data {
				fmt.Printf("  Hashed Key: %d, Value: %s\n", key, value)
			}
		} else {
			fmt.Println("  No data stored on this node.")
		}

		// Display the finger table for routing information
		fmt.Println("Finger Table:")
		for i, finger := range node.FingerTable {
			if finger != nil {
				fmt.Printf("  Entry %d -> Node ID: %d\n", i, finger.ID)
			} else {
				fmt.Printf("  Entry %d -> (nil)\n", i)
			}
		}

		// Display predecessor and successor information
		if node.Predecessor != nil {
			fmt.Printf("Predecessor: %d\n", node.Predecessor.ID)
		} else {
			fmt.Println("Predecessor: (nil)")
		}

		if node.Successor != nil {
			fmt.Printf("Successor: %d\n", node.Successor.ID)
		} else {
			fmt.Println("Successor: (nil)")
		}

		fmt.Println("----------------------------------------")
	}
}

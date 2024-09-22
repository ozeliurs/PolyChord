package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: <program name> [simple|stress_keys|network_growth|node_leave]")
		return
	}

	scenario := os.Args[1]

	switch scenario {
	case "simple":
		runSimpleScenario()
	case "stress_keys":
		runStressKeysScenario()
	case "network_growth":
		runNetworkGrowthScenario()
	case "node_leave":
		runNodeLeaveScenario()
	default:
		fmt.Println("Invalid scenario. Use 'simple', 'stress_keys', 'network_growth', or 'node_leave'")
	}
}

func runSimpleScenario() {
	network := NewNetwork(true)

	// Create the first node with a specified ID
	node1 := NewNodeWithRandomID(network)

	// Create two additional nodes and join the network
	node2 := NewNodeWithRandomID(network)
	node3 := NewNodeWithRandomID(network)

	err := node2.Join(node1)
	if err != nil {
		fmt.Printf("Error joining node with ID %d: %v\n", node2.ID, err)
	}

	err = node3.Join(node1)
	if err != nil {
		fmt.Printf("Error joining node with ID %d: %v\n", node3.ID, err)
	}

	time.Sleep(100 * time.Millisecond)

	// Insert 9 key-value pairs into the network
	node1.Put("key1", "value1")
	node1.Put("key2", "value2")
	node1.Put("key3", "value3")
	node2.Put("key4", "value4")
	node2.Put("key5", "value5")
	node2.Put("key6", "value6")
	node3.Put("key7", "value7")
	node3.Put("key8", "value8")
	node3.Put("key9", "value9")

	time.Sleep(100 * time.Millisecond)

	jsonInfo, err := network.PrintNetworkInfoJSON()
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println(jsonInfo)
	}

	network.Stop()
}

func runStressKeysScenario() {
	network := NewNetwork(true)

	// Create the first node with a specified ID
	node1 := NewNodeWithRandomID(network)

	numberOfNodes := 10

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
	}

	time.Sleep(100 * time.Millisecond)

	// Insert many keys into the network
	startTime := time.Now()
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("key%d", i)
		value := fmt.Sprintf("value%d", i)
		node1.Put(key, value)
	}
	elapsedTime := time.Since(startTime)
	fmt.Printf("Time taken to store 1000 keys: %v\n", elapsedTime)

	time.Sleep(1 * time.Second)

	jsonInfo, err := network.PrintNetworkInfoJSON()
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println(jsonInfo)
	}

	network.Stop()
}

func runNetworkGrowthScenario() {
	network := NewNetwork()

	// Create the first three nodes
	node1 := NewNodeWithRandomID(network)
	node2 := NewNodeWithRandomID(network)
	node3 := NewNodeWithRandomID(network)

	err := node2.Join(node1)
	if err != nil {
		fmt.Printf("Error joining node with ID %d: %v\n", node2.ID, err)
	}

	err = node3.Join(node1)
	if err != nil {
		fmt.Printf("Error joining node with ID %d: %v\n", node3.ID, err)
	}

	// Wait for the network to stabilize
	time.Sleep(1 * time.Second)

	// Insert some initial keys
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("initial_key%d", i)
		value := fmt.Sprintf("initial_value%d", i)
		node1.Put(key, value)
	}

	// Try getting keys before network growth
	for i := 0; i < 5; i++ {
		key := fmt.Sprintf("initial_key%d", i)
		value, found := node1.Get(key)
		if !found {
			fmt.Printf("Error getting key %s before network growth: %v\n", key)
		} else {
			fmt.Printf("Got key %s with value %s before network growth\n", key, value)
		}
	}

	// Add 25 nodes instantly
	for i := 0; i < 25; i++ {
		node := NewNodeWithRandomID(network)
		if node == nil {
			continue
		}
		err := node.Join(node1)
		if err != nil {
			fmt.Printf("Error joining node with ID %d: %v\n", node.ID, err)
			continue
		}
	}

	// Try getting keys during network growth
	for i := 0; i < 5; i++ {
		key := fmt.Sprintf("initial_key%d", i)
		value, found := node1.Get(key)
		if !found {
			fmt.Printf("Error getting key %s during network growth: %v\n", key)
		} else {
			fmt.Printf("Got key %s with value %s during network growth\n", key, value)
		}
	}

	// Wait for the network to stabilize
	time.Sleep(2 * time.Second)

	// Try getting keys after network stabilization
	for i := 0; i < 5; i++ {
		key := fmt.Sprintf("initial_key%d", i)
		value, found := node1.Get(key)
		if !found {
			fmt.Printf("Error getting key %s after network stabilization: %v\n", key)
		} else {
			fmt.Printf("Got key %s with value %s after network stabilization\n", key, value)
		}
	}

	jsonInfo, err := network.PrintNetworkInfoJSON()
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println(jsonInfo)
	}

	network.Stop()
}

func runNodeLeaveScenario() {
	network := NewNetwork(true)

	// Create three nodes
	node1 := NewNodeWithRandomID(network)
	node2 := NewNodeWithRandomID(network)
	node3 := NewNodeWithRandomID(network)

	// Join the network
	err := node2.Join(node1)
	if err != nil {
		fmt.Printf("Error joining node with ID %d: %v\n", node2.ID, err)
	}

	err = node3.Join(node1)
	if err != nil {
		fmt.Printf("Error joining node with ID %d: %v\n", node3.ID, err)
	}

	time.Sleep(100 * time.Millisecond)

	// Insert 9 random key-value pairs
	for i := 0; i < 9; i++ {
		key := fmt.Sprintf("key%d", i)
		value := fmt.Sprintf("value%d", i)
		randomNode := []*Node{node1, node2, node3}[rand.Intn(3)]
		randomNode.Put(key, value)
	}

	time.Sleep(100 * time.Millisecond)

	// Check all keys are present
	fmt.Println("Checking keys before node leaves:")
	keysFound := 0
	for i := 0; i < 9; i++ {
		key := fmt.Sprintf("key%d", i)
		value, found := node1.Get(key)
		if found {
			fmt.Printf("Key %s found with value %s\n", key, value)
			keysFound++
		} else {
			fmt.Printf("Key %s not found\n", key)
		}
	}

	// One node leaves the network
	network.DisconnectNode(node3.ID)
	time.Sleep(100 * time.Millisecond)

	// Check keys after node leaves
	fmt.Println("\nChecking keys after node leaves:")
	keysFoundAfterLeave := 0
	for i := 0; i < 9; i++ {
		key := fmt.Sprintf("key%d", i)
		value, found := node1.Get(key)
		if found {
			fmt.Printf("Key %s found with value %s\n", key, value)
			keysFoundAfterLeave++
		} else {
			fmt.Printf("Key %s not found\n", key)
		}
	}

	// Calculate percentage of keys lost
	percentLost := float64(keysFound-keysFoundAfterLeave) / float64(keysFound) * 100
	fmt.Printf("\nPercentage of keys lost: %.2f%%\n", percentLost)

	time.Sleep(100 * time.Millisecond)

	jsonInfo, err := network.PrintNetworkInfoJSON()
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println(jsonInfo)
	}

	network.Stop()
}

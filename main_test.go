package main

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestSimpleChordOperations(t *testing.T) {
	network := NewNetwork(true)
	node1 := NewNodeWithRandomID(network)
	node2 := NewNodeWithRandomID(network)
	node3 := NewNodeWithRandomID(network)

	// Join nodes to the network
	node2.Join(node1)
	node3.Join(node1)

	// Allow some time for stabilization
	time.Sleep(1 * time.Second)

	// Test Put and Get operations
	err := node1.Put("key1", "value1")
	if err != nil {
		t.Errorf("Failed to put key1: %v", err)
	}

	value, found := node2.Get("key1")
	if !found || value != "value1" {
		t.Errorf("Failed to get key1 or incorrect value. Found: %v, Value: %s", found, value)
	}

	// Test non-existent key
	_, found = node3.Get("nonexistent")
	if found {
		t.Errorf("Found a non-existent key")
	}

	network.Stop()
}

func TestKeyDistribution(t *testing.T) {
	network := NewNetwork(true)
	initialNode := NewNodeWithRandomID(network)

	numNodes := 10
	nodes := make([]*Node, numNodes)
	nodes[0] = initialNode

	// Slowly add nodes and keys
	for i := 1; i < numNodes; i++ {
		nodes[i] = NewNodeWithRandomID(network)
		err := nodes[i].Join(initialNode)
		if err != nil {
			t.Errorf("Failed to join node %d: %v", i, err)
		}

		// Allow some time for stabilization
		time.Sleep(1 * time.Second)

		// Add 10 keys after each node joins
		for j := 0; j < 10; j++ {
			key := fmt.Sprintf("key%d_%d", i, j)
			value := fmt.Sprintf("value%d_%d", i, j)
			err := initialNode.Put(key, value)
			if err != nil {
				t.Errorf("Failed to put %s: %v", key, err)
			}
		}
	}

	// Allow some time for final stabilization
	time.Sleep(5 * time.Second)

	// Check key distribution
	for i := 1; i < numNodes; i++ {
		for j := 0; j < 10; j++ {
			key := fmt.Sprintf("key%d_%d", i, j)
			value := fmt.Sprintf("value%d_%d", i, j)
			correctNode := nodes[0].FindSuccessor(hash(key))

			// Check if the key is stored on the correct node
			storedValue, found := correctNode.Get(key)
			if !found || storedValue != value {
				t.Errorf("Key %s not stored on the correct node or has incorrect value. Expected node: %v, Found: %v, Value: %s", key, correctNode.ID, found, storedValue)
			}

			// Check if the key is not stored on other nodes
			for _, node := range nodes {
				if node.ID != correctNode.ID {
					_, found := node.Get(key)
					if found {
						t.Errorf("Key %s incorrectly stored on node %v", key, node.ID)
					}
				}
			}
		}
	}

	network.Stop()
}

func TestStressKeys(t *testing.T) {
	network := NewNetwork()
	node := NewNodeWithRandomID(network)
	node2 := NewNodeWithRandomID(network)
	node3 := NewNodeWithRandomID(network)
	node4 := NewNodeWithRandomID(network)

	node2.Join(node)
	node3.Join(node)
	node4.Join(node)

	// Allow some time for stabilization
	time.Sleep(1 * time.Second)

	numKeys := 10000
	var wg sync.WaitGroup

	// Stress test Put operations
	for i := 0; i < numKeys; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("key%d", i)
			value := fmt.Sprintf("value%d", i)
			err := node.Put(key, value)
			if err != nil {
				t.Errorf("Failed to put %s: %v", key, err)
			}
		}(i)
	}
	wg.Wait()

	// Stress test Get operations
	for i := 0; i < numKeys; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("key%d", i)
			expectedValue := fmt.Sprintf("value%d", i)
			value, found := node.Get(key)
			if !found || value != expectedValue {
				t.Errorf("Failed to get %s or incorrect value. Found: %v, Value: %s, Expected: %s", key, found, value, expectedValue)
			}
		}(i)
	}
	wg.Wait()
}

func TestStressNodes(t *testing.T) {
	network := NewNetwork(true)
	initialNode := NewNodeWithRandomID(network)

	numNodes := 25
	nodes := make([]*Node, numNodes)
	nodes[0] = initialNode

	// Stress test node creation and joining
	for i := 1; i < numNodes; i++ {
		nodes[i] = NewNodeWithRandomID(network)
		err := nodes[i].Join(initialNode)
		if err != nil {
			t.Errorf("Failed to join node %d: %v", i, err)
		}
	}

	// Allow some time for stabilization
	time.Sleep(10 * time.Second)

	// Add 10 data to 10 random nodes
	for i := 0; i < 10; i++ {
		randomNode := nodes[rand.Intn(numNodes)]
		key := fmt.Sprintf("testkey%d", i)
		value := fmt.Sprintf("testvalue%d", i)
		err := randomNode.Put(key, value)
		if err != nil {
			t.Errorf("Failed to put %s: %v", key, err)
		}
	}

	// Allow some time for data propagation
	time.Sleep(1 * time.Second)

	// Verify data can be retrieved from any node
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("testkey%d", i)
		expectedValue := fmt.Sprintf("testvalue%d", i)
		randomNode := nodes[rand.Intn(numNodes)]
		value, found := randomNode.Get(key)
		if !found || value != expectedValue {
			t.Errorf("Failed to get %s or incorrect value. Found: %v, Value: %s, Expected: %s", key, found, value, expectedValue)
		}
	}

	// Display the network
	info, err := network.PrintNetworkInfoJSON()
	if err != nil {
		t.Errorf("Failed to print network info: %v", err)
	}
	fmt.Println("Network Info:")
	fmt.Println(info)

	network.Stop()
}

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

// Visualize the entire network and each node's data
func VisualizeNetwork(network *Network) interface{} {
	networkVisualization := struct {
		NumberOfNodes int             `json:"numberOfNodes"`
		Nodes         []VisualizeNode `json:"nodes"`
	}{
		NumberOfNodes: len(network.Nodes),
		Nodes:         make([]VisualizeNode, 0, len(network.Nodes)),
	}

	// Iterate through all nodes in the network
	for _, node := range network.Nodes {
		nodeData := VisualizeNode{
			ID:   node.ID,
			Data: node.Data,
		}

		// Prepare finger table
		fingerTable := make([]int, len(node.FingerTable))
		for i, finger := range node.FingerTable {
			if finger != nil {
				fingerTable[i] = finger.ID
			} else {
				fingerTable[i] = -1 // Use -1 to represent nil
			}
		}
		nodeData.FingerTable = fingerTable

		// Set predecessor and successor
		if node.Predecessor != nil {
			nodeData.Predecessor = node.Predecessor.ID
		} else {
			nodeData.Predecessor = -1
		}

		if node.Successor != nil {
			nodeData.Successor = node.Successor.ID
		} else {
			nodeData.Successor = -1
		}

		networkVisualization.Nodes = append(networkVisualization.Nodes, nodeData)
	}

	return networkVisualization
}

type VisualizeNode struct {
	ID          int            `json:"id"`
	Data        map[int]string `json:"data"`
	FingerTable []int          `json:"fingerTable"`
	Predecessor int            `json:"predecessor"`
	Successor   int            `json:"successor"`
}

// SaveNetworkState periodically saves the network state to a file
func SaveNetworkState(network *Network, interval time.Duration, filename string) {
	var networkStates []interface{}
	ticker := time.NewTicker(interval)
	mutex := &sync.Mutex{}

	go func() {
		for range ticker.C {
			mutex.Lock()
			networkState := VisualizeNetwork(network)
			networkStates = append(networkStates, networkState)
			mutex.Unlock()
		}
	}()

	go func() {
		for {
			time.Sleep(10 * time.Second) // Save to file every 10 seconds
			mutex.Lock()
			jsonData, err := json.MarshalIndent(networkStates, "", "  ")
			if err != nil {
				fmt.Println("Error marshalling network states:", err)
			} else {
				err = os.WriteFile(filename, jsonData, 0644)
				if err != nil {
					fmt.Println("Error writing to file:", err)
				}
			}
			mutex.Unlock()
		}
	}()
}

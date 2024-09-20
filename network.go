package main

import (
	"encoding/json"
	"fmt"
	"sync"
)

// Network struct manages nodes and simulates network operations
type Network struct {
	Nodes map[int]*Node
	mu    sync.Mutex
}

// NewNetwork creates a new Chord network
func NewNetwork() *Network {
	return &Network{Nodes: make(map[int]*Node)}
}

// AddNode adds a node to the network
func (net *Network) AddNode(node *Node) {
	net.mu.Lock()
	defer net.mu.Unlock()
	net.Nodes[node.ID] = node
}

// PrintNetworkInfoJSON prints all the information of the Chord DHT to JSON
func (net *Network) PrintNetworkInfoJSON() (string, error) {
	net.mu.Lock()
	defer net.mu.Unlock()

	type KeyPair struct {
		Key   int    `json:"key"`
		Value string `json:"value"`
	}

	type NodeInfo struct {
		ID          int       `json:"id"`
		Predecessor *int      `json:"predecessor"`
		Successor   int       `json:"successor"`
		FingerTable []int     `json:"fingerTable"`
		Data        []KeyPair `json:"data"`
	}

	type NetworkInfo struct {
		Nodes []NodeInfo `json:"nodes"`
	}

	networkInfo := NetworkInfo{Nodes: make([]NodeInfo, 0, len(net.Nodes))}

	for _, node := range net.Nodes {
		nodeInfo := NodeInfo{
			ID:          node.ID,
			Successor:   node.Successor.ID,
			FingerTable: make([]int, M),
			Data:        make([]KeyPair, 0),
		}

		if node.Predecessor != nil {
			predecessorID := node.Predecessor.ID
			nodeInfo.Predecessor = &predecessorID
		}

		for i, finger := range node.FingerTable {
			if finger != nil {
				nodeInfo.FingerTable[i] = finger.ID
			} else {
				nodeInfo.FingerTable[i] = -1 // Use -1 to represent nil
			}
		}

		for key, value := range node.Data {
			nodeInfo.Data = append(nodeInfo.Data, KeyPair{Key: key, Value: value})
		}

		networkInfo.Nodes = append(networkInfo.Nodes, nodeInfo)
	}

	jsonData, err := json.MarshalIndent(networkInfo, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error marshaling network info to JSON: %v", err)
	}

	return string(jsonData), nil
}

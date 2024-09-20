package main

import "sync"

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

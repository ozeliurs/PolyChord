package main

import (
	"crypto/rand"
	"fmt"
	"hash/fnv"
	"math/big"
	"sync"
	"time"
)

const M = 5 // Defines the size of the identifier space (2^M)

type Node struct {
	ID          int
	Predecessor *Node
	Successor   *Node
	FingerTable []*Node
	Data        map[int]string
	mu          sync.Mutex
	Network     *Network
}

// Create a new node
func NewNode(id int, network *Network) *Node {
	node := &Node{
		ID:          id,
		Predecessor: nil,
		Successor:   nil,
		FingerTable: make([]*Node, M),
		Data:        make(map[int]string),
		Network:     network,
	}
	node.Successor = node // Initially points to itself
	network.AddNode(node)
	return node
}

// NewNodeWithRandomID creates a new node with a random ID
func NewNodeWithRandomID(network *Network) *Node {
	// Generate a random ID within the ID space (e.g., 0 to 2^M - 1)
	maxID := big.NewInt(1 << M) // M is the number of bits, i.e., ID space 0 to 2^M-1
	randomID, err := rand.Int(rand.Reader, maxID)
	if err != nil {
		fmt.Println("Error generating random ID:", err)
		return nil
	}

	return NewNode(int(randomID.Int64()), network)
}

// Hash function to map keys to the identifier space
func hash(s string) int {
	h := fnv.New32a()
	h.Write([]byte(s))
	return int(h.Sum32() % (1 << M))
}

// Find the successor for a given id
func (n *Node) FindSuccessor(id int) *Node {
	if n.Successor == nil || between(id, n.ID, n.Successor.ID) || id == n.Successor.ID {
		return n.Successor
	} else {
		closest := n.ClosestPrecedingNode(id)
		return closest.FindSuccessor(id)
	}
}

// Helper function to determine if 'id' is between two ids in a ring
func between(id, start, end int) bool {
	if start < end {
		return id > start && id < end
	}
	return id > start || id < end
}

// Find the closest preceding node in the finger table
func (n *Node) ClosestPrecedingNode(id int) *Node {
	for i := M - 1; i >= 0; i-- {
		if n.FingerTable[i] != nil && between(n.FingerTable[i].ID, n.ID, id) {
			return n.FingerTable[i]
		}
	}
	return n
}

// Join a new node to the network
func (n *Node) Join(existing *Node) error {
	if existing != nil {
		n.Predecessor = nil
		n.Successor = existing.FindSuccessor(n.ID)
		if n.Successor == nil {
			return nil
		}
	} else {
		n.Predecessor = nil
		n.Successor = n
	}
	return nil
}

// Stabilize the network by fixing successor/predecessor relationships
func (n *Node) Stabilize() {
	n.mu.Lock()
	defer n.mu.Unlock()

	x := n.Successor.Predecessor
	if x != nil && between(x.ID, n.ID, n.Successor.ID) {
		n.Successor = x
	}
	n.Successor.Notify(n)
}

// Notify updates the predecessor
func (n *Node) Notify(p *Node) {
	if n.Predecessor == nil || between(p.ID, n.Predecessor.ID, n.ID) {
		n.Predecessor = p
	}
}

// Fix finger table periodically
func (n *Node) FixFingers() {
	for i := 0; i < M; i++ {
		start := (n.ID + (1 << i)) % (1 << M)
		n.FingerTable[i] = n.FindSuccessor(start)
	}
}

// Put a key-value pair in the DHT
func (n *Node) Put(key, value string) error {
	hashedKey := hash(key)
	successor := n.FindSuccessor(hashedKey)
	if successor == nil {
		return nil
	}
	successor.mu.Lock()
	defer successor.mu.Unlock()
	successor.Data[hashedKey] = value
	return nil
}

// Get a value by key from the DHT
func (n *Node) Get(key string) (string, bool) {
	hashedKey := hash(key)
	successor := n.FindSuccessor(hashedKey)
	if successor == nil {
		return "", false
	}
	successor.mu.Lock()
	defer successor.mu.Unlock()
	value, found := successor.Data[hashedKey]
	return value, found
}

// Periodic stabilization and finger table fixing
func (n *Node) RunStabilizer(interval time.Duration) {
	for {
		n.Stabilize()
		n.FixFingers()
		time.Sleep(interval)
	}
}

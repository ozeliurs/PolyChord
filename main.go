package main

import (
	"fmt"
	"hash/fnv"
	"sync"
)

const M = 5 // Defines the size of the identifier space (2^M)

type Node struct {
	ID          int
	Predecessor *Node
	Successor   *Node
	FingerTable []int
	Data        map[int]string
	mu          sync.Mutex
}

// Hash function to map keys to the identifier space
func hash(s string) int {
	h := fnv.New32a()
	h.Write([]byte(s))
	return int(h.Sum32() % (1 << M))
}

// Create a new node
func NewNode(id int) *Node {
	node := &Node{
		ID:          id,
		Predecessor: nil,
		Successor:   nil,
		FingerTable: make([]int, M),
		Data:        make(map[int]string),
	}
	node.Successor = node // Initially points to itself
	return node
}

// Find the successor for a given id
func (n *Node) FindSuccessor(id int) *Node {
	if n.Successor == nil {
		return n
	}

	if between(id, n.ID, n.Successor.ID) || id == n.Successor.ID {
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
		if between(n.FingerTable[i], n.ID, id) {
			return &Node{ID: n.FingerTable[i]}
		}
	}
	return n
}

// Join a new node to the network
func (n *Node) Join(existing *Node) {
	if existing != nil {
		n.Predecessor = nil
		n.Successor = existing.FindSuccessor(n.ID)
	} else {
		// First node in the network, points to itself
		n.Predecessor = nil
		n.Successor = n
	}
}

// Stabilize the network by fixing successor/predecessor relationships
func (n *Node) Stabilize() {
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
		n.FingerTable[i] = n.FindSuccessor(start).ID
	}
}

// Put a key-value pair in the DHT
func (n *Node) Put(key, value string) {
	hashedKey := hash(key)
	successor := n.FindSuccessor(hashedKey)
	successor.mu.Lock()
	defer successor.mu.Unlock()
	successor.Data[hashedKey] = value
}

// Get a value by key from the DHT
func (n *Node) Get(key string) (string, bool) {
	hashedKey := hash(key)
	successor := n.FindSuccessor(hashedKey)
	successor.mu.Lock()
	defer successor.mu.Unlock()
	value, found := successor.Data[hashedKey]
	return value, found
}

func main() {
	// Create the first node
	node1 := NewNode(1)

	// Create a second node and join the network
	node2 := NewNode(2)
	node2.Join(node1)

	// Stabilize the network
	node1.Stabilize()
	node2.Stabilize()

	// Insert and retrieve data
	node1.Put("foo", "bar")
	value, found := node2.Get("foo")
	if found {
		fmt.Println("Found:", value)
	} else {
		fmt.Println("Key not found")
	}
}

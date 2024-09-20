package main

import (
	"fmt"
	"hash/fnv"
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

// Network struct manages nodes and simulates network operations
type Network struct {
	Nodes map[int]*Node
	mu    sync.Mutex
}

// NewNetwork creates a new Chord network
func NewNetwork() *Network {
	return &Network{Nodes: make(map[int]*Node)}
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

// AddNode adds a node to the network
func (net *Network) AddNode(node *Node) {
	net.mu.Lock()
	defer net.mu.Unlock()
	net.Nodes[node.ID] = node
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
			return fmt.Errorf("could not find successor")
		}
	} else {
		// First node in the network, points to itself
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
		return fmt.Errorf("could not find successor to store key: %s", key)
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

func main() {
	network := NewNetwork()

	// Create the first node
	node1 := NewNode(1, network)
	go node1.RunStabilizer(5 * time.Second)

	// Create additional nodes and join the network
	for i := 2; i <= 5; i++ {
		node := NewNode(i, network)
		err := node.Join(node1)
		if err != nil {
			fmt.Printf("Error joining node %d: %v\n", i, err)
			continue
		}
		go node.RunStabilizer(5 * time.Second)
	}

	// Insert and retrieve data
	node1.Put("foo", "bar")
	time.Sleep(2 * time.Second) // Allow time for stabilization

	value, found := node1.Get("foo")
	if found {
		fmt.Printf("Found value: %s\n", value)
	} else {
		fmt.Println("Key not found")
	}

	// Simulate more keys and values
	node1.Put("baz", "qux")
	node1.Put("apple", "banana")

	// Retrieve the values from the network
	time.Sleep(2 * time.Second) // Allow time for stabilization
	value, found = node1.Get("baz")
	if found {
		fmt.Printf("Found value for 'baz': %s\n", value)
	} else {
		fmt.Println("Key 'baz' not found")
	}

	value, found = node1.Get("apple")
	if found {
		fmt.Printf("Found value for 'apple': %s\n", value)
	} else {
		fmt.Println("Key 'apple' not found")
	}
}

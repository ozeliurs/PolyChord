# Chord DHT Implementation

This project implements a Chord Distributed Hash Table (DHT) in Go. It simulates various scenarios to demonstrate the functionality and robustness of the Chord protocol.

## Features

- Simple Chord network creation and key-value storage
- Stress testing with multiple nodes and a large number of keys
- Network growth simulation
- Node leave scenario and data integrity check

## Prerequisites

- Go 1.15 or higher

## Installation

1. Clone the repository:
   ```
   git clone https://github.com/ozeliurs/PolyChord.git
   ```
2. Navigate to the project directory:
   ```
   cd PolyChord
   ```

## Usage

The program supports four different scenarios. To run a scenario, use the following command:

```
go run ./chord simple|stress_keys|network_growth|node_leave
```

### Scenarios

1. **Simple Scenario** (`simple`):
   Creates a small network with three nodes and demonstrates basic key-value insertion and retrieval.

   ```
   go run ./chord simple
   ```

2. **Stress Keys Scenario** (`stress_keys`):
   Creates a network with 10 nodes and inserts 1000 key-value pairs to test performance under load.

   ```
   go run ./chord stress_keys
   ```

3. **Network Growth Scenario** (`network_growth`):
   Simulates network growth by starting with 3 nodes, inserting some initial keys, and then rapidly adding 25 more nodes. It checks key retrieval before, during, and after network growth.

   ```
   go run ./chord network_growth
   ```

4. **Node Leave Scenario** (`node_leave`):
   Creates a network with three nodes, inserts random key-value pairs, and then simulates one node leaving the network. It checks data integrity before and after the node leaves.

   ```
   go run ./chord node_leave
   ```

## Output

Each scenario will print relevant information to the console, including:
- Network structure
- Key-value operations
- Timing information (for stress testing)
- Data integrity checks

At the end of each scenario, a JSON representation of the network state is printed.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License.

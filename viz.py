import json
import networkx as nx
import matplotlib.pyplot as plt
import math

# Load the JSON data from the file
with open('data.json', 'r') as f:
    data = json.load(f)

# Create a directed graph
G = nx.DiGraph()

# Add nodes and edges based on the successor and predecessor relationships
for node in data['nodes']:
    node_id = node['id']
    successor = node['successor']
    predecessor = node['predecessor']

    G.add_node(node_id)

    # Add edges (successor and predecessor connections)
    if successor is not None:
        G.add_edge(node_id, successor, color='red')
    if predecessor is not None:
        G.add_edge(predecessor, node_id, color='blue')

# Custom circular layout function
def custom_circular_layout(G, radius=1):
    nodes = sorted(G.nodes())
    node_count = len(nodes)
    pos = {}
    for i, node in enumerate(nodes):
        theta = 2 * math.pi * i / node_count
        x = radius * math.cos(theta)
        y = radius * math.sin(theta)
        pos[node] = (x, y)
    return pos

# Set up custom circular layout
pos = custom_circular_layout(G)

# Draw the graph
plt.figure(figsize=(12, 12))

# Draw nodes
nx.draw_networkx_nodes(G, pos, node_size=1000, node_color="lightblue")
nx.draw_networkx_labels(G, pos, labels={node: str(node) for node in G.nodes()},
                        font_size=10, font_weight='bold')

# Draw edges with different colors
edge_colors = [G[u][v]['color'] for u, v in G.edges()]
nx.draw_networkx_edges(G, pos, edge_color=edge_colors, arrowsize=20, arrows=True)

# Display the graph
plt.title("DHT Chord Ring Visualization")
plt.axis('equal')  # Ensure the circle is not distorted
plt.axis('off')  # Turn off the axis
plt.tight_layout()

# Add a legend
red_patch = plt.Line2D([0], [0], color="red", lw=2, label='Successor')
blue_patch = plt.Line2D([0], [0], color="blue", lw=2, label='Predecessor')
plt.legend(handles=[red_patch, blue_patch], loc='upper right')

plt.show()

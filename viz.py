import json
import networkx as nx
import matplotlib.pyplot as plt
import matplotlib.animation as animation
import math
from matplotlib.lines import Line2D

def read_json_lines(file_path):
    with open(file_path, 'r') as file:
        return [json.loads(line) for line in file]

data_frames = read_json_lines('data.json')

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

def create_graph(data):
    G = nx.MultiDiGraph()  # Use MultiDiGraph to allow multiple arrows between nodes
    for node in data['nodes']:
        node_id = node['id']
        successor = node['successor']
        predecessor = node['predecessor']
        G.add_node(node_id)
        if successor is not None:
            G.add_edge(node_id, successor, color='red', connection_type='successor')
        if predecessor is not None:
            G.add_edge(node_id, predecessor, color='blue', connection_type='predecessor')
    return G

fig, ax = plt.subplots(figsize=(12, 12))

def update(frame):
    ax.clear()
    ax.set_axis_off()

    G = create_graph(data_frames[frame])
    pos = custom_circular_layout(G)

    # Draw nodes
    nx.draw_networkx_nodes(G, pos, ax=ax, node_size=1000, node_color="lightblue")
    nx.draw_networkx_labels(G, pos, labels={node: str(node) for node in G.nodes()},
                            font_size=10, font_weight='bold', ax=ax)

    # Draw edges with different arcs for successors and predecessors
    for u, v, key in G.edges(keys=True):  # Iterate over each edge in MultiDiGraph
        edge_color = G[u][v][key]['color']
        connection_type = G[u][v][key]['connection_type']

        # Successor arcs
        if connection_type == 'successor':
            nx.draw_networkx_edges(G, pos, ax=ax, edgelist=[(u, v)], edge_color=edge_color,
                                   connectionstyle='arc3,rad=0.2', arrows=True, arrowsize=20, width=1.5, alpha=0.7)

        # Predecessor arcs
        elif connection_type == 'predecessor':
            nx.draw_networkx_edges(G, pos, ax=ax, edgelist=[(u, v)], edge_color=edge_color,
                                   connectionstyle='arc3,rad=0.3', arrows=True, arrowsize=20, width=1.5, alpha=0.7)

    ax.set_title(f"DHT Chord Ring Visualization - Frame {frame+1}")

    # Add legend
    red_patch = Line2D([0], [0], color="red", lw=2, label='Successor')
    blue_patch = Line2D([0], [0], color="blue", lw=2, label='Predecessor')
    ax.legend(handles=[red_patch, blue_patch], loc='upper right')

    ax.set_xlim(-1.1, 1.1)
    ax.set_ylim(-1.1, 1.1)

    return ax.get_children()

anim = animation.FuncAnimation(fig, update, frames=len(data_frames), interval=20, repeat=False)

plt.tight_layout()
plt.show()

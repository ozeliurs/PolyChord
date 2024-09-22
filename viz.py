import json
import networkx as nx
import matplotlib.pyplot as plt
import matplotlib.animation as animation
import math
from matplotlib.lines import Line2D
import argparse
import random

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

def create_graph(data, show_fingers):
    G = nx.MultiDiGraph()  # Use MultiDiGraph to allow multiple arrows between nodes
    for node in data['nodes']:
        node_id = node['id']
        successor = node['successor']
        predecessor = node['predecessor']
        keys_stored = len(node['data'])
        G.add_node(node_id, keys_stored=keys_stored)
        if successor is not None:
            G.add_edge(node_id, successor, color='red', connection_type='successor')
        if predecessor is not None:
            G.add_edge(node_id, predecessor, color='blue', connection_type='predecessor')
        if show_fingers and 'fingerTable' in node:
            for finger in node['fingerTable']:
                if finger is not None and finger != -1:
                    G.add_edge(node_id, finger, color='green', connection_type='finger')
    return G

fig, ax = plt.subplots(figsize=(12, 12))

parser = argparse.ArgumentParser()
parser.add_argument('--show-fingers', action='store_true', help='Show finger table connections')
args = parser.parse_args()

def update(frame):
    ax.clear()
    ax.set_axis_off()

    G = create_graph(data_frames[frame], args.show_fingers)
    pos = custom_circular_layout(G)

    # Draw nodes
    nx.draw_networkx_nodes(G, pos, ax=ax, node_size=1000, node_color="lightblue")

    nx.draw_networkx_labels(G, pos, labels={node: f"{node}\n({G.nodes[node]['keys_stored']} keys)" for node in G.nodes()},
                            font_size=10, font_weight='bold', ax=ax)

    # Draw edges with different arcs for successors, predecessors, and fingers
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

        # Finger arcs
        elif connection_type == 'finger':
            arc_rad = random.uniform(0.4, 0.5)
            nx.draw_networkx_edges(G, pos, ax=ax, edgelist=[(u, v)], edge_color=edge_color,
                                connectionstyle=f'arc3,rad={arc_rad}', arrows=True, arrowsize=20, width=1, alpha=0.5)

    ax.set_title(f"DHT Chord Ring Visualization - Frame {frame+1}")

    # Add legend
    red_patch = Line2D([0], [0], color="red", lw=2, label='Successor')
    blue_patch = Line2D([0], [0], color="blue", lw=2, label='Predecessor')
    legend_handles = [red_patch, blue_patch]
    if args.show_fingers:
        green_patch = Line2D([0], [0], color="green", lw=2, label='Finger')
        legend_handles.append(green_patch)
    ax.legend(handles=legend_handles, loc='upper right')

    ax.set_xlim(-1.1, 1.1)
    ax.set_ylim(-1.1, 1.1)

    return ax.get_children()

anim = animation.FuncAnimation(fig, update, frames=len(data_frames), interval=20, repeat=False)

anim.save('chord_ring_visualization.gif', writer='pillow', fps=5, progress_callback=(lambda i, n: print(f'Saving frame {i+1} of {n}')))
print("Saving gif, this is so fuuuuuuuuuuuuuucking long... (especially when you've just run my shitty code and a goddamn goroutine is `while True {}` ¯\\_(ツ)_/¯")

plt.tight_layout()
plt.show()

import json
from pathlib import Path

import matplotlib.animation as animation
import matplotlib.pyplot as plt
from matplotlib.colors import to_hex

DATA_FILE = (
    Path(__file__).resolve().parent.resolve().parent
    / "go"
    / "logs"
    / "exp_3_medium_run_1_brute_force.json"
)
with DATA_FILE.open("r", encoding="utf-8") as f:
    data = json.load(f)

points = {int(k): v for k, v in data["points"].items()}
steps = data["steps"]

fig, ax = plt.subplots(figsize=(8, 8))
fig.patch.set_facecolor("#f6f7fb")

xs, ys = zip(*points.values())
client_ids = [i for i in points if i]
max_routes = max(len(step["routes"]) for step in steps)


def build_route_colors(route_count):
    if route_count == 0:
        return []

    cmap = plt.colormaps["turbo"]
    if route_count == 1:
        return [to_hex(cmap(0.2))]

    samples = [0.08 + 0.84 * i / (route_count - 1) for i in range(route_count)]
    return [to_hex(cmap(sample)) for sample in samples]


route_colors = build_route_colors(max_routes)


def draw_points():
    client_xs = [points[i][0] for i in client_ids]
    client_ys = [points[i][1] for i in client_ids]
    depot_x, depot_y = points[0]

    ax.scatter(
        client_xs,
        client_ys,
        c="#1f1f1f",
        s=90,
        edgecolors="white",
        linewidths=1.2,
        zorder=3,
        label="Clients",
    )
    ax.scatter(
        depot_x,
        depot_y,
        c="#e63946",
        s=180,
        marker="*",
        edgecolors="white",
        linewidths=1.2,
        zorder=4,
        label="Depot",
    )

    for i, (x, y) in points.items():
        ax.text(
            x + 0.12,
            y + 0.12,
            str(i),
            fontsize=11,
            weight="bold",
            color="#2b2d42",
            zorder=5,
        )


def draw_routes(routes):
    for idx, route in enumerate(routes):
        color = route_colors[idx % len(route_colors)]

        route_body = [0] + route
        body_xs = [points[node_id][0] for node_id in route_body]
        body_ys = [points[node_id][1] for node_id in route_body]

        ax.plot(
            body_xs,
            body_ys,
            color=color,
            linewidth=3.0,
            alpha=0.9,
            marker="o",
            markersize=5.5,
            label=f"Route {idx + 1}",
            zorder=2,
        )

        if route:
            return_segment = [route[-1], 0]
            return_xs = [points[node_id][0] for node_id in return_segment]
            return_ys = [points[node_id][1] for node_id in return_segment]
            ax.plot(
                return_xs,
                return_ys,
                color=color,
                linewidth=3.0,
                alpha=0.9,
                linestyle="--",
                label="_nolegend_",
                zorder=2,
            )


def update(frame):
    ax.clear()
    ax.set_facecolor("#ffffff")

    step = steps[frame]

    draw_points()
    draw_routes(step["routes"])

    ax.set_title(
        f"Step {step['step_id']} | Cost: {step['cost']:.2f}",
        fontsize=15,
        weight="bold",
        pad=14,
    )

    ax.text(
        0.02,
        0.98,
        f"Routes: {len(step['routes'])}  |  Clients: {len(client_ids)}",
        transform=ax.transAxes,
        va="top",
        ha="left",
        fontsize=10,
        color="#495057",
        bbox={
            "boxstyle": "round,pad=0.35",
            "facecolor": "#f8f9fa",
            "edgecolor": "#dee2e6",
            "alpha": 0.95,
        },
    )

    ax.legend(loc="upper right", frameon=True, framealpha=0.96)
    ax.grid(True, linestyle="--", linewidth=0.8, alpha=0.25)
    ax.set_aspect("equal", adjustable="box")
    ax.margins(0.15)

    ax.set_xlim(min(xs) - 1, max(xs) + 1)
    ax.set_ylim(min(ys) - 1, max(ys) + 1)
    return ax.lines


ani = animation.FuncAnimation(
    fig, update, frames=len(steps), interval=1500, repeat=True
)

plt.show()

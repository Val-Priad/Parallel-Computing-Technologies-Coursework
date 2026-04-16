import json
import math
from pathlib import Path

import matplotlib.pyplot as plt

LOGS_DIR = (
    Path(__file__).resolve().parent.parent / "go" / "results" / "brute_vs_aco"
)
SOLUTIONS_DIR = Path(__file__).resolve().parent / "results" / "routes"


def build_colors(n):
    cmap = plt.colormaps["turbo"]
    return [cmap(i / max(1, n - 1)) for i in range(n)]


def plot_points(ax, points):
    xs, ys = zip(*points.values())

    ax.scatter(*points[0], marker="*", s=150, label="Depot")

    client_points = [p for i, p in points.items() if i != 0]
    cx, cy = zip(*client_points)
    ax.scatter(cx, cy, s=70, label="Clients")

    for i, (x, y) in points.items():
        ax.annotate(str(i), (x, y), textcoords="offset points", xytext=(4, 4))


def plot_routes(ax, routes, points, colors):
    for i, route in enumerate(routes):
        if not route:
            continue

        path = [0] + route + [0]
        xs = [points[n][0] for n in path]
        ys = [points[n][1] for n in path]

        ax.plot(
            xs,
            ys,
            marker="o",
            linewidth=2.5,
            color=colors[i],
            label=f"Route {i + 1}",
        )


def plot_step(ax, step, points, colors):
    plot_points(ax, points)
    plot_routes(ax, step["routes"], points, colors)

    ax.set_title(f"Step {step['step_id'] + 1} | Cost: {step['cost']:.2f}")
    ax.set_aspect("equal")
    ax.grid(True)

    ax.legend(fontsize=8)


def visualize(data, name):
    points = {int(k): v for k, v in data["points"].items()}
    steps = data["steps"]

    if not steps:
        return

    n = len(steps)
    cols = math.ceil(math.sqrt(n))
    rows = math.ceil(n / cols)

    colors = build_colors(max(len(s["routes"]) for s in steps))

    fig, axes = plt.subplots(
        rows, cols, figsize=(cols * 5, rows * 5), layout="constrained"
    )

    axes = (
        axes.flatten()  # type: ignore
        if isinstance(axes, (list, tuple)) or hasattr(axes, "flat")
        else [axes]
    )

    for ax, step in zip(axes, steps):
        plot_step(ax, step, points, colors)

    for ax in axes[len(steps) :]:
        ax.set_visible(False)

    fig.suptitle(f"VRP Solution\n{name}", fontsize=16)

    output_path = SOLUTIONS_DIR / f"{Path(name).stem}.png"
    fig.savefig(output_path, dpi=200)
    plt.close(fig)


def main():
    files = sorted(LOGS_DIR.glob("*.json"))
    SOLUTIONS_DIR.mkdir(parents=True, exist_ok=True)

    if not files:
        raise FileNotFoundError("No logs found")

    for f in files:
        with open(f, encoding="utf-8") as fp:
            data = json.load(fp)

        visualize(data, f.name)


if __name__ == "__main__":
    main()

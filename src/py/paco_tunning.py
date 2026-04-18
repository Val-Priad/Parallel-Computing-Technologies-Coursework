import csv
from pathlib import Path

import matplotlib.pyplot as plt

CSV_PATH = (
    Path(__file__).resolve().parent.parent
    / "go"
    / "results"
    / "paco_process_tuning.csv"
)

OUTPUT_DIR = Path(__file__).resolve().parent / "results" / "paco_tunning"


def read_results():
    with open(CSV_PATH, newline="", encoding="utf-8") as fp:
        return list(csv.DictReader(fp))


def prepare_data(rows):
    data = [
        (
            int(r["workers"]),
            float(r["time_ms"]),
            float(r["cost"]),
        )
        for r in rows
    ]

    if not data:
        raise ValueError("No PACO tuning data")

    data.sort()
    workers, time_ms, cost = zip(*data)
    return workers, time_ms, cost


def plot_cost(workers, cost):
    plt.figure(figsize=(8, 5))
    plt.plot(workers, cost, marker="o")

    for x, y in zip(workers, cost):
        plt.annotate(
            f"{y:.1f}",
            (x, y),
            textcoords="offset points",
            xytext=(-10, 5),
        )

    plt.xlabel("Workers")
    plt.ylabel("Cost")
    plt.title("PACO scaling: workers vs cost")
    plt.grid()

    plt.savefig(OUTPUT_DIR / "cost_vs_workers.png", dpi=200)
    plt.close()


def plot_speedup(workers, time_ms):
    t1 = time_ms[0]
    speedup = [t1 / t for t in time_ms]

    plt.figure(figsize=(8, 5))
    plt.plot(workers, speedup, marker="o")

    for x, y in zip(workers, speedup):
        plt.annotate(
            f"{y:.2f}x",
            (x, y),
            textcoords="offset points",
            xytext=(-10, 5),
        )

    plt.xlabel("Workers")
    plt.ylabel("Speedup (T1 / Tp)")
    plt.title("PACO speedup")
    plt.grid()

    plt.savefig(OUTPUT_DIR / "speedup.png", dpi=200)
    plt.close()


def main():
    OUTPUT_DIR.mkdir(parents=True, exist_ok=True)

    rows = read_results()
    workers, time_ms, cost = prepare_data(rows)

    plot_cost(workers, cost)
    plot_speedup(workers, time_ms)


if __name__ == "__main__":
    main()

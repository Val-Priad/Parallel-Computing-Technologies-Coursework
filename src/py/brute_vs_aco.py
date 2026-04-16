import csv
from pathlib import Path

import matplotlib.pyplot as plt

CSV_PATH = (
    Path(__file__).resolve().parent.parent
    / "go"
    / "results"
    / "brute_force_vs_aco.csv"
)

OUTPUT_DIR = Path(__file__).resolve().parent / "results" / "brute_vs_aco"


def read_results():
    with open(CSV_PATH, newline="", encoding="utf-8") as fp:
        return list(csv.DictReader(fp))


def plot_algorithm(rows, algorithm, filename):
    data = [
        (int(r["customers"]), float(r["time_ms"]))
        for r in rows
        if r["algorithm"] == algorithm
    ]

    if not data:
        raise ValueError(f"No data for {algorithm}")

    data.sort()
    customers, times = zip(*data)

    plt.figure(figsize=(8, 5))
    plt.plot(customers, times, marker="o")

    for x, y in zip(customers, times):
        plt.annotate(
            f"{y:.3f}",
            (x, y),
            textcoords="offset points",
            xytext=(-15, 4),
        )

    plt.title(f"{algorithm} time growth")
    plt.xlabel("Customers")
    plt.ylabel("Time (ms)")
    plt.grid()

    plt.savefig(OUTPUT_DIR / filename, dpi=200)
    plt.close()


def main():
    OUTPUT_DIR.mkdir(parents=True, exist_ok=True)

    rows = read_results()

    plot_algorithm(rows, "aco", "aco.png")
    plot_algorithm(rows, "brute_force", "brute_force.png")


if __name__ == "__main__":
    main()

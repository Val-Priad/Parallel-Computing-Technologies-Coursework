import csv
from pathlib import Path

import matplotlib.pyplot as plt

CSV_PATH = (
    Path(__file__).resolve().parent.parent / "go" / "results" / "aco.csv"
)

OUTPUT_DIR = Path(__file__).resolve().parent / "results" / "aco"


def read_results():
    with open(CSV_PATH, newline="", encoding="utf-8") as fp:
        return list(csv.DictReader(fp))


def plot_aco(rows):
    data = [
        (
            int(r["customers"]),
            float(r["time_ms"]),
            int(r["vehicles"]),
        )
        for r in rows
    ]

    if not data:
        raise ValueError("No ACO data")

    data.sort()

    customers, times, vehicles = zip(*data)

    plt.figure(figsize=(8, 5))
    plt.plot(customers, times, marker="o")

    for x, y, v in zip(customers, times, vehicles):
        plt.annotate(
            f"{y:.3f}\n(v={v})",
            (x, y),
            textcoords="offset points",
            xytext=(-15, 4),
        )

    plt.title("ACO execution time growth")
    plt.xlabel("Customers")
    plt.ylabel("Time (ms)")
    plt.grid()

    plt.savefig(OUTPUT_DIR / "aco_growth.png", dpi=200)
    plt.close()


def main():
    OUTPUT_DIR.mkdir(parents=True, exist_ok=True)

    rows = read_results()
    plot_aco(rows)


if __name__ == "__main__":
    main()

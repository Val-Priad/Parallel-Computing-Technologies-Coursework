import csv
from pathlib import Path

import matplotlib.pyplot as plt

CSV_PATH = (
    Path(__file__).resolve().parent.parent
    / "go"
    / "results"
    / "aco_vs_paco.csv"
)

OUTPUT_DIR = Path(__file__).resolve().parent / "solutions" / "aco_vs_paco"


def read_results():
    with open(CSV_PATH, newline="", encoding="utf-8") as fp:
        return list(csv.DictReader(fp))


def prepare_pairs(rows):
    grouped = {}

    for r in rows:
        key = r["experiment"]
        grouped.setdefault(key, {})[r["algorithm"]] = r

    pairs = []
    for run, algos in grouped.items():
        if "aco" in algos and "paco" in algos:
            aco = algos["aco"]
            paco = algos["paco"]

            pairs.append(
                (
                    int(aco["customers"]),
                    float(aco["time_ms"]),
                    float(paco["time_ms"]),
                    float(aco["cost"]),
                    float(paco["cost"]),
                )
            )

    return sorted(pairs)


def plot_speedup(pairs):
    customers = []
    speedups = []

    for c, t_aco, t_paco, _, _ in pairs:
        customers.append(c)
        speedups.append(t_aco / t_paco)

    plt.figure(figsize=(8, 5))
    plt.plot(customers, speedups, marker="o")

    for x, y in zip(customers, speedups):
        plt.annotate(
            f"{y:.2f}x",
            (x, y),
            textcoords="offset points",
            xytext=(-10, 5),
        )

    plt.title("PACO speedup over ACO")
    plt.xlabel("Customers")
    plt.ylabel("Speedup (ACO / PACO)")
    plt.grid()

    plt.savefig(OUTPUT_DIR / "speedup.png", dpi=200)
    plt.close()


def plot_quality(pairs):
    customers = []
    diffs = []

    for c, _, _, cost_aco, cost_paco in pairs:
        customers.append(c)
        diffs.append(cost_paco - cost_aco)

    plt.figure(figsize=(8, 5))
    plt.plot(customers, diffs, marker="o")

    for x, y in zip(customers, diffs):
        plt.annotate(
            f"{y:.2f}",
            (x, y),
            textcoords="offset points",
            xytext=(-10, 5),
        )

    plt.axhline(0, linestyle="--")  # линия равенства

    plt.title("Solution quality difference (PACO - ACO)")
    plt.xlabel("Customers")
    plt.ylabel("Cost difference")
    plt.grid()

    plt.savefig(OUTPUT_DIR / "quality.png", dpi=200)
    plt.close()


def main():
    OUTPUT_DIR.mkdir(parents=True, exist_ok=True)

    rows = read_results()
    pairs = prepare_pairs(rows)

    plot_speedup(pairs)
    plot_quality(pairs)


if __name__ == "__main__":
    main()

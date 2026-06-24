# Parallel ACO for Vehicle Routing Problem

## ⚙️ Technologies

[![Technologies](https://skillicons.dev/icons?i=go,python&perline=5)](https://skillicons.dev)

- Go
- Python
- Matplotlib
- CSV
- JSON

## 🧠 What I Learned

### 🐜 Ant Colony Optimization

- Implemented the Ant Colony Optimization algorithm for solving the Vehicle Routing Problem.
- Improved understanding of pheromone-based search, heuristic distance influence, evaporation, and elite solution reinforcement.
- Practiced working with randomized algorithms using configurable seeds for reproducible experiments.

### ⚡ Parallel Algorithm Design

- Implemented a parallel version of ACO using multiple workers.
- Practiced splitting computation between goroutines.
- Improved understanding of synchronization, shared best solution exchange, and performance scaling.

### 📊 Experimental Analysis

- Compared brute force search, sequential ACO, and parallel ACO.
- Measured execution time and solution cost across different generated VRP instances.
- Generated CSV reports for further analysis and visualization.

### 📈 Data Visualization

- Created Python scripts for plotting algorithm performance.
- Visualized execution time growth, speedup, solution quality difference, and route construction.
- Practiced separating algorithm implementation from result analysis.

## ✨ Features

- Generates random Vehicle Routing Problem instances.
- Supports different vehicle capacity modes.
- Solves small VRP instances using brute force search.
- Solves larger VRP instances using Ant Colony Optimization.
- Provides a parallel ACO implementation with configurable worker count.
- Compares ACO and PACO by execution time and solution quality.
- Exports experiment results to CSV files.
- Saves route construction logs in JSON format.
- Builds plots from experiment results using Python.

## 🚀 How to Run

### 1. Clone the repository

```bash
git clone <repository-url>
cd parallel-aco
```

### 2. Run Go tests

```bash
cd src/go
go test ./...
```

This command checks that the implemented solvers return valid VRP solutions and compares small instances against the brute force baseline.

### 3. Run brute force vs ACO comparison

```bash
go run ./cmd/brute_vs_aco
```

This command runs experiments on small VRP instances.

It compares:

- brute force search
- sequential ACO

It produces:

```text
src/go/results/brute_force_vs_aco.csv
src/go/results/brute_vs_aco/*.json
```

The CSV file contains time and cost metrics.
The JSON files contain route improvement steps that can later be visualized.

### 4. Run ACO experiment on larger instances

```bash
go run ./cmd/aco_experiment
```

This command runs sequential ACO on larger generated VRP instances.

It produces:

```text
src/go/results/aco.csv
```

This file is used to analyze how ACO execution time grows when the number of customers increases.

### 5. Run ACO vs PACO comparison

```bash
go run ./cmd/aco_vs_paco
```

This command compares:

- sequential ACO
- parallel ACO

It runs multiple trials for each generated instance and saves average results.

It produces:

```text
src/go/results/aco_vs_paco.csv
```

This file is used to compare execution time, speedup, and solution quality.

### 6. Run ACO parameter tuning

```bash
go run ./cmd/tune_aco
```

This command evaluates several ACO configurations and prints the best configuration according to average solution cost and execution time.

It shows:

- tested configurations
- average cost
- average execution time
- selected best configuration

### 7. Run PACO worker tuning

```bash
go run ./cmd/tune_paco
```

This command tests PACO with different numbers of workers.

It produces:

```text
src/go/results/paco_process_tuning.csv
```

This file is used to analyze how the number of workers affects speed and solution cost.

## 📊 How to Generate Visualizations

Before running Python scripts, install the required dependency:

```bash
pip install matplotlib
```

Run the scripts from the repository root or directly from the `src/py` directory.

### ACO time growth

```bash
python src/py/aco.py
```

Input:

```text
src/go/results/aco.csv
```

Output:

```text
src/py/results/aco/aco_growth.png
```

This plot shows how sequential ACO execution time changes as the number of customers increases.

### Brute force vs ACO plots

```bash
python src/py/brute_vs_aco.py
```

Input:

```text
src/go/results/brute_force_vs_aco.csv
```

Output:

```text
src/py/results/brute_vs_aco/aco.png
src/py/results/brute_vs_aco/brute_force.png
```

These plots show the execution time growth of brute force and ACO separately.

### ACO vs PACO plots

```bash
python src/py/aco_vs_paco.py
```

Input:

```text
src/go/results/aco_vs_paco.csv
```

Output:

```text
src/py/results/aco_vs_paco/time_comparison.png
src/py/results/aco_vs_paco/speedup.png
src/py/results/aco_vs_paco/quality.png
```

These plots show:

- execution time comparison
- PACO speedup over ACO
- solution cost difference between PACO and ACO

### PACO worker scaling plots

```bash
python src/py/paco_tunning.py
```

Input:

```text
src/go/results/paco_process_tuning.csv
```

Output:

```text
src/py/results/paco_tunning/cost_vs_workers.png
src/py/results/paco_tunning/speedup.png
```

These plots show how the number of workers affects PACO speed and solution quality.

### Route visualization

```bash
python src/py/routes.py
```

Input:

```text
src/go/results/brute_vs_aco/*.json
```

Output:

```text
src/py/results/routes/*.png
```

This script visualizes generated VRP instances and route improvement steps saved during the brute force vs ACO experiment.

## 🧭 Recommended Execution Order

For a complete experiment workflow, run the project in this order:

```bash
cd src/go
go test ./...
go run ./cmd/brute_vs_aco
go run ./cmd/aco_experiment
go run ./cmd/aco_vs_paco
go run ./cmd/tune_aco
go run ./cmd/tune_paco
cd ../..
python src/py/brute_vs_aco.py
python src/py/aco.py
python src/py/aco_vs_paco.py
python src/py/paco_tunning.py
python src/py/routes.py
```

### What each step shows

1. `go test ./...`
   Verifies that the solvers produce valid routes and respect VRP constraints.

2. `go run ./cmd/brute_vs_aco`
   Shows how ACO compares with the exact brute force method on small instances.

3. `go run ./cmd/aco_experiment`
   Shows how sequential ACO performs on larger VRP instances.

4. `go run ./cmd/aco_vs_paco`
   Shows whether the parallel version is faster than sequential ACO and how solution quality changes.

5. `go run ./cmd/tune_aco`
   Helps choose a better ACO parameter configuration.

6. `go run ./cmd/tune_paco`
   Shows how the number of parallel workers affects PACO performance.

7. Python visualization scripts
   Convert CSV and JSON experiment outputs into plots for performance analysis and route visualization.

## 📁 Project Structure

```text
.
├── go
│   ├── cmd
│   │   ├── aco_experiment
│   │   ├── aco_vs_paco
│   │   ├── brute_vs_aco
│   │   ├── tune_aco
│   │   └── tune_paco
│   └── internal
│       ├── experiment
│       ├── logging
│       ├── solver
│       └── vrp
└── py
    ├── aco.py
    ├── aco_vs_paco.py
    ├── brute_vs_aco.py
    ├── paco_tunning.py
    └── routes.py
```

## 📌 Results

The project generates several types of results:

- CSV files with execution time and solution cost.
- JSON logs with route improvement steps.
- PNG charts for algorithm comparison.
- Route visualizations for selected VRP instances.

These outputs provide a clear overview of algorithm performance, solution quality, speedup, and route construction behavior.

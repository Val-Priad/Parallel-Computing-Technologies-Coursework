# Parallel ACO for Vehicle Routing Problem

This project implements sequential and parallel Ant Colony Optimization algorithms for solving the Vehicle Routing Problem.

It compares brute force search, sequential ACO, and parallel ACO using generated VRP instances, CSV reports, JSON route logs, and Python visualizations.

## ⚙️ Technologies

[![Technologies](https://skillicons.dev/icons?i=go,python&perline=5)](https://skillicons.dev)

- Go
- Python
- Matplotlib
- CSV
- JSON

## 📝 Project Background

I developed this project as my coursework for the subject **Parallel Computing Technologies**.

I first encountered the Ant Colony Optimization algorithm during my second year of university. At that time, it seemed extremely complex to me, both conceptually and mathematically. In my third year, I had the opportunity to return to this topic, study it more deeply, and finally understand how the algorithm works in practice.

### Why Ant Colony Optimization

The main idea of the project was to implement and analyze the Ant Colony Optimization algorithm for solving the Vehicle Routing Problem. I wanted not only to reproduce the algorithm, but also to understand its internal logic: pheromone updates, heuristic influence, evaporation, route construction, and solution improvement over multiple iterations.

To be honest, the mathematical part was the most challenging aspect for me. I understood the general principles of the algorithm, but some of the formulas were difficult to keep fully in my head. Because of that, I worked with scientific articles and other academic sources while implementing the project. In the written coursework report, I explained the theoretical background, the formulas, and the reasoning behind the algorithm in more detail.

### Why I Chose Go

For the implementation, I chose **Go** deliberately. There were many possible options, such as Python, CUDA, Java, C#, C++, and others. However, Go seemed like the best choice for this project because of its simplicity, strong performance, and built-in support for concurrency through goroutines.

This project also helped me understand why Go is often described as a simpler and more lightweight alternative to C++. I really enjoyed how direct and practical the language feels. Its concurrency model made it much easier to experiment with parallel execution and worker-based computation.

### What I Learned

During this project, I gained a much deeper understanding of Ant Colony Optimization and parallel algorithm design. I learned how to split computation between multiple workers, synchronize shared data, exchange the best solution between goroutines, and measure the performance difference between sequential and parallel implementations.

The main goal was to speed up the sequential version of the algorithm by implementing a parallel version. As a result, I was able to achieve a noticeable performance improvement while preserving the quality of the solution.

### Result

Overall, I am satisfied with this project. It was difficult, but I put a lot of effort into understanding the algorithm, implementing it correctly, and analyzing the results. This coursework helped me improve both my knowledge of parallel computing and my practical experience with Go.

Most importantly, this project made me realize that I really enjoy working with Go. I like its simplicity, its performance, and the way it handles concurrency. Because of this project, I became much more confident in using Go for algorithmic and performance-oriented tasks.

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
python src/py/paco_tuning.py
```

Input:

```text
src/go/results/paco_process_tuning.csv
```

Output:

```text
src/py/results/paco_tuning/cost_vs_workers.png
src/py/results/paco_tuning/speedup.png
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
python src/py/paco_tuning.py
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
└── src
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
        ├── paco_tuning.py
        └── routes.py
```

## 📌 Results

The project generates several types of results:

- CSV files with execution time and solution cost.
- JSON logs with route improvement steps.
- PNG charts for algorithm comparison.
- Route visualizations for selected VRP instances.

These outputs provide a clear overview of algorithm performance, solution quality, speedup, and route construction behavior.

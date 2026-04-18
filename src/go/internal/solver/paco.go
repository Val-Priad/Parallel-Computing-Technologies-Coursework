package solver

import (
	"math"
	"math/rand"
	"parallel-aco/internal/vrp"
	"runtime"
	"sync"
	"time"
)

type PACOConfig struct {
	BaseConfig ACOConfig
	NumWorkers int
}

func SolvePACO(instance vrp.VRPInstance, cfg PACOConfig) vrp.Solution {
	startTime := time.Now()

	if len(instance.Customers) == 0 || instance.Vehicles == 0 {
		return vrp.Solution{
			Routes:     []vrp.Route{},
			Cost:       0,
			DurationMS: float64(time.Since(startTime).Milliseconds()),
		}
	}

	applyConfigDefaults(&cfg.BaseConfig)
	if !validateInstance(instance) {
		return solutionWithMetrics([]vrp.Route{}, math.Inf(1), startTime)
	}

	cfg.NumWorkers = resolvePACOWorkers(cfg.NumWorkers, cfg.BaseConfig.NumAnts)

	type result struct {
		best vrp.Solution
	}

	results := make([]result, cfg.NumWorkers)
	antsPerWorker := splitAnts(cfg.BaseConfig.NumAnts, cfg.NumWorkers)

	var wg sync.WaitGroup

	for w := 0; w < cfg.NumWorkers; w++ {
		wg.Add(1)

		go func(workerID int) {
			defer wg.Done()

			localCfg := cfg.BaseConfig
			localCfg.Seed = cfg.BaseConfig.Seed + int64(workerID)*1000
			localCfg.NumAnts = antsPerWorker[workerID]

			rng := rand.New(rand.NewSource(localCfg.Seed))
			n := len(instance.Dist)
			pheromone := makeMatrix(n, n, localCfg.InitialPheromone)
			mem := NewWorkerMemory(n, len(instance.Customers))

			best := vrp.Solution{
				Routes: []vrp.Route{},
				Cost:   math.Inf(1),
			}

			for iter := 0; iter < localCfg.Iterations; iter++ {
				ants := make([]antSolution, 0, localCfg.NumAnts)

				for ant := 0; ant < localCfg.NumAnts; ant++ {
					sol, feasible := buildSolution(instance, pheromone, localCfg, rng, mem)

					ants = append(ants, antSolution{
						Solution: sol,
						Feasible: feasible,
					})

					if feasible && sol.Cost < best.Cost {
						best = cloneSolution(sol)
					}
				}

				evaporate(pheromone, localCfg.Evaporation)

				for _, ant := range ants {
					if !ant.Feasible || ant.Solution.Cost <= 0 {
						continue
					}
					depositSolution(pheromone, ant.Solution, localCfg.Q/ant.Solution.Cost)
				}

				if !math.IsInf(best.Cost, 1) && best.Cost > 0 {
					depositSolution(pheromone, best, localCfg.EliteWeight*localCfg.Q/best.Cost)
				}
			}

			results[workerID] = result{best: best}
		}(w)
	}

	wg.Wait()

	globalBest := vrp.Solution{
		Routes: []vrp.Route{},
		Cost:   math.Inf(1),
	}

	for _, r := range results {
		if r.best.Cost < globalBest.Cost {
			globalBest = cloneSolution(r.best)
		}
	}

	if math.IsInf(globalBest.Cost, 1) {
		return solutionWithMetrics([]vrp.Route{}, math.Inf(1), startTime)
	}

	return solutionWithMetrics(globalBest.Routes, globalBest.Cost, startTime)
}

func resolvePACOWorkers(numWorkers, numAnts int) int {
	if numWorkers <= 0 {
		numWorkers = runtime.NumCPU()
	}
	if numWorkers > numAnts {
		numWorkers = numAnts
	}
	if numWorkers <= 0 {
		numWorkers = 1
	}

	return numWorkers
}

func splitAnts(total, workers int) []int {
	result := make([]int, workers)

	base := total / workers
	rem := total % workers

	for i := 0; i < workers; i++ {
		result[i] = base
		if i < rem {
			result[i]++
		}
	}

	return result
}

package main

import (
	"math"
	"math/rand"
	"time"
)

type ACOConfig struct {
	NumAnts          int
	Iterations       int
	Alpha            float64
	Beta             float64
	Evaporation      float64
	Q                float64
	InitialPheromone float64
	EliteWeight      float64
	Seed             int64
}

type antSolution struct {
	Solution Solution
	Feasible bool
}

func DefaultACOConfig() ACOConfig {
	return ACOConfig{
		NumAnts:          20,
		Iterations:       100,
		Alpha:            1.0,
		Beta:             3.0,
		Evaporation:      0.5,
		Q:                100.0,
		InitialPheromone: 1.0,
		EliteWeight:      2.0,
		Seed:             time.Now().UnixNano(),
	}
}

func SolveACO(instance VRPInstance, logger *Logger, cfg ACOConfig) Solution {
	startTime := time.Now()

	if len(instance.Customers) == 0 || instance.Vehicles == 0 {
		return Solution{
			Routes: []Route{},
			Cost:   0,
			Metrics: SearchMetrics{
				DurationMS: float64(time.Since(startTime).Nanoseconds()) / 1e6,
			},
		}
	}

	applyACOConfigDefaults(&cfg)

	n := len(instance.Dist)
	pheromone := makeMatrix(n, n, cfg.InitialPheromone)

	rng := rand.New(rand.NewSource(cfg.Seed))

	best := Solution{
		Routes: []Route{},
		Cost:   math.Inf(1),
	}

	stepID := 0

	for iter := 0; iter < cfg.Iterations; iter++ {
		ants := make([]antSolution, 0, cfg.NumAnts)
		iterBest := Solution{Cost: math.Inf(1)}

		for ant := 0; ant < cfg.NumAnts; ant++ {
			sol, feasible := buildAntSolution(instance, pheromone, cfg, rng)

			ants = append(ants, antSolution{
				Solution: sol,
				Feasible: feasible,
			})

			if feasible && sol.Cost < iterBest.Cost {
				iterBest = sol
			}
			if feasible && sol.Cost < best.Cost {
				best = cloneSolution(sol)

				if logger != nil {
					loggedRoutes := make([][]int, len(best.Routes))
					for i, route := range best.Routes {
						loggedRoutes[i] = append([]int{}, route.Nodes...)
					}

					logger.Log(Step{
						StepID: stepID,
						Routes: loggedRoutes,
						Cost:   best.Cost,
					})
					stepID++
				}
			}
		}

		evaporate(pheromone, cfg.Evaporation)

		for _, ant := range ants {
			if !ant.Feasible || math.IsInf(ant.Solution.Cost, 1) || ant.Solution.Cost <= 0 {
				continue
			}
			depositSolution(pheromone, ant.Solution, cfg.Q/ant.Solution.Cost)
		}

		if !math.IsInf(iterBest.Cost, 1) && iterBest.Cost > 0 {
			depositSolution(pheromone, iterBest, cfg.EliteWeight*cfg.Q/iterBest.Cost)
		}

		if !math.IsInf(best.Cost, 1) && best.Cost > 0 {
			depositSolution(pheromone, best, cfg.EliteWeight*cfg.Q/best.Cost)
		}
	}

	if math.IsInf(best.Cost, 1) {
		return Solution{
			Routes: []Route{},
			Cost:   math.Inf(1),
			Metrics: SearchMetrics{
				DurationMS: float64(time.Since(startTime).Nanoseconds()) / 1e6,
			},
		}
	}

	best.Metrics = SearchMetrics{
		DurationMS: float64(time.Since(startTime).Nanoseconds()) / 1e6,
	}

	return best
}

func applyACOConfigDefaults(cfg *ACOConfig) {
	if cfg.NumAnts <= 0 {
		cfg.NumAnts = 20
	}
	if cfg.Iterations <= 0 {
		cfg.Iterations = 100
	}
	if cfg.Alpha <= 0 {
		cfg.Alpha = 1.0
	}
	if cfg.Beta <= 0 {
		cfg.Beta = 3.0
	}
	if cfg.Evaporation <= 0 || cfg.Evaporation >= 1 {
		cfg.Evaporation = 0.5
	}
	if cfg.Q <= 0 {
		cfg.Q = 100.0
	}
	if cfg.InitialPheromone <= 0 {
		cfg.InitialPheromone = 1.0
	}
	if cfg.EliteWeight < 0 {
		cfg.EliteWeight = 0
	}
	if cfg.Seed == 0 {
		cfg.Seed = time.Now().UnixNano()
	}
}

func buildAntSolution(
	instance VRPInstance,
	pheromone [][]float64,
	cfg ACOConfig,
	rng *rand.Rand,
) (Solution, bool) {
	n := len(instance.Dist)
	capacity := instance.VehicleCapacity
	vehicles := instance.Vehicles

	if n == 0 || vehicles == 0 {
		return Solution{Routes: []Route{}, Cost: 0}, true
	}

	demandByID := make([]int, n)
	customerSet := make([]bool, n)
	for _, c := range instance.Customers {
		if c.ID >= 0 && c.ID < n {
			demandByID[c.ID] = c.Demand
			customerSet[c.ID] = true
		}
	}

	visited := make([]bool, n)
	remaining := len(instance.Customers)

	routes := make([]Route, 0, vehicles)

	for v := 0; v < vehicles; v++ {
		current := 0
		currentLoad := 0
		nodes := make([]int, 0)

		for {
			candidates := make([]int, 0)
			for _, c := range instance.Customers {
				if c.ID < 0 || c.ID >= n {
					continue
				}
				if visited[c.ID] {
					continue
				}
				if currentLoad+demandByID[c.ID] > capacity {
					continue
				}
				candidates = append(candidates, c.ID)
			}

			if len(candidates) == 0 {
				break
			}

			next := selectNextCustomer(current, candidates, instance.Dist, pheromone, cfg, rng)
			if next == -1 {
				break
			}

			nodes = append(nodes, next)
			visited[next] = true
			currentLoad += demandByID[next]
			current = next
			remaining--
		}

		routes = append(routes, Route{
			VehicleID: v,
			Nodes:     nodes,
		})

		if remaining == 0 {
			for vv := v + 1; vv < vehicles; vv++ {
				routes = append(routes, Route{
					VehicleID: vv,
					Nodes:     []int{},
				})
			}
			break
		}
	}

	if remaining > 0 {
		return Solution{
			Routes: routes,
			Cost:   math.Inf(1),
		}, false
	}

	for _, c := range instance.Customers {
		if c.ID < 0 || c.ID >= n || !customerSet[c.ID] || !visited[c.ID] {
			return Solution{
				Routes: routes,
				Cost:   math.Inf(1),
			}, false
		}
	}

	totalCost := 0.0
	for _, route := range routes {
		totalCost += computeRouteCost(route.Nodes, instance.Dist)
	}

	return Solution{
		Routes: routes,
		Cost:   totalCost,
	}, true
}

func selectNextCustomer(
	current int,
	candidates []int,
	dist [][]float64,
	pheromone [][]float64,
	cfg ACOConfig,
	rng *rand.Rand,
) int {
	if len(candidates) == 0 {
		return -1
	}
	if len(candidates) == 1 {
		return candidates[0]
	}

	type weightedCandidate struct {
		Node   int
		Weight float64
	}

	weighted := make([]weightedCandidate, 0, len(candidates))
	totalWeight := 0.0

	for _, node := range candidates {
		d := dist[current][node]
		if d <= 0 {
			d = 1e-9
		}

		tau := math.Pow(pheromone[current][node], cfg.Alpha)
		eta := math.Pow(1.0/d, cfg.Beta)
		w := tau * eta

		if math.IsNaN(w) || math.IsInf(w, 0) || w <= 0 {
			w = 1e-12
		}

		weighted = append(weighted, weightedCandidate{
			Node:   node,
			Weight: w,
		})
		totalWeight += w
	}

	if totalWeight <= 0 {
		bestNode := candidates[0]
		bestDist := dist[current][bestNode]
		for _, node := range candidates[1:] {
			if dist[current][node] < bestDist {
				bestDist = dist[current][node]
				bestNode = node
			}
		}
		return bestNode
	}

	r := rng.Float64() * totalWeight
	acc := 0.0
	for _, item := range weighted {
		acc += item.Weight
		if r <= acc {
			return item.Node
		}
	}

	return weighted[len(weighted)-1].Node
}

func evaporate(pheromone [][]float64, evaporation float64) {
	factor := 1.0 - evaporation
	minPheromone := 1e-6

	for i := 0; i < len(pheromone); i++ {
		for j := 0; j < len(pheromone[i]); j++ {
			pheromone[i][j] *= factor
			if pheromone[i][j] < minPheromone {
				pheromone[i][j] = minPheromone
			}
		}
	}
}

func depositSolution(pheromone [][]float64, solution Solution, amount float64) {
	if amount <= 0 || math.IsNaN(amount) || math.IsInf(amount, 0) {
		return
	}

	for _, route := range solution.Routes {
		prev := 0
		for _, node := range route.Nodes {
			pheromone[prev][node] += amount
			pheromone[node][prev] += amount
			prev = node
		}
		pheromone[prev][0] += amount
		pheromone[0][prev] += amount
	}
}

func makeMatrix(rows, cols int, value float64) [][]float64 {
	m := make([][]float64, rows)
	for i := 0; i < rows; i++ {
		m[i] = make([]float64, cols)
		for j := 0; j < cols; j++ {
			if i == j {
				m[i][j] = 0
			} else {
				m[i][j] = value
			}
		}
	}
	return m
}

func cloneSolution(sol Solution) Solution {
	clonedRoutes := make([]Route, len(sol.Routes))
	for i, route := range sol.Routes {
		clonedRoutes[i] = Route{
			VehicleID: route.VehicleID,
			Nodes:     append([]int{}, route.Nodes...),
		}
	}

	return Solution{
		Routes: clonedRoutes,
		Cost:   sol.Cost,
		Metrics: SearchMetrics{
			DurationMS: sol.Metrics.DurationMS,
		},
	}
}

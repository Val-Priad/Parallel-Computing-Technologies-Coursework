package main

import (
	"math"
	"math/rand"
	"time"
)

const (
	selectionNoise    = 0.05
	candidateListSize = 5
	violationPenalty  = 1000.0
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
		Beta:             2.0,
		Evaporation:      0.5,
		Q:                100.0,
		InitialPheromone: 1.0,
		EliteWeight:      1.0,
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

		for ant := 0; ant < cfg.NumAnts; ant++ {
			sol, feasible := buildAntSolution(instance, pheromone, cfg, rng)

			ants = append(ants, antSolution{
				Solution: sol,
				Feasible: feasible,
			})

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
		cfg.Beta = 2.0
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
		cfg.EliteWeight = 1.0
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
	if n == 0 {
		return Solution{Routes: []Route{}, Cost: 0}, true
	}

	demandByID := make([]int, n)
	customers := make([]int, 0, len(instance.Customers))
	seen := make([]bool, n)
	for _, c := range instance.Customers {
		if c.ID < 0 || c.ID >= n {
			return Solution{Routes: []Route{}, Cost: math.Inf(1)}, false
		}
		if seen[c.ID] {
			return Solution{Routes: []Route{}, Cost: math.Inf(1)}, false
		}
		if c.Demand > instance.VehicleCapacity {
			return Solution{Routes: []Route{}, Cost: math.Inf(1)}, false
		}

		demandByID[c.ID] = c.Demand
		seen[c.ID] = true
		customers = append(customers, c.ID)
	}

	if len(customers) == 0 {
		return Solution{Routes: []Route{}, Cost: 0}, true
	}
	unvisited := append([]int{}, customers...)
	routes := make([]Route, 0, instance.Vehicles)
	totalPenalty := 0.0

	for len(unvisited) > 0 && len(routes) < instance.Vehicles {
		routeNodes := make([]int, 0)
		load := 0
		current := 0

		for {
			candidates := feasibleCustomers(unvisited, load, instance.VehicleCapacity, demandByID)
			if len(candidates) == 0 {
				break
			}

			next := selectNextCustomer(current, candidates, instance.Dist, pheromone, cfg, rng)
			if next == -1 {
				break
			}

			routeNodes = append(routeNodes, next)
			load += demandByID[next]
			current = next
			unvisited = removeCustomerValue(unvisited, next)
		}

		routes = append(routes, Route{Nodes: routeNodes})
	}

	if len(unvisited) > 0 {
		if len(routes) == 0 {
			routes = append(routes, Route{Nodes: []int{}})
		}

		last := len(routes) - 1
		for _, node := range unvisited {
			routes[last].Nodes = append(routes[last].Nodes, node)
			totalPenalty += violationPenalty
		}
		unvisited = unvisited[:0]
	}

	for i := range routes {
		routes[i].VehicleID = i
	}

	for len(routes) < instance.Vehicles {
		routes = append(routes, Route{VehicleID: len(routes), Nodes: []int{}})
	}

	totalCost := 0.0
	for _, route := range routes {
		totalCost += computeRouteCost(route.Nodes, instance.Dist)

		load := 0
		for _, node := range route.Nodes {
			load += demandByID[node]
		}
		if load > instance.VehicleCapacity {
			totalPenalty += violationPenalty
		}
	}

	return Solution{
		Routes: routes,
		Cost:   totalCost + totalPenalty,
	}, true
}

func feasibleCustomers(unvisited []int, currentLoad, capacity int, demandByID []int) []int {
	result := make([]int, 0, len(unvisited))
	for _, node := range unvisited {
		demand := demandByID[node]
		if currentLoad+demand <= capacity {
			result = append(result, node)
		}
	}
	return result
}

func removeCustomerAt(values []int, index int) []int {
	copy(values[index:], values[index+1:])
	return values[:len(values)-1]
}

func removeCustomerValue(values []int, value int) []int {
	for index, candidate := range values {
		if candidate == value {
			return removeCustomerAt(values, index)
		}
	}

	return values
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
	candidates = nearestNeighbors(current, candidates, dist, candidateListSize)
	if len(candidates) == 0 {
		return -1
	}
	if len(candidates) == 1 {
		return candidates[0]
	}
	if rng.Float64() < 0.1 {
		return candidates[rng.Intn(len(candidates))]
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
		w *= 1.0 + selectionNoise*rng.Float64()

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

func nearestNeighbors(current int, candidates []int, dist [][]float64, k int) []int {
	if k <= 0 || len(candidates) <= k {
		return candidates
	}

	selected := make([]int, 0, k)
	used := make([]bool, len(candidates))

	for len(selected) < k {
		bestIndex := -1
		bestDist := math.Inf(1)

		for i, node := range candidates {
			if used[i] {
				continue
			}
			d := dist[current][node]
			if d < bestDist {
				bestDist = d
				bestIndex = i
			}
		}

		if bestIndex == -1 {
			break
		}

		used[bestIndex] = true
		selected = append(selected, candidates[bestIndex])
	}

	if len(selected) == 0 {
		return candidates
	}

	return selected
}

func evaporate(pheromone [][]float64, evaporation float64) {
	factor := 1.0 - evaporation

	for i := 0; i < len(pheromone); i++ {
		for j := 0; j < len(pheromone[i]); j++ {
			if i == j {
				pheromone[i][j] = 0
				continue
			}

			pheromone[i][j] *= factor
		}
	}
}

func depositSolution(pheromone [][]float64, solution Solution, amount float64) {
	if amount <= 0 || math.IsNaN(amount) || math.IsInf(amount, 0) {
		return
	}

	for _, route := range solution.Routes {
		if len(route.Nodes) == 0 {
			continue
		}

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

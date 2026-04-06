package main

import (
	"math"
	"math/rand"
	"time"
)

const (
	defaultNumAnts          = 25
	defaultIterations       = 200
	defaultAlpha            = 0.5
	defaultBeta             = 1.0
	defaultEvaporation      = 0.6
	defaultQ                = 100.0
	defaultInitialPheromone = 1.0
	defaultEliteWeight      = 1.0
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
		NumAnts:          defaultNumAnts,
		Iterations:       defaultIterations,
		Alpha:            defaultAlpha,
		Beta:             defaultBeta,
		Evaporation:      defaultEvaporation,
		Q:                defaultQ,
		InitialPheromone: defaultInitialPheromone,
		EliteWeight:      defaultEliteWeight,
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

	applyConfigDefaults(&cfg)
	if !validateInstance(instance) {
		return solutionWithMetrics([]Route{}, math.Inf(1), startTime)
	}

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
			sol, feasible := buildSolution(instance, pheromone, cfg, rng)

			ants = append(ants, antSolution{
				Solution: sol,
				Feasible: feasible,
			})

			if feasible && sol.Cost < best.Cost {
				best = cloneSolution(sol)
				stepID = logSolutionStep(logger, stepID, best)
			}
		}

		evaporate(pheromone, cfg.Evaporation)

		for _, ant := range ants {
			if !ant.Feasible || ant.Solution.Cost <= 0 {
				continue
			}
			depositSolution(pheromone, ant.Solution, cfg.Q/ant.Solution.Cost)
		}

		if !math.IsInf(best.Cost, 1) && best.Cost > 0 {
			depositSolution(pheromone, best, cfg.EliteWeight*cfg.Q/best.Cost)
		}
	}

	if math.IsInf(best.Cost, 1) {
		return solutionWithMetrics([]Route{}, math.Inf(1), startTime)
	}

	return solutionWithMetrics(best.Routes, best.Cost, startTime)
}

func applyConfigDefaults(cfg *ACOConfig) {
	if cfg.NumAnts <= 0 {
		cfg.NumAnts = defaultNumAnts
	}
	if cfg.Iterations <= 0 {
		cfg.Iterations = defaultIterations
	}
	if cfg.Alpha <= 0 {
		cfg.Alpha = defaultAlpha
	}
	if cfg.Beta <= 0 {
		cfg.Beta = defaultBeta
	}
	if cfg.Evaporation <= 0 || cfg.Evaporation >= 1 {
		cfg.Evaporation = defaultEvaporation
	}
	if cfg.Q <= 0 {
		cfg.Q = defaultQ
	}
	if cfg.InitialPheromone <= 0 {
		cfg.InitialPheromone = defaultInitialPheromone
	}
	if cfg.EliteWeight < 0 {
		cfg.EliteWeight = defaultEliteWeight
	}
	if cfg.Seed == 0 {
		cfg.Seed = time.Now().UnixNano()
	}
}

func validateInstance(instance VRPInstance) bool {
	n := len(instance.Dist)
	if n == 0 {
		return len(instance.Customers) == 0
	}

	seen := make([]bool, n)
	for _, c := range instance.Customers {
		if c.ID < 0 || c.ID >= n {
			return false
		}
		if seen[c.ID] {
			return false
		}
		if c.Demand > instance.VehicleCapacity {
			return false
		}

		seen[c.ID] = true
	}

	return true
}

func solutionWithMetrics(routes []Route, cost float64, startTime time.Time) Solution {
	return Solution{
		Routes: routes,
		Cost:   cost,
		Metrics: SearchMetrics{
			DurationMS: float64(time.Since(startTime).Nanoseconds()) / 1e6,
		},
	}
}

func logSolutionStep(logger *Logger, stepID int, sol Solution) int {
	if logger == nil {
		return stepID
	}

	loggedRoutes := make([][]int, len(sol.Routes))
	for i, route := range sol.Routes {
		loggedRoutes[i] = append([]int{}, route.Nodes...)
	}

	logger.Log(Step{
		StepID: stepID,
		Routes: loggedRoutes,
		Cost:   sol.Cost,
	})

	return stepID + 1
}

func buildSolution(
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
	for _, c := range instance.Customers {
		demandByID[c.ID] = c.Demand
		customers = append(customers, c.ID)
	}

	if len(customers) == 0 {
		return Solution{Routes: []Route{}, Cost: 0}, true
	}
	unvisited := append([]int{}, customers...)
	routes := make([]Route, 0, instance.Vehicles)

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
		return Solution{Routes: []Route{}, Cost: math.Inf(1)}, false
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
	}

	return Solution{
		Routes: routes,
		Cost:   totalCost,
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

		if w <= 0 {
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
	if amount <= 0 {
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

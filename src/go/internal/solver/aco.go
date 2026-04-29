package solver

import (
	"math"
	"math/rand"
	"parallel-aco/internal/logging"
	"parallel-aco/internal/vrp"
	"time"
)

type WorkerMemory struct {
	demandByID []int
	visited    []bool
	candidates []int
	weights    []float64
}

func NewWorkerMemory(n int, numCustomers int) *WorkerMemory {
	return &WorkerMemory{
		demandByID: make([]int, n),
		visited:    make([]bool, n),
		candidates: make([]int, 0, numCustomers),
		weights:    make([]float64, 0, numCustomers),
	}
}

const (
	defaultNumAnts          = 80
	defaultIterations       = 500
	defaultAlpha            = 1.0
	defaultBeta             = 3.5
	defaultEvaporation      = 0.15
	defaultQ                = 100.0
	defaultInitialPheromone = 0.50
	defaultEliteWeight      = 3.0
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
	Solution vrp.Solution
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

func SolveACO(instance vrp.VRPInstance, logger *logging.Logger, cfg ACOConfig) vrp.Solution {
	startTime := time.Now()

	if len(instance.Customers) == 0 {
		return solutionWithMetrics([]vrp.Route{}, 0, startTime)
	}

	if instance.Vehicles == 0 {
		return solutionWithMetrics([]vrp.Route{}, math.Inf(1), startTime)
	}

	applyConfigDefaults(&cfg)
	if !validateInstance(instance) {
		return solutionWithMetrics([]vrp.Route{}, math.Inf(1), startTime)
	}

	n := len(instance.Dist)
	pheromone := makeMatrix(n, n, cfg.InitialPheromone)
	mem := NewWorkerMemory(n, len(instance.Customers))

	rng := rand.New(rand.NewSource(cfg.Seed))

	best := vrp.Solution{
		Routes: []vrp.Route{},
		Cost:   math.Inf(1),
	}

	stepID := 0

	for iter := 0; iter < cfg.Iterations; iter++ {
		ants := make([]antSolution, 0, cfg.NumAnts)

		for ant := 0; ant < cfg.NumAnts; ant++ {
			sol, feasible := buildSolution(instance, pheromone, cfg, rng, mem)

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
		return solutionWithMetrics([]vrp.Route{}, math.Inf(1), startTime)
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

func validateInstance(instance vrp.VRPInstance) bool {
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

func solutionWithMetrics(routes []vrp.Route, cost float64, startTime time.Time) vrp.Solution {
	return vrp.Solution{
		Routes:     routes,
		Cost:       cost,
		DurationMS: float64(time.Since(startTime).Milliseconds()),
	}
}

func logSolutionStep(logger *logging.Logger, stepID int, sol vrp.Solution) int {
	if logger == nil {
		return stepID
	}

	loggedRoutes := make([][]int, len(sol.Routes))
	for i, route := range sol.Routes {
		loggedRoutes[i] = append([]int{}, route.Nodes...)
	}

	logger.Log(logging.Step{
		StepID: stepID,
		Routes: loggedRoutes,
		Cost:   sol.Cost,
	})

	return stepID + 1
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

func depositSolution(pheromone [][]float64, solution vrp.Solution, amount float64) {
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

func cloneSolution(sol vrp.Solution) vrp.Solution {
	clonedRoutes := make([]vrp.Route, len(sol.Routes))
	for i, route := range sol.Routes {
		clonedRoutes[i] = vrp.Route{
			VehicleID: route.VehicleID,
			Nodes:     append([]int{}, route.Nodes...),
		}
	}

	return vrp.Solution{
		Routes:     clonedRoutes,
		Cost:       sol.Cost,
		DurationMS: sol.DurationMS,
	}
}

func buildSolution(
	instance vrp.VRPInstance,
	pheromone [][]float64,
	cfg ACOConfig,
	rng *rand.Rand,
	mem *WorkerMemory,
) (vrp.Solution, bool) {
	n := len(instance.Dist)
	if n == 0 {
		return vrp.Solution{Routes: []vrp.Route{}, Cost: math.Inf(1)}, false
	}

	for i := 0; i < n; i++ {
		mem.demandByID[i] = 0
		mem.visited[i] = false
	}

	remaining := len(instance.Customers)
	for _, c := range instance.Customers {
		mem.demandByID[c.ID] = c.Demand
	}

	routes := make([]vrp.Route, 0, instance.Vehicles)
	routeNodes := make([]int, 0, len(instance.Customers))

	for remaining > 0 && len(routes) < instance.Vehicles {
		routeNodes = routeNodes[:0]
		load := 0
		current := 0

		for {
			mem.candidates = mem.candidates[:0]
			for _, c := range instance.Customers {
				id := c.ID
				if mem.visited[id] {
					continue
				}
				if load+mem.demandByID[id] <= instance.VehicleCapacity {
					mem.candidates = append(mem.candidates, id)
				}
			}

			if len(mem.candidates) == 0 {
				break
			}

			next := selectNextCustomer(
				current,
				mem.candidates,
				instance.Dist,
				pheromone,
				cfg,
				rng,
				&mem.weights,
			)
			if next == -1 {
				break
			}

			routeNodes = append(routeNodes, next)
			load += mem.demandByID[next]
			current = next
			mem.visited[next] = true
			remaining--
		}

		routes = append(routes, vrp.Route{
			Nodes:     append([]int(nil), routeNodes...),
			VehicleID: len(routes),
		})
	}

	if remaining > 0 {
		return vrp.Solution{Routes: []vrp.Route{}, Cost: math.Inf(1)}, false
	}

	for len(routes) < instance.Vehicles {
		routes = append(routes, vrp.Route{VehicleID: len(routes), Nodes: []int{}})
	}

	totalCost := 0.0
	for _, route := range routes {
		totalCost += computeRouteCost(route.Nodes, instance.Dist)
	}

	return vrp.Solution{Routes: routes, Cost: totalCost}, true
}

func selectNextCustomer(
	current int,
	candidates []int,
	dist [][]float64,
	pheromone [][]float64,
	cfg ACOConfig,
	rng *rand.Rand,
	weights *[]float64,
) int {
	if len(candidates) == 0 {
		return -1
	}
	if len(candidates) == 1 {
		return candidates[0]
	}

	w := *weights
	if cap(w) < len(candidates) {
		w = make([]float64, len(candidates))
	}
	w = w[:len(candidates)]
	*weights = w

	totalWeight := 0.0
	for i, node := range candidates {
		d := dist[current][node]
		if d <= 0 {
			d = 1e-9
		}
		weight := math.Pow(pheromone[current][node], cfg.Alpha) * math.Pow(1.0/d, cfg.Beta)
		if weight <= 0 {
			weight = 1e-12
		}
		w[i] = weight
		totalWeight += weight
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
	for i, weight := range w {
		acc += weight
		if r <= acc {
			return candidates[i]
		}
	}

	return candidates[len(candidates)-1]
}

func computeRouteCost(nodes []int, dist [][]float64) float64 {
	if len(nodes) == 0 {
		return 0
	}

	cost := 0.0
	prev := 0

	for _, node := range nodes {
		cost += dist[prev][node]
		prev = node
	}

	cost += dist[prev][0]
	return cost
}

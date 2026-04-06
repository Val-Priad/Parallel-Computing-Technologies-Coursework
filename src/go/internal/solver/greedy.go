package solver

import (
	"math"
	"parallel-aco/internal/logging"
	"parallel-aco/internal/vrp"
	"time"
)

func SolveGreedy(instance vrp.VRPInstance, logger *logging.Logger) vrp.Solution {
	n := len(instance.Customers)
	k := instance.Vehicles
	startTime := time.Now()

	if n == 0 || k == 0 {
		return vrp.Solution{
			Routes: []vrp.Route{},
			Cost:   0,
			Metrics: vrp.SearchMetrics{
				DurationMS: float64(time.Since(startTime).Nanoseconds()) / 1e6,
			},
		}
	}

	dist := instance.Dist
	capacity := instance.VehicleCapacity

	maxID := 0
	for _, customer := range instance.Customers {
		if customer.ID > maxID {
			maxID = customer.ID
		}
	}

	visited := make([]bool, maxID+1)
	demandByID := make([]int, maxID+1)
	for _, customer := range instance.Customers {
		if customer.ID >= 0 && customer.ID < len(demandByID) {
			demandByID[customer.ID] = customer.Demand
		}
	}
	routes := make([]vrp.Route, k)
	current := make([]int, k)
	currentLoad := make([]int, k)

	for i := 0; i < k; i++ {
		current[i] = 0
		routes[i] = vrp.Route{
			VehicleID: i,
			Nodes:     []int{},
		}
	}

	remaining := n

	for remaining > 0 {
		bestVehicle := -1
		bestNode := -1
		bestScore := math.Inf(1)
		bestDist := math.Inf(1)

		for v := 0; v < k; v++ {
			for _, c := range instance.Customers {
				if c.ID < 0 || c.ID >= len(visited) || visited[c.ID] {
					continue
				}
				if currentLoad[v]+demandByID[c.ID] > capacity {
					continue
				}

				delta := dist[current[v]][c.ID] + dist[c.ID][0] - dist[current[v]][0]
				candidateDist := dist[current[v]][c.ID]

				if delta < bestScore || (delta == bestScore && candidateDist < bestDist) {
					bestScore = delta
					bestDist = candidateDist
					bestVehicle = v
					bestNode = c.ID
				}
			}
		}

		if bestVehicle == -1 || bestNode == -1 {
			break
		}

		routes[bestVehicle].Nodes = append(routes[bestVehicle].Nodes, bestNode)
		visited[bestNode] = true
		currentLoad[bestVehicle] += demandByID[bestNode]
		current[bestVehicle] = bestNode
		remaining--
	}

	totalCost := 0.0
	for _, route := range routes {
		totalCost += computeRouteCost(route.Nodes, dist)
	}

	if logger != nil {
		loggedRoutes := make([][]int, len(routes))
		for i, route := range routes {
			loggedRoutes[i] = append([]int{}, route.Nodes...)
		}

		logger.Log(logging.Step{
			StepID: 0,
			Routes: loggedRoutes,
			Cost:   totalCost,
		})
	}

	return vrp.Solution{
		Routes: routes,
		Cost:   totalCost,
		Metrics: vrp.SearchMetrics{
			DurationMS: float64(time.Since(startTime).Nanoseconds()) / 1e6,
		},
	}
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

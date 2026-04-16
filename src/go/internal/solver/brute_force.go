package solver

import (
	"math"
	"parallel-aco/internal/logging"
	"parallel-aco/internal/vrp"
	"time"
)

func SolveBruteForce(instance vrp.VRPInstance, logger *logging.Logger) vrp.Solution {
	n := len(instance.Customers)
	startTime := time.Now()

	if n == 0 {
		return vrp.Solution{
			Routes:     []vrp.Route{},
			Cost:       0,
			DurationMS: max(0.001, float64(time.Since(startTime).Milliseconds())),
		}
	}

	customerIDs := make([]int, n)
	for i, c := range instance.Customers {
		customerIDs[i] = c.ID
	}

	best := vrp.Solution{Cost: math.Inf(1)}
	dist := instance.Dist
	capacity := instance.VehicleCapacity
	k := instance.Vehicles
	stepID := 0

	demandByID := make([]int, len(dist))
	for _, customer := range instance.Customers {
		demandByID[customer.ID] = customer.Demand
	}

	permute(customerIDs, func(order []int) {
		findBestPartition(order, k, capacity, demandByID, dist, logger, &stepID, &best)
	})

	best.DurationMS = max(0.001, float64(time.Since(startTime).Milliseconds()))

	return best
}

func findBestPartition(order []int, numVehicles int, capacity int, demand []int, dist [][]float64, logger *logging.Logger, stepID *int, best *vrp.Solution) {
	n := len(order)
	if n == 0 {
		return
	}

	assignment := make([]int, n)
	findPartitionRec(order, assignment, 0, numVehicles, capacity, demand, dist, logger, stepID, best)
}

func findPartitionRec(order []int, assignment []int, customerIdx int, numVehicles int, capacity int, demand []int, dist [][]float64, logger *logging.Logger, stepID *int, best *vrp.Solution) {
	if customerIdx == len(order) {
		evaluatePartition(order, assignment, numVehicles, capacity, demand, dist, logger, stepID, best)
		return
	}

	for vehicle := 0; vehicle < numVehicles; vehicle++ {
		assignment[customerIdx] = vehicle
		findPartitionRec(order, assignment, customerIdx+1, numVehicles, capacity, demand, dist, logger, stepID, best)
	}
}

func evaluatePartition(order []int, assignment []int, numVehicles int, capacity int, demand []int, dist [][]float64, logger *logging.Logger, stepID *int, best *vrp.Solution) {
	routes := make([]vrp.Route, numVehicles)
	vehicleLoads := make([]int, numVehicles)

	for i, customerID := range order {
		vehicle := assignment[i]
		routes[vehicle].Nodes = append(routes[vehicle].Nodes, customerID)
		routes[vehicle].VehicleID = vehicle
		vehicleLoads[vehicle] += demand[customerID]
	}

	for _, load := range vehicleLoads {
		if load > capacity {
			return
		}
	}

	totalCost := 0.0
	for _, route := range routes {
		cost := routeCostWithReturn(route.Nodes, dist)
		totalCost += cost
	}

	if totalCost < best.Cost {
		best.Cost = totalCost
		best.Routes = copyRoutes(routes)
		if stepID != nil {
			*stepID = logSolutionStep(logger, *stepID, *best)
		}
	}
}

func routeCostWithReturn(route []int, dist [][]float64) float64 {
	if len(route) == 0 {
		return 0
	}
	cost := 0.0
	prev := 0
	for _, node := range route {
		cost += dist[prev][node]
		prev = node
	}
	cost += dist[prev][0]
	return cost
}

func copyRoutes(routes []vrp.Route) []vrp.Route {
	result := make([]vrp.Route, len(routes))
	for i, r := range routes {
		result[i] = vrp.Route{
			VehicleID: r.VehicleID,
			Nodes:     append([]int{}, r.Nodes...),
		}
	}
	return result
}

func permute(arr []int, callback func([]int)) {
	if len(arr) == 0 {
		callback([]int{})
		return
	}

	var generate func(int)
	generate = func(n int) {
		if n == 1 {
			tmp := make([]int, len(arr))
			copy(tmp, arr)
			callback(tmp)
			return
		}
		for i := 0; i < n; i++ {
			generate(n - 1)
			if n%2 == 1 {
				arr[0], arr[n-1] = arr[n-1], arr[0]
			} else {
				arr[i], arr[n-1] = arr[n-1], arr[i]
			}
		}
	}
	generate(len(arr))
}

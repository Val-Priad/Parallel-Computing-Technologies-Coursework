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
			Routes:  []vrp.Route{},
			Cost:    0,
			Metrics: vrp.SearchMetrics{DurationMS: float64(time.Since(startTime).Nanoseconds()) / 1e6},
		}
	}

	perm := make([]int, n-1)
	for i := 0; i < n-1; i++ {
		perm[i] = i + 2
	}

	best := vrp.Solution{Cost: math.Inf(1)}

	stepID := 0

	permute(perm, func(p []int) {
		order := make([]int, n)
		order[0] = 1
		copy(order[1:], p)

		sol := evaluate(instance, order)
		if sol.Cost < best.Cost {
			loggedRoutes := make([][]int, len(sol.Routes))
			for i, route := range sol.Routes {
				loggedRoutes[i] = append([]int{}, route.Nodes...)
			}

			best = sol

			logger.Log(logging.Step{
				StepID: stepID,
				Routes: loggedRoutes,
				Cost:   sol.Cost,
			})
			stepID++

		}
	})

	best.Metrics = vrp.SearchMetrics{
		DurationMS: float64(time.Since(startTime).Nanoseconds()) / 1e6,
	}

	return best
}

func permute(arr []int, f func([]int)) {
	if len(arr) == 0 {
		f([]int{})
		return
	}

	var generate func(int)
	generate = func(n int) {
		if n == 1 {
			tmp := make([]int, len(arr))
			copy(tmp, arr)
			f(tmp)
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

func evaluate(instance vrp.VRPInstance, order []int) vrp.Solution {
	k := instance.Vehicles
	n := len(order)
	if n == 0 || k == 0 {
		return vrp.Solution{Routes: []vrp.Route{}, Cost: 0}
	}

	best := vrp.Solution{Cost: math.Inf(1), Routes: []vrp.Route{}}
	dist := instance.Dist
	capacity := instance.VehicleCapacity
	demandByID := make([]int, len(dist))
	for _, customer := range instance.Customers {
		if customer.ID >= 0 && customer.ID < len(demandByID) {
			demandByID[customer.ID] = customer.Demand
		}
	}
	routeEnds := make([]int, k)

	var buildRoutesFromEnds = func() []vrp.Route {
		routes := make([]vrp.Route, k)
		segmentStart := 0
		for i := 0; i < k; i++ {
			segmentEnd := routeEnds[i]
			nodes := make([]int, segmentEnd-segmentStart)
			copy(nodes, order[segmentStart:segmentEnd])
			routes[i] = vrp.Route{
				VehicleID: i,
				Nodes:     nodes,
			}
			segmentStart = segmentEnd
		}
		return routes
	}

	var search func(vehicleIdx, start int, currentCost float64, prevFirstNode int)
	search = func(vehicleIdx, start int, currentCost float64, prevFirstNode int) {
		if currentCost >= best.Cost {
			return
		}

		if vehicleIdx == k {
			if start == n {
				if currentCost < best.Cost {
					best.Cost = currentCost
					best.Routes = buildRoutesFromEnds()
				}
			}
			return
		}

		if start >= n {
			for i := vehicleIdx; i < k; i++ {
				routeEnds[i] = start
			}
			if currentCost < best.Cost {
				best.Cost = currentCost
				best.Routes = buildRoutesFromEnds()
			}
			return
		}

		prevNode := 0
		routeCostWithoutReturn := 0.0
		routeLoad := 0
		for end := start; end <= n; end++ {
			if end == start {
				if start != 0 {
					continue
				}
				routeEnds[vehicleIdx] = end
				search(vehicleIdx+1, end, currentCost, prevFirstNode)
				continue
			}

			firstNode := order[start]
			node := order[end-1]

			routeLoad += demandByID[node]
			if routeLoad > capacity {
				break
			}

			if end-start > 1 && firstNode > node {
				continue
			}

			if prevFirstNode != 0 && firstNode < prevFirstNode {
				continue
			}

			nextRouteCostWithoutReturn := routeCostWithoutReturn + dist[prevNode][node]
			nextPrevNode := node

			routeCost := nextRouteCostWithoutReturn + dist[nextPrevNode][0]
			nextCost := currentCost + routeCost
			if nextCost >= best.Cost {
				continue
			}

			routeCostWithoutReturn = nextRouteCostWithoutReturn
			prevNode = nextPrevNode

			routeEnds[vehicleIdx] = end
			search(vehicleIdx+1, end, nextCost, firstNode)
		}
	}

	search(0, 0, 0, 0)

	if math.IsInf(best.Cost, 1) {
		return vrp.Solution{
			Routes: []vrp.Route{},
			Cost:   math.Inf(1),
		}
	}
	return best
}

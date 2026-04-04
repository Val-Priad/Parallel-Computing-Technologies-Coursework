package main

import (
	"math"
	"time"
)

func SolveVRP(instance VRPInstance, logger *Logger) Solution {
	n := len(instance.Customers)
	startTime := time.Now()

	if n == 0 {
		return Solution{
			Routes:  []Route{},
			Cost:    0,
			Metrics: SearchMetrics{DurationMS: float64(time.Since(startTime).Microseconds()) / 1000.0},
		}
	}

	perm := make([]int, n-1)
	for i := 0; i < n-1; i++ {
		perm[i] = i + 2
	}

	best := Solution{Cost: math.Inf(1)}

	stepID := 0

	Permute(perm, func(p []int) {
		order := make([]int, n)
		order[0] = 1
		copy(order[1:], p)

		sol := Evaluate(instance, order)
		if sol.Cost < best.Cost {
			loggedRoutes := make([][]int, len(sol.Routes))
			for i, route := range sol.Routes {
				loggedRoutes[i] = append([]int{}, route.Nodes...)
			}

			best = sol

			logger.Log(Step{
				StepID: stepID,
				Routes: loggedRoutes,
				Cost:   sol.Cost,
			})
			stepID++

		}
	})

	best.Metrics = SearchMetrics{
		DurationMS: float64(time.Since(startTime).Microseconds()) / 1000.0,
	}

	return best
}

func Permute(arr []int, f func([]int)) {
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

func Evaluate(instance VRPInstance, order []int) Solution {
	k := instance.Vehicles
	n := len(order)
	if n == 0 || k == 0 {
		return Solution{Routes: []Route{}, Cost: 0}
	}

	best := Solution{Cost: math.Inf(1), Routes: []Route{}}
	dist := instance.Dist
	capacity := instance.VehicleCapacity
	demandByID := make([]int, len(dist))
	for _, customer := range instance.Customers {
		if customer.ID >= 0 && customer.ID < len(demandByID) {
			demandByID[customer.ID] = customer.Demand
		}
	}
	routeEnds := make([]int, k)

	var buildRoutesFromEnds = func() []Route {
		routes := make([]Route, k)
		segmentStart := 0
		for i := 0; i < k; i++ {
			segmentEnd := routeEnds[i]
			nodes := make([]int, segmentEnd-segmentStart)
			copy(nodes, order[segmentStart:segmentEnd])
			routes[i] = Route{
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
		return Solution{
			Routes: []Route{},
			Cost:   math.Inf(1),
		}
	}
	return best
}

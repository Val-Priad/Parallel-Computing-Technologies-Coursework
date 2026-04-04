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

		remainingVehicles := k - vehicleIdx
		remainingCustomers := n - start
		maxEnd := n
		if remainingCustomers >= remainingVehicles {
			maxEnd = n - (remainingVehicles - 1)
		}

		prevNode := 0
		routeCostWithoutReturn := 0.0
		for end := start + 1; end <= maxEnd; end++ {
			firstNode := order[start]
			node := order[end-1]

			if end-start > 1 && firstNode > node {
				continue
			}

			if prevFirstNode != 0 && firstNode < prevFirstNode {
				continue
			}

			routeCostWithoutReturn += dist[prevNode][node]
			prevNode = node

			routeCost := routeCostWithoutReturn + dist[prevNode][0]
			nextCost := currentCost + routeCost
			if nextCost >= best.Cost {
				continue
			}

			routeEnds[vehicleIdx] = end
			search(vehicleIdx+1, end, nextCost, firstNode)
		}
	}

	search(0, 0, 0, 0)

	if math.IsInf(best.Cost, 1) {
		best = Solution{Routes: []Route{}, Cost: 0}
	}
	return best
}

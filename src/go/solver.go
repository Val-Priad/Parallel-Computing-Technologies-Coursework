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
	totalChecked := 0

	stepID := 0

	Permute(perm, func(p []int) {
		order := make([]int, n)
		order[0] = 1
		copy(order[1:], p)

		sol := Evaluate(instance, order)
		totalChecked += sol.Metrics.CheckedSolutions
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
		CheckedSolutions: totalChecked,
		DurationMS:       float64(time.Since(startTime).Microseconds()) / 1000.0,
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

	best := Solution{Cost: math.Inf(1)}
	checkedSolutions := 0

	var search func(vehicleIdx, start int, currentCost float64, routes []Route)
	search = func(vehicleIdx, start int, currentCost float64, routes []Route) {
		if currentCost >= best.Cost {
			return
		}

		if vehicleIdx == k {
			if start == n {
				checkedSolutions++
				candidate := Solution{
					Routes: cloneRoutes(routes),
					Cost:   currentCost,
				}
				if candidate.Cost < best.Cost {
					best = candidate
				}
			}
			return
		}

		if start >= n {
			finalRoutes := cloneRoutes(routes)
			for i := vehicleIdx; i < k; i++ {
				finalRoutes = append(finalRoutes, Route{VehicleID: i, Nodes: []int{}})
			}
			checkedSolutions++
			candidate := Solution{
				Routes: finalRoutes,
				Cost:   currentCost,
			}
			if candidate.Cost < best.Cost {
				best = candidate
			}
			return
		}

		remainingVehicles := k - vehicleIdx
		remainingCustomers := n - start
		maxEnd := n
		if remainingCustomers >= remainingVehicles {
			maxEnd = n - (remainingVehicles - 1)
		}

		for end := start + 1; end <= maxEnd; end++ {
			nodes := append([]int{}, order[start:end]...)
			routeCost := RouteCost(instance, nodes)
			nextRoutes := append(cloneRoutes(routes), Route{
				VehicleID: vehicleIdx,
				Nodes:     nodes,
			})
			search(vehicleIdx+1, end, currentCost+routeCost, nextRoutes)
		}
	}

	search(0, 0, 0, []Route{})

	if math.IsInf(best.Cost, 1) {
		best = Solution{Routes: []Route{}, Cost: 0}
	}

	best.Metrics.CheckedSolutions = checkedSolutions
	return best
}

func cloneRoutes(routes []Route) []Route {
	cloned := make([]Route, len(routes))
	for i, route := range routes {
		nodes := append([]int{}, route.Nodes...)
		cloned[i] = Route{
			VehicleID: route.VehicleID,
			Nodes:     nodes,
		}
	}
	return cloned
}

func RouteCost(instance VRPInstance, nodes []int) float64 {
	if len(nodes) == 0 {
		return 0
	}

	depot := 0
	cost := 0.0

	prev := depot

	for _, node := range nodes {
		cost += instance.Dist[prev][node]
		prev = node
	}

	cost += instance.Dist[prev][depot]

	return cost
}

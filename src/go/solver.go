package main

import (
	"math"
)

func SolveVRP(instance VRPInstance, logger *Logger) Solution {
	n := len(instance.Customers)

	perm := make([]int, n)
	for i := 0; i < n; i++ {
		perm[i] = i + 1 // ID клиентов
	}

	best := Solution{Cost: math.Inf(1)}

	stepID := 0

	Permute(perm, func(p []int) {
		sol := Evaluate(instance, p)

		loggedRoutes := make([][]int, len(sol.Routes))
		for i, route := range sol.Routes {
			loggedRoutes[i] = append([]int{}, route.Nodes...)
		}

		logger.Log(Step{
			StepID: stepID,
			Routes: loggedRoutes,
			Cost:   sol.Cost,
		})

		stepID++

		if sol.Cost < best.Cost {
			best = sol
		}
	})

	return best
}

func Permute(arr []int, f func([]int)) {
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

	chunkSize := int(math.Ceil(float64(n) / float64(k)))

	routes := []Route{}
	totalCost := 0.0

	for i := 0; i < k; i++ {
		start := i * chunkSize
		end := (i + 1) * chunkSize
		if end > n {
			end = n
		}
		if start >= n {
			break
		}

		nodes := order[start:end]
		cost := RouteCost(instance, nodes)

		routes = append(routes, Route{
			VehicleID: i,
			Nodes:     nodes,
		})

		totalCost += cost
	}

	return Solution{
		Routes: routes,
		Cost:   totalCost,
	}
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

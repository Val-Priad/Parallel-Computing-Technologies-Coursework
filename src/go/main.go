package main

import "fmt"

func main() {
	points := []Point{
		{0, 0, 0},
		{1, 1, 3},
		{2, 4, 4},
		{3, 6, 1},
		{4, 3, 7},
	}

	dist := BuildDistanceMatrix(points)

	instance := VRPInstance{
		Depot:     points[0],
		Customers: points[1:],
		Vehicles:  2,
		Dist:      dist,
	}

	logger := NewLogger(true)

	solution := SolveVRP(instance, logger)

	fmt.Println("Best cost:", solution.Cost)
	fmt.Println("Routes:", solution.Routes)
	fmt.Println("Checked solutions:", solution.Metrics.CheckedSolutions)
	fmt.Printf("Execution time: %.3f ms\n", solution.Metrics.DurationMS)

	logger.SaveToFile("log.json", points, solution.Metrics)
}

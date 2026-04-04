package main

import (
	"fmt"
	"path/filepath"
)

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
	fmt.Printf("Execution time: %.3f ms\n", solution.Metrics.DurationMS)

	exactLogName := filepath.Join("logs", fmt.Sprintf("brute-force__vehicles-%d_customers-%d.json", instance.Vehicles, len(instance.Customers)))
	if err := logger.SaveToFile(exactLogName, points, solution.Metrics); err != nil {
		fmt.Println("Failed to save exact log:", err)
	}

	greedyLogger := NewLogger(true)
	greedySolution := SolveGreedy(instance, greedyLogger)

	fmt.Println("\n--- Greedy ---")
	fmt.Println("Cost:", greedySolution.Cost)
	fmt.Println("Routes:", greedySolution.Routes)
	fmt.Printf("Execution time: %.3f ms\n", greedySolution.Metrics.DurationMS)

	greedyLogName := filepath.Join("logs", fmt.Sprintf("greedy__vehicles-%d_customers-%d.json", instance.Vehicles, len(instance.Customers)))
	if err := greedyLogger.SaveToFile(greedyLogName, points, greedySolution.Metrics); err != nil {
		fmt.Println("Failed to save greedy log:", err)
	}
}

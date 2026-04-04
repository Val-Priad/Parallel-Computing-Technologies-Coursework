package main

import (
	"fmt"
	"path/filepath"
)

func main() {
	points, instance := GenerateInstance(GeneratorConfig{
		NumCustomers: 10,
		Vehicles:     3,
		Width:        50,
		Height:       50,
		Seed:         42,
	})

	logger := NewLogger(true)

	solution := SolveVRP(instance, logger)

	fmt.Println("Best cost:", solution.Cost)
	fmt.Println("Routes:", solution.Routes)
	fmt.Printf("Execution time: %.3f ms\n", solution.Metrics.DurationMS)

	exactLogName := filepath.Join("logs",
		fmt.Sprintf("brute-force__vehicles-%d_customers-%d.json",
			instance.Vehicles, len(instance.Customers)),
	)
	logger.SaveToFile(exactLogName, points, solution.Metrics)

	greedyLogger := NewLogger(true)
	greedySolution := SolveGreedy(instance, greedyLogger)

	fmt.Println("\n--- Greedy ---")
	fmt.Println("Cost:", greedySolution.Cost)
	fmt.Println("Routes:", greedySolution.Routes)
	fmt.Printf("Execution time: %.3f ms\n", greedySolution.Metrics.DurationMS)

	greedyLogName := filepath.Join("logs",
		fmt.Sprintf("greedy__vehicles-%d_customers-%d.json",
			instance.Vehicles, len(instance.Customers)),
	)
	greedyLogger.SaveToFile(greedyLogName, points, greedySolution.Metrics)
}

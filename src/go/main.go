package main

import (
	"fmt"
	"path/filepath"
)

func main() {
	experiments := GetExperiments()

	csvLogger, err := NewCSVLogger("results.csv")
	if err != nil {
		panic(err)
	}
	defer csvLogger.Close()

	for _, exp := range experiments {
		baseSeed := exp.Seed
		for i := 0; i < 1; i++ {
			runSeed := baseSeed + int64(i)
			runName := fmt.Sprintf("%s_run_%d", exp.Name, i+1)

			fmt.Println("\n=== Running", runName, "(seed", runSeed, ")===")

			points, instance := GenerateInstance(GeneratorConfig{
				NumCustomers: exp.NumCustomers,
				Vehicles:     exp.Vehicles,
				Width:        exp.Width,
				Height:       exp.Height,
				Seed:         runSeed,
				CapacityMode: exp.CapacityMode,
			})

			greedyLogger := NewLogger(true)
			greedySolution := SolveGreedy(instance, greedyLogger)

			fmt.Println("Greedy:", greedySolution.Cost)
			greedyLogPath := filepath.Join("logs", fmt.Sprintf("%s_greedy.json", runName))
			if err := greedyLogger.SaveToFile(greedyLogPath, points, greedySolution.Metrics); err != nil {
				fmt.Println("Failed to save greedy log:", err)
			}

			csvLogger.Log(runName, "greedy", instance, greedySolution)

			bruteForceLogger := NewLogger(true)
			exactSolution := SolveBruteForce(instance, bruteForceLogger)

			fmt.Println("Brute:", exactSolution.Cost)
			exactLogPath := filepath.Join("logs", fmt.Sprintf("%s_brute_force.json", runName))
			if err := bruteForceLogger.SaveToFile(exactLogPath, points, exactSolution.Metrics); err != nil {
				fmt.Println("Failed to save brute-force log:", err)
			}

			csvLogger.Log(runName, "brute_force", instance, exactSolution)
		}
	}

	fmt.Println("\nDone. Results saved to results.csv")
}

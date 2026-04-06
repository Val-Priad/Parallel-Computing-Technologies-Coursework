package main

import "fmt"

func main() {
	experiments := GetExperiments()

	csvLogger, err := NewCSVLogger("results.csv")
	if err != nil {
		panic(err)
	}
	defer csvLogger.Close()

	for _, exp := range experiments {
		baseSeed := exp.Seed
		for i := 0; i < 5; i++ {
			runSeed := baseSeed + int64(i)
			runName := fmt.Sprintf("%s_run_%d", exp.Name, i+1)

			fmt.Println("\n=== Running", runName, "(seed", runSeed, ")===")

			_, instance := GenerateInstance(GeneratorConfig{
				NumCustomers: exp.NumCustomers,
				Vehicles:     exp.Vehicles,
				Width:        exp.Width,
				Height:       exp.Height,
				Seed:         runSeed,
				CapacityMode: exp.CapacityMode,
			})

			greedyLogger := NewLogger(false)
			greedySolution := SolveGreedy(instance, greedyLogger)

			fmt.Println("Greedy:", greedySolution.Cost)

			csvLogger.Log(runName, "greedy", instance, greedySolution)

			exactLogger := NewLogger(false)
			exactSolution := SolveVRP(instance, exactLogger)

			fmt.Println("Brute:", exactSolution.Cost)

			csvLogger.Log(runName, "brute_force", instance, exactSolution)
		}
	}

	fmt.Println("\nDone. Results saved to results.csv")
}

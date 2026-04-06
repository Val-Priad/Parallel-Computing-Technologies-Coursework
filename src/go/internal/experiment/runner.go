package experiment

import (
	"fmt"
	"parallel-aco/internal/logging"
	"parallel-aco/internal/solver"
	"parallel-aco/internal/vrp"
	"path/filepath"
)

func Run() {
	experiments := GetExperiments()

	csvLogger, err := logging.NewCSVLogger("results.csv")
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

			points, instance := vrp.GenerateInstance(vrp.GeneratorConfig{
				NumCustomers: exp.NumCustomers,
				Vehicles:     exp.Vehicles,
				Width:        exp.Width,
				Height:       exp.Height,
				Seed:         runSeed,
				CapacityMode: exp.CapacityMode,
			})

			greedyLogger := logging.NewLogger(true)
			greedySolution := solver.SolveGreedy(instance, greedyLogger)

			fmt.Println("Greedy:", greedySolution.Cost)
			greedyLogPath := filepath.Join("logs", fmt.Sprintf("%s_greedy.json", runName))
			if err := greedyLogger.SaveToFile(greedyLogPath, points, greedySolution.Metrics); err != nil {
				fmt.Println("Failed to save greedy log:", err)
			}

			if exp.NumCustomers <= 12 {
				csvLogger.Log(runName, "greedy", instance, greedySolution)

				bruteForceLogger := logging.NewLogger(true)
				exactSolution := solver.SolveBruteForce(instance, bruteForceLogger)

				fmt.Println("Brute:", exactSolution.Cost)
				exactLogPath := filepath.Join("logs", fmt.Sprintf("%s_brute_force.json", runName))
				if err := bruteForceLogger.SaveToFile(exactLogPath, points, exactSolution.Metrics); err != nil {
					fmt.Println("Failed to save brute-force log:", err)
				}

				csvLogger.Log(runName, "brute_force", instance, exactSolution)
			}

			acoLogger := logging.NewLogger(true)
			acoSolution := solver.SolveACO(instance, acoLogger, solver.DefaultACOConfig())

			fmt.Println("ACO:", acoSolution.Cost)
			acoLogPath := filepath.Join("logs", fmt.Sprintf("%s_aco.json", runName))
			if err := acoLogger.SaveToFile(acoLogPath, points, acoSolution.Metrics); err != nil {
				fmt.Println("Failed to save ACO log:", err)
			}

			csvLogger.Log(runName, "aco", instance, acoSolution)
		}
	}

	fmt.Println("\nDone. Results saved to results.csv")
}

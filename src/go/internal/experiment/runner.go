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

	logSolution := func(runName, algo string, logger *logging.Logger, points []vrp.Point, metrics vrp.SearchMetrics) {
		logPath := filepath.Join("logs", fmt.Sprintf("%s_%s.json", runName, algo))
		if err := logger.SaveToFile(logPath, points, metrics); err != nil {
			fmt.Printf("Failed to save %s log: %v\n", algo, err)
		}
	}

	for expIdx, exp := range experiments {
		for i := 0; i < 1; i++ {
			runSeed := exp.Seed + int64(i)
			runName := fmt.Sprintf(
				"exp_%d_c%d_v%d_w%.0f_h%.0f_run_%d",
				expIdx+1,
				exp.NumCustomers,
				exp.Vehicles,
				exp.Width,
				exp.Height,
				i+1,
			)

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
			logSolution(runName, "greedy", greedyLogger, points, greedySolution.Metrics)
			csvLogger.Log(runName, "greedy", instance, greedySolution)

			acoLogger := logging.NewLogger(true)
			acoSolution := solver.SolveACO(instance, acoLogger, solver.DefaultACOConfig())
			fmt.Println("ACO:", acoSolution.Cost)
			logSolution(runName, "aco", acoLogger, points, acoSolution.Metrics)
			csvLogger.Log(runName, "aco", instance, acoSolution)

			if exp.NumCustomers <= 8 {
				bruteForceLogger := logging.NewLogger(true)
				exactSolution := solver.SolveBruteForce(instance, bruteForceLogger)
				fmt.Println("Brute:", exactSolution.Cost)
				logSolution(runName, "brute_force", bruteForceLogger, points, exactSolution.Metrics)
				csvLogger.Log(runName, "brute_force", instance, exactSolution)
			}
		}
	}

	fmt.Println("\nDone. Results saved to results.csv")
}

package experiment

import (
	"fmt"
	"os"
	"parallel-aco/internal/logging"
	"parallel-aco/internal/solver"
	"parallel-aco/internal/vrp"
	"path/filepath"
)

func RunACOExperiment() {
	experiments := GetLargeComparisonExperiments()

	if err := os.MkdirAll(resultsDir, 0o755); err != nil {
		panic(err)
	}

	csvLogger, err := logging.NewCSVLogger(filepath.Join(resultsDir, "aco.csv"))
	if err != nil {
		panic(err)
	}
	defer csvLogger.Close()

	fmt.Println("Running ACO experiment")

	for expIdx, exp := range experiments {
		runSeed := exp.Seed
		runName := fmt.Sprintf("aco_exp_%d_c%d_v%d_w%.0f_h%.0f", expIdx+1, exp.NumCustomers, exp.Vehicles, exp.Width, exp.Height)

		fmt.Println("\n=== Running", runName, "(seed", runSeed, ")===")

		points, instance := vrp.GenerateInstance(vrp.GeneratorConfig{
			NumCustomers: exp.NumCustomers,
			Vehicles:     exp.Vehicles,
			Width:        exp.Width,
			Height:       exp.Height,
			Seed:         runSeed,
			CapacityMode: exp.CapacityMode,
		})

		acoLogger := logging.NewLogger(false)
		acoSolution := solver.SolveACO(instance, acoLogger, solver.DefaultACOConfig())
		writeExperimentResult(csvLogger, runName, "aco", instance, acoSolution, acoLogger, points)
	}

	fmt.Println("\nDone.")
}

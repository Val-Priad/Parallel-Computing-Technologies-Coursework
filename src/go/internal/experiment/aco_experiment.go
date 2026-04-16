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
	experiments := GetLargeComparisonExperiments(10)

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

		acoSolution := solver.SolveACO(instance, nil, solver.DefaultACOConfig())
		writeExperimentResult(csvLogger, runName, "aco", "aco", instance, acoSolution, nil, points)
	}

	fmt.Println("\nDone.")
}

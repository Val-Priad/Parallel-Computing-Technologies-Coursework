package experiment

import (
	"fmt"
	"os"
	"parallel-aco/internal/logging"
	"parallel-aco/internal/solver"
	"parallel-aco/internal/vrp"
	"path/filepath"
)

func RunACOvsPACO() {
	experiments := GetLargeComparisonExperiments()

	if err := os.MkdirAll(resultsDir, 0o755); err != nil {
		panic(err)
	}

	csvLogger, err := logging.NewCSVLogger(filepath.Join(resultsDir, "aco_vs_paco.csv"))
	if err != nil {
		panic(err)
	}
	defer csvLogger.Close()

	fmt.Println("Running ACO vs PACO comparison")

	for expIdx, exp := range experiments {
		runSeed := exp.Seed
		runName := fmt.Sprintf("paco_exp_%d_c%d_v%d_w%.0f_h%.0f", expIdx+1, exp.NumCustomers, exp.Vehicles, exp.Width, exp.Height)

		fmt.Println("\n=== Running", runName, "(seed", runSeed, ")===")

		points, instance := vrp.GenerateInstance(vrp.GeneratorConfig{
			NumCustomers: exp.NumCustomers,
			Vehicles:     exp.Vehicles,
			Width:        exp.Width,
			Height:       exp.Height,
			Seed:         runSeed,
			CapacityMode: exp.CapacityMode,
		})

		acoLogger := logging.NewLogger(true)
		acoSolution := solver.SolveACO(instance, acoLogger, solver.DefaultACOConfig())

		pacoSolution := solver.SolvePACO(instance, solver.PACOConfig{BaseConfig: solver.DefaultACOConfig()})

		writeExperimentResult(csvLogger, runName, "aco", instance, acoSolution, acoLogger, points)
		writeExperimentResult(csvLogger, runName, "paco", instance, pacoSolution, nil, points)
	}

	fmt.Println("\nDone.")
}

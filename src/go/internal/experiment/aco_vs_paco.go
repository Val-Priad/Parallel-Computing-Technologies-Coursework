package experiment

import (
	"fmt"
	"os"
	"parallel-aco/internal/logging"
	"parallel-aco/internal/solver"
	"parallel-aco/internal/vrp"
	"path/filepath"
)

const pacoResultsDir = "results"

func RunACOvsPACO() {
	experiments := GetACOvsPACOExperiments()

	if err := os.MkdirAll(pacoResultsDir, 0o755); err != nil {
		panic(err)
	}

	csvLogger, err := logging.NewCSVLogger(filepath.Join(pacoResultsDir, "aco_vs_paco.csv"))
	if err != nil {
		panic(err)
	}
	defer csvLogger.Close()

	fmt.Println("Running ACO vs PACO comparison")

	logSolution := func(runName, algo string, logger *logging.Logger, points []vrp.Point, durationMS float64) {
		if logger == nil {
			return
		}

		logPath := filepath.Join(pacoResultsDir, fmt.Sprintf("%s_%s.json", runName, algo))
		if err := logger.SaveToFile(logPath, points, durationMS); err != nil {
			fmt.Printf("Failed to save %s log: %v\n", algo, err)
		}
	}

	logAndWrite := func(
		runName string,
		algo string,
		instance vrp.VRPInstance,
		solution vrp.Solution,
		logger *logging.Logger,
		points []vrp.Point,
	) {
		fmt.Printf("%s: cost=%.3f, time=%.3fms\n", algo, solution.Cost, solution.DurationMS)
		logSolution(runName, algo, logger, points, solution.DurationMS)
		csvLogger.Log(runName, algo, instance, solution)
	}

	for expIdx, exp := range experiments {
		for i := 0; i < 1; i++ {
			runSeed := exp.Seed + int64(i)
			runName := fmt.Sprintf(
				"paco_exp_%d_c%d_v%d_w%.0f_h%.0f_run_%d",
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

			acoLogger := logging.NewLogger(true)
			acoSolution := solver.SolveACO(instance, acoLogger, solver.DefaultACOConfig())

			pacoSolution := solver.SolvePACO(instance, solver.PACOConfig{BaseConfig: solver.DefaultACOConfig()})

			logAndWrite(runName, "aco", instance, acoSolution, acoLogger, points)
			logAndWrite(runName, "paco", instance, pacoSolution, nil, points)
		}
	}

	fmt.Println("\nDone.")
}

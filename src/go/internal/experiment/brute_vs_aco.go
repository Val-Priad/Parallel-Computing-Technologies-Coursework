package experiment

import (
	"fmt"
	"os"
	"parallel-aco/internal/logging"
	"parallel-aco/internal/solver"
	"parallel-aco/internal/vrp"
	"path/filepath"
)

const (
	resultsDir = "results"
)

func saveSolutionLog(runName, algo string, logger *logging.Logger, points []vrp.Point, durationMS float64) {
	if logger == nil {
		return
	}

	logPath := filepath.Join(resultsDir, fmt.Sprintf("%s_%s.json", runName, algo))
	if err := logger.SaveToFile(logPath, points, durationMS); err != nil {
		fmt.Printf("Failed to save %s log: %v\n", algo, err)
	}
}

func logResult(
	csvLogger *logging.CSVLogger,
	runName string,
	algo string,
	instance vrp.VRPInstance,
	solution vrp.Solution,
	logger *logging.Logger,
	points []vrp.Point,
) {
	fmt.Printf("%s: cost=%.3f, time=%.3fms\n", algo, solution.Cost, solution.DurationMS)
	saveSolutionLog(runName, algo, logger, points, solution.DurationMS)
	csvLogger.Log(runName, algo, instance, solution)
}

func RunBruteForceVsACO() {
	experiments := GetExperiments()

	if err := os.MkdirAll(resultsDir, 0o755); err != nil {
		panic(err)
	}

	csvLogger, err := logging.NewCSVLogger(filepath.Join(resultsDir, "brute_force_vs_aco.csv"))
	if err != nil {
		panic(err)
	}
	defer csvLogger.Close()

	fmt.Println("Running brute force vs ACO comparison")

	for expIdx, exp := range experiments {
		runSeed := exp.Seed + int64(expIdx)
		runName := fmt.Sprintf("exp_%d_c%d_v%d_w%.0f_h%.0f", expIdx+1, exp.NumCustomers, exp.Vehicles, exp.Width, exp.Height)

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
		logResult(csvLogger, runName, "aco", instance, acoSolution, acoLogger, points)

		bruteForceLogger := logging.NewLogger(true)
		exactSolution := solver.SolveBruteForce(instance, bruteForceLogger)
		logResult(csvLogger, runName, "brute_force", instance, exactSolution, bruteForceLogger, points)
	}

}

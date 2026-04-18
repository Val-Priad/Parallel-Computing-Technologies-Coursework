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
	experiments := GetLargeComparisonExperiments(5)
	const runsPerExperiment = 20
	const warmupRuns = 3

	if len(experiments) == 0 {
		fmt.Println("No experiments configured.")
		return
	}

	if err := os.MkdirAll(resultsDir, 0o755); err != nil {
		panic(err)
	}

	csvLogger, err := logging.NewCSVLogger(filepath.Join(resultsDir, "aco_vs_paco.csv"))
	if err != nil {
		panic(err)
	}
	defer csvLogger.Close()

	fmt.Println("Running ACO vs PACO comparison")

	firstExp := experiments[0]
	_, firstInstance := vrp.GenerateInstance(vrp.GeneratorConfig{
		NumCustomers: firstExp.NumCustomers,
		Vehicles:     firstExp.Vehicles,
		Width:        firstExp.Width,
		Height:       firstExp.Height,
		Seed:         firstExp.Seed,
		CapacityMode: firstExp.CapacityMode,
	})

	warmupACOAndPACO(firstInstance, firstExp.Seed, warmupRuns)

	for expIdx, exp := range experiments {
		runSeed := exp.Seed
		runName := fmt.Sprintf("aco_vs_paco_exp_%d_c%d_v%d_w%.0f_h%.0f_runs%d", expIdx+1, exp.NumCustomers, exp.Vehicles, exp.Width, exp.Height, runsPerExperiment)

		fmt.Println("\n=== Running", runName, "(base seed", runSeed, ")===")

		points, instance := vrp.GenerateInstance(vrp.GeneratorConfig{
			NumCustomers: exp.NumCustomers,
			Vehicles:     exp.Vehicles,
			Width:        exp.Width,
			Height:       exp.Height,
			Seed:         runSeed,
			CapacityMode: exp.CapacityMode,
		})

		acoTotalCost := 0.0
		acoTotalDuration := 0.0
		pacoTotalCost := 0.0
		pacoTotalDuration := 0.0

		for run := 0; run < runsPerExperiment; run++ {
			seed := runSeed + int64(run)

			acoCfg := solver.DefaultACOConfig()
			acoCfg.Seed = seed
			acoSolution := solver.SolveACO(instance, nil, acoCfg)
			acoTotalCost += acoSolution.Cost
			acoTotalDuration += acoSolution.DurationMS

			pacoCfg := solver.PACOConfig{BaseConfig: solver.DefaultACOConfig()}
			pacoCfg.BaseConfig.Seed = seed
			pacoSolution := solver.SolvePACO(instance, pacoCfg)
			pacoTotalCost += pacoSolution.Cost
			pacoTotalDuration += pacoSolution.DurationMS
		}

		avgAcoSolution := vrp.Solution{
			Cost:       acoTotalCost / float64(runsPerExperiment),
			DurationMS: acoTotalDuration / float64(runsPerExperiment),
		}
		avgPacoSolution := vrp.Solution{
			Cost:       pacoTotalCost / float64(runsPerExperiment),
			DurationMS: pacoTotalDuration / float64(runsPerExperiment),
		}

		writeExperimentResult(csvLogger, runName, "aco", "aco_vs_paco", instance, avgAcoSolution, nil, points)
		writeExperimentResult(csvLogger, runName, "paco", "aco_vs_paco", instance, avgPacoSolution, nil, points)
	}

	fmt.Println("\nDone.")
}

func warmupACOAndPACO(instance vrp.VRPInstance, seed int64, warmupRuns int) {
	for warmup := 0; warmup < warmupRuns; warmup++ {
		currentSeed := seed + int64(warmup)

		acoCfg := solver.DefaultACOConfig()
		acoCfg.Seed = currentSeed
		_ = solver.SolveACO(instance, nil, acoCfg)

		pacoCfg := solver.PACOConfig{BaseConfig: solver.DefaultACOConfig()}
		pacoCfg.BaseConfig.Seed = currentSeed
		_ = solver.SolvePACO(instance, pacoCfg)
	}
}

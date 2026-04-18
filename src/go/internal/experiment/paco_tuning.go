package experiment

import (
	"encoding/csv"
	"fmt"
	"os"
	"parallel-aco/internal/solver"
	"parallel-aco/internal/vrp"
	"path/filepath"
	"strconv"
)

func RunPACOProcessTuning() {
	experiments := GetLargeComparisonExperiments(1)

	exp := experiments[0]

	if err := os.MkdirAll(resultsDir, 0o755); err != nil {
		panic(err)
	}

	file, err := os.Create(filepath.Join(resultsDir, "paco_process_tuning.csv"))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write([]string{
		"workers",
		"time_ms",
		"cost",
	}); err != nil {
		panic(err)
	}

	fmt.Println("Running PACO process scaling")

	_, instance := vrp.GenerateInstance(vrp.GeneratorConfig{
		NumCustomers: exp.NumCustomers,
		Vehicles:     exp.Vehicles,
		Width:        exp.Width,
		Height:       exp.Height,
		Seed:         exp.Seed,
		CapacityMode: exp.CapacityMode,
	})

	baseCfg := solver.DefaultACOConfig()

	for workers := 1; workers <= 20; workers++ {
		fmt.Printf("workers=%d...\n", workers)

		runCfg := solver.PACOConfig{
			BaseConfig: baseCfg,
			NumWorkers: workers,
		}
		runCfg.BaseConfig.Seed = exp.Seed

		solution := solver.SolvePACO(instance, runCfg)

		if err := writer.Write([]string{
			strconv.Itoa(workers),
			fmt.Sprintf("%.3f", solution.DurationMS),
			fmt.Sprintf("%.3f", solution.Cost),
		}); err != nil {
			panic(err)
		}
	}

	fmt.Println("\nDone.")
}

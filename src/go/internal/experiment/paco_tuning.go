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

type pacoProcessTuningRow struct {
	Experiment        string
	Workers           int
	AverageCost       float64
	AverageDurationMS float64
	Selected          bool
}

const pacoProcessTuningRunsPerSetting = 3

func RunPACOProcessTuning() {
	experiments := GetLargeComparisonExperiments(4)

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
		"experiment",
		"customers",
		"vehicles",
		"capacity",
		"capacity_mode",
		"workers",
		"avg_time_ms",
		"avg_cost",
		"selected",
	}); err != nil {
		panic(err)
	}

	fmt.Println("Running PACO process tuning")

	for expIdx, exp := range experiments {
		runSeed := exp.Seed
		runName := fmt.Sprintf("paco_tuning_exp_%d_c%d_v%d_w%.0f_h%.0f", expIdx+1, exp.NumCustomers, exp.Vehicles, exp.Width, exp.Height)

		fmt.Println("\n=== Running", runName, "(seed", runSeed, ")===")

		_, instance := vrp.GenerateInstance(vrp.GeneratorConfig{
			NumCustomers: exp.NumCustomers,
			Vehicles:     exp.Vehicles,
			Width:        exp.Width,
			Height:       exp.Height,
			Seed:         runSeed,
			CapacityMode: exp.CapacityMode,
		})

		rows := evaluatePACOProcessRows(runName, exp, instance)
		bestIdx := selectBestPACOProcessRow(rows)
		if bestIdx >= 0 {
			rows[bestIdx].Selected = true
		}

		for _, row := range rows {
			if err := writer.Write([]string{
				row.Experiment,
				strconv.Itoa(exp.NumCustomers),
				strconv.Itoa(exp.Vehicles),
				strconv.Itoa(instance.VehicleCapacity),
				string(exp.CapacityMode),
				strconv.Itoa(row.Workers),
				fmt.Sprintf("%.3f", row.AverageDurationMS),
				fmt.Sprintf("%.3f", row.AverageCost),
				strconv.FormatBool(row.Selected),
			}); err != nil {
				panic(err)
			}
		}

		if bestIdx >= 0 {
			best := rows[bestIdx]
			fmt.Printf("selected workers=%d avg_cost=%.3f avg_time=%.3fms\n",
				best.Workers,
				best.AverageCost,
				best.AverageDurationMS,
			)
		}
	}

	fmt.Println("\nDone.")
}

func evaluatePACOProcessRows(runName string, exp ExperimentConfig, instance vrp.VRPInstance) []pacoProcessTuningRow {
	processCounts := candidatePACOProcessCounts()
	rows := make([]pacoProcessTuningRow, 0, len(processCounts))

	baseCfg := solver.DefaultACOConfig()
	for _, workers := range processCounts {
		row := pacoProcessTuningRow{
			Experiment: runName,
			Workers:    workers,
		}

		totalCost := 0.0
		totalDuration := 0.0

		for trial := 0; trial < pacoProcessTuningRunsPerSetting; trial++ {
			runCfg := solver.PACOConfig{BaseConfig: baseCfg, NumWorkers: workers}
			runCfg.BaseConfig.Seed = exp.Seed + int64(trial)

			solution := solver.SolvePACO(instance, runCfg)
			totalCost += solution.Cost

			totalDuration += solution.DurationMS
		}

		row.AverageCost = totalCost / float64(pacoProcessTuningRunsPerSetting)
		row.AverageDurationMS = totalDuration / float64(pacoProcessTuningRunsPerSetting)

		rows = append(rows, row)
	}

	return rows
}

func selectBestPACOProcessRow(rows []pacoProcessTuningRow) int {
	bestIdx := -1
	for i := range rows {
		if bestIdx < 0 || pacoProcessRowLess(rows[i], rows[bestIdx]) {
			bestIdx = i
		}
	}

	return bestIdx
}

func pacoProcessRowLess(a, b pacoProcessTuningRow) bool {
	if a.AverageCost != b.AverageCost {
		return a.AverageCost < b.AverageCost
	}
	if a.AverageDurationMS != b.AverageDurationMS {
		return a.AverageDurationMS < b.AverageDurationMS
	}
	return a.Workers < b.Workers
}

func candidatePACOProcessCounts() []int {
	// maxWorkers := runtime.NumCPU()
	maxWorkers := 15
	baseAnts := solver.DefaultACOConfig().NumAnts
	if baseAnts > 0 && baseAnts < maxWorkers {
		maxWorkers = baseAnts
	}
	if maxWorkers < 1 {
		maxWorkers = 1
	}

	counts := make([]int, 0, maxWorkers)
	for workers := 1; workers <= maxWorkers; workers++ {
		counts = append(counts, workers)
	}

	return counts
}

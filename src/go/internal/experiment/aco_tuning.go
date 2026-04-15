package experiment

import (
	"fmt"
	"math"
	"sort"
	"sync"

	"parallel-aco/internal/solver"
	"parallel-aco/internal/vrp"
)

type acoTuningCase struct {
	exp      ExperimentConfig
	instance vrp.VRPInstance
}

type acoTuningRow struct {
	Config            solver.ACOConfig
	FeasibleRuns      int
	TotalRuns         int
	AverageCost       float64
	AverageDurationMS float64
}

const acoTuningWorkers = 6
const infeasibleRunPenaltyCost = 1e15

func TuneACOConfig(experiments []ExperimentConfig) (solver.ACOConfig, []acoTuningRow) {

	cases := buildACOOnceCases(experiments)
	if len(cases) == 0 {
		return solver.DefaultACOConfig(), nil
	}

	rows := evaluateACOConfigs(cases, candidateACOConfigs())
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].FeasibleRuns != rows[j].FeasibleRuns {
			return rows[i].FeasibleRuns > rows[j].FeasibleRuns
		}
		if rows[i].AverageCost != rows[j].AverageCost {
			return rows[i].AverageCost < rows[j].AverageCost
		}
		return rows[i].AverageDurationMS < rows[j].AverageDurationMS
	})

	return rows[0].Config, rows
}

func RunACOTuning() {
	experiments := GetLargeComparisonExperiments()
	bestConfig, rows := TuneACOConfig(experiments)
	fmt.Printf("(sequential ACO, %d workers in config pool)\n", acoTuningWorkers)
	printTuningResults("ACO (sequential, pooled)", bestConfig, rows)
}

func printTuningResults(label string, bestConfig solver.ACOConfig, rows []acoTuningRow) {
	fmt.Printf("\n=== %s parameter tuning on large instances ===\n", label)
	fmt.Printf("%-4s %-48s %12s %12s %14s %12s\n",
		"rank",
		"config",
		"feasible",
		"avg_cost",
		"avg_ms",
		"runs",
	)

	limit := 5
	if len(rows) < limit {
		limit = len(rows)
	}

	for i := 0; i < limit; i++ {
		row := rows[i]
		fmt.Printf("%-4d %-48s %6d/%-5d %12.3f %14.3f %12d\n",
			i+1,
			formatACOConfig(row.Config),
			row.FeasibleRuns,
			row.TotalRuns,
			row.AverageCost,
			row.AverageDurationMS,
			row.TotalRuns,
		)
	}

	fmt.Printf("\nselected config: %s\n", formatACOConfig(bestConfig))
}

func buildACOOnceCases(experiments []ExperimentConfig) []acoTuningCase {
	cases := make([]acoTuningCase, 0, len(experiments))

	for _, exp := range experiments {
		_, instance := vrp.GenerateInstance(vrp.GeneratorConfig{
			NumCustomers: exp.NumCustomers,
			Vehicles:     exp.Vehicles,
			Width:        exp.Width,
			Height:       exp.Height,
			Seed:         exp.Seed,
			CapacityMode: exp.CapacityMode,
		})

		cases = append(cases, acoTuningCase{
			exp:      exp,
			instance: instance,
		})
	}

	return cases
}

func evaluateACOConfigs(cases []acoTuningCase, configs []solver.ACOConfig) []acoTuningRow {
	if len(configs) == 0 {
		return nil
	}

	workers := acoTuningWorkers
	if workers > len(configs) {
		workers = len(configs)
	}
	if workers < 1 {
		workers = 1
	}

	jobs := make(chan solver.ACOConfig, len(configs))
	results := make(chan acoTuningRow, len(configs))

	var wg sync.WaitGroup
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			for job := range jobs {
				results <- evaluateSingleACOConfig(cases, job)
			}
		}()
	}

	for _, cfg := range configs {
		jobs <- cfg
	}
	close(jobs)

	wg.Wait()
	close(results)

	rows := make([]acoTuningRow, 0, len(configs))
	for row := range results {
		rows = append(rows, row)
	}

	return rows
}

func evaluateSingleACOConfig(cases []acoTuningCase, cfg solver.ACOConfig) acoTuningRow {
	row := acoTuningRow{Config: cfg}
	totalCost := 0.0
	totalDuration := 0.0
	for _, tc := range cases {
		for trial := 0; trial < 3; trial++ {
			runCfg := cfg
			runCfg.Seed = tc.exp.Seed + int64(trial)

			solution := solver.SolveACO(tc.instance, nil, runCfg)

			if isFiniteFeasibleCost(solution.Cost) {
				row.FeasibleRuns++
				totalCost += solution.Cost
			} else {
				totalCost += infeasibleRunPenaltyCost
			}

			totalDuration += solution.DurationMS
			row.TotalRuns++
		}
	}

	if row.TotalRuns > 0 {
		row.AverageCost = totalCost / float64(row.TotalRuns)
		row.AverageDurationMS = totalDuration / float64(row.TotalRuns)
	}

	return row
}

func candidateACOConfigs() []solver.ACOConfig {
	base := solver.DefaultACOConfig()

	configs := make([]solver.ACOConfig, 0, 7)
	seen := make(map[string]struct{})

	addConfig := func(cfg solver.ACOConfig) {
		cfg.Seed = 0
		key := acoConfigKey(cfg)
		if _, ok := seen[key]; ok {
			return
		}

		seen[key] = struct{}{}
		configs = append(configs, cfg)
	}

	presets := []solver.ACOConfig{
		{
			NumAnts:          80,
			Iterations:       500,
			Alpha:            1.0,
			Beta:             3.5,
			Evaporation:      0.15,
			Q:                100.0,
			InitialPheromone: 0.50,
			EliteWeight:      3.0,
		},
		{
			NumAnts:          60,
			Iterations:       250,
			Alpha:            1.0,
			Beta:             3.0,
			Evaporation:      0.20,
			Q:                base.Q,
			InitialPheromone: base.InitialPheromone,
			EliteWeight:      base.EliteWeight,
		},
		{
			NumAnts:          90,
			Iterations:       220,
			Alpha:            0.9,
			Beta:             4.2,
			Evaporation:      0.15,
			Q:                base.Q,
			InitialPheromone: base.InitialPheromone,
			EliteWeight:      base.EliteWeight,
		},
		{
			NumAnts:          40,
			Iterations:       360,
			Alpha:            1.2,
			Beta:             3.5,
			Evaporation:      0.25,
			Q:                base.Q,
			InitialPheromone: base.InitialPheromone,
			EliteWeight:      base.EliteWeight,
		},
		{
			NumAnts:          110,
			Iterations:       180,
			Alpha:            0.8,
			Beta:             4.8,
			Evaporation:      0.18,
			Q:                base.Q,
			InitialPheromone: base.InitialPheromone,
			EliteWeight:      base.EliteWeight,
		},
		{
			NumAnts:          70,
			Iterations:       300,
			Alpha:            1.1,
			Beta:             3.8,
			Evaporation:      0.22,
			Q:                base.Q,
			InitialPheromone: base.InitialPheromone,
			EliteWeight:      base.EliteWeight,
		},
		{
			NumAnts:          120,
			Iterations:       150,
			Alpha:            1.0,
			Beta:             4.5,
			Evaporation:      0.12,
			Q:                base.Q,
			InitialPheromone: base.InitialPheromone,
			EliteWeight:      base.EliteWeight,
		},
	}

	for _, cfg := range presets {
		addConfig(cfg)
	}

	return configs
}

func acoConfigKey(cfg solver.ACOConfig) string {
	return fmt.Sprintf(
		"ants=%d;iter=%d;a=%.3f;b=%.3f;evap=%.3f;q=%.3f;init=%.3f;elite=%.3f",
		cfg.NumAnts,
		cfg.Iterations,
		cfg.Alpha,
		cfg.Beta,
		cfg.Evaporation,
		cfg.Q,
		cfg.InitialPheromone,
		cfg.EliteWeight,
	)
}

func formatACOConfig(cfg solver.ACOConfig) string {
	return fmt.Sprintf(
		"ants=%d iter=%d a=%.1f b=%.1f evap=%.2f init=%.2f elite=%.1f",
		cfg.NumAnts,
		cfg.Iterations,
		cfg.Alpha,
		cfg.Beta,
		cfg.Evaporation,
		cfg.InitialPheromone,
		cfg.EliteWeight,
	)
}

func isFiniteFeasibleCost(cost float64) bool {
	return !math.IsInf(cost, 0) && !math.IsNaN(cost) && cost > 0
}

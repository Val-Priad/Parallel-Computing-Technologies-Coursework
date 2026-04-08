package experiment

import (
	"fmt"
	"math"
	"sort"

	"parallel-aco/internal/solver"
	"parallel-aco/internal/vrp"
)

type acoTuningCase struct {
	exp        ExperimentConfig
	instance   vrp.VRPInstance
	greedyCost float64
}

type acoTuningRow struct {
	Config            solver.ACOConfig
	FeasibleRuns      int
	TotalRuns         int
	AverageCost       float64
	AverageDurationMS float64
	AverageGapPct     float64
}

func TuneACOConfig(experiments []ExperimentConfig, numWorkers int) (solver.ACOConfig, []acoTuningRow, float64) {
	if numWorkers <= 0 {
		numWorkers = 4
	}

	cases, averageGreedyCost := buildACOOnceCases(experiments)
	if len(cases) == 0 {
		return solver.DefaultACOConfig(), nil, 0
	}

	rows := evaluateACOConfigs(cases, averageGreedyCost, candidateACOConfigs(), numWorkers)
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].FeasibleRuns != rows[j].FeasibleRuns {
			return rows[i].FeasibleRuns > rows[j].FeasibleRuns
		}
		if rows[i].AverageCost != rows[j].AverageCost {
			return rows[i].AverageCost < rows[j].AverageCost
		}
		return rows[i].AverageDurationMS < rows[j].AverageDurationMS
	})

	return rows[0].Config, rows, averageGreedyCost
}

func RunACOTuning() {
	experiments := GetLargeComparisonExperiments()
	numWorkers := 12
	bestConfig, rows, averageGreedyCost := TuneACOConfig(experiments, numWorkers)
	fmt.Printf("(parallel PACO, %d workers per candidate)\n", numWorkers)
	printTuningResults("ACO (parallel)", bestConfig, rows, averageGreedyCost)
}

func printTuningResults(label string, bestConfig solver.ACOConfig, rows []acoTuningRow, averageGreedyCost float64) {
	fmt.Printf("\n=== %s parameter tuning on large instances ===\n", label)
	fmt.Printf("average greedy cost: %.3f\n", averageGreedyCost)
	fmt.Printf("%-4s %-48s %12s %12s %14s %14s %12s\n",
		"rank",
		"config",
		"feasible",
		"avg_cost",
		"gap_vs_greedy",
		"avg_ms",
		"runs",
	)

	limit := 5
	if len(rows) < limit {
		limit = len(rows)
	}

	for i := 0; i < limit; i++ {
		row := rows[i]
		fmt.Printf("%-4d %-48s %6d/%-5d %12.3f %12s %14.3f %12d\n",
			i+1,
			formatACOConfig(row.Config),
			row.FeasibleRuns,
			row.TotalRuns,
			row.AverageCost,
			formatGap(row.AverageGapPct),
			row.AverageDurationMS,
			row.TotalRuns,
		)
	}

	fmt.Printf("\nselected config: %s\n", formatACOConfig(bestConfig))
}

func buildACOOnceCases(experiments []ExperimentConfig) ([]acoTuningCase, float64) {
	cases := make([]acoTuningCase, 0, len(experiments))
	totalGreedyCost := 0.0

	for _, exp := range experiments {
		_, instance := vrp.GenerateInstance(vrp.GeneratorConfig{
			NumCustomers: exp.NumCustomers,
			Vehicles:     exp.Vehicles,
			Width:        exp.Width,
			Height:       exp.Height,
			Seed:         exp.Seed,
			CapacityMode: exp.CapacityMode,
		})

		greedySolution := solver.SolveGreedy(instance, nil)
		totalGreedyCost += greedySolution.Cost

		cases = append(cases, acoTuningCase{
			exp:        exp,
			instance:   instance,
			greedyCost: greedySolution.Cost,
		})
	}

	averageGreedyCost := 0.0
	if len(cases) > 0 {
		averageGreedyCost = totalGreedyCost / float64(len(cases))
	}

	return cases, averageGreedyCost
}

func evaluateACOConfigs(cases []acoTuningCase, averageGreedyCost float64, configs []solver.ACOConfig, numWorkers int) []acoTuningRow {
	rows := make([]acoTuningRow, 0, len(configs))

	for _, cfg := range configs {
		row := acoTuningRow{Config: cfg}
		totalCost := 0.0
		totalDuration := 0.0
		totalRuns := 0

		for _, tc := range cases {
			for trial := 0; trial < 3; trial++ {
				runCfg := cfg
				runCfg.Seed = tc.exp.Seed + int64(trial)

				pacoCfg := solver.PACOConfig{BaseConfig: runCfg, NumWorkers: numWorkers}
				solution := solver.SolvePACO(tc.instance, pacoCfg)

				if isFiniteFeasibleCost(solution.Cost) {
					row.FeasibleRuns++
					totalCost += solution.Cost
				} else {
					totalCost += tc.greedyCost * 4.0
				}

				totalDuration += solution.Metrics.DurationMS
				totalRuns++
			}
		}

		row.TotalRuns = totalRuns
		if totalRuns > 0 {
			row.AverageCost = totalCost / float64(totalRuns)
			row.AverageDurationMS = totalDuration / float64(totalRuns)
		}
		if averageGreedyCost > 0 {
			row.AverageGapPct = (row.AverageCost/averageGreedyCost - 1.0) * 100.0
		}

		rows = append(rows, row)
	}

	return rows
}

func candidateACOConfigs() []solver.ACOConfig {
	base := solver.DefaultACOConfig()
	base.Seed = 0

	return []solver.ACOConfig{
		base,
		{
			NumAnts:          20,
			Iterations:       150,
			Alpha:            0.5,
			Beta:             1.0,
			Evaporation:      0.60,
			Q:                100.0,
			InitialPheromone: 1.0,
			EliteWeight:      1.0,
		},
		{
			NumAnts:          40,
			Iterations:       200,
			Alpha:            0.8,
			Beta:             2.0,
			Evaporation:      0.35,
			Q:                100.0,
			InitialPheromone: 1.0,
			EliteWeight:      1.0,
		},
		{
			NumAnts:          60,
			Iterations:       250,
			Alpha:            1.0,
			Beta:             3.0,
			Evaporation:      0.30,
			Q:                100.0,
			InitialPheromone: 0.50,
			EliteWeight:      1.5,
		},
		{
			NumAnts:          80,
			Iterations:       300,
			Alpha:            1.0,
			Beta:             4.0,
			Evaporation:      0.25,
			Q:                100.0,
			InitialPheromone: 0.25,
			EliteWeight:      2.0,
		},
		{
			NumAnts:          100,
			Iterations:       350,
			Alpha:            1.2,
			Beta:             3.5,
			Evaporation:      0.20,
			Q:                100.0,
			InitialPheromone: 0.25,
			EliteWeight:      2.0,
		},
		{
			NumAnts:          120,
			Iterations:       400,
			Alpha:            1.0,
			Beta:             4.5,
			Evaporation:      0.20,
			Q:                100.0,
			InitialPheromone: 0.10,
			EliteWeight:      2.5,
		},
		{
			NumAnts:          60,
			Iterations:       400,
			Alpha:            1.5,
			Beta:             2.5,
			Evaporation:      0.40,
			Q:                100.0,
			InitialPheromone: 1.0,
			EliteWeight:      1.0,
		},
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
	}
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

func formatGap(gapPct float64) string {
	if math.IsNaN(gapPct) || math.IsInf(gapPct, 0) {
		return "n/a"
	}

	return fmt.Sprintf("%+.1f%%", gapPct)
}

func isFiniteFeasibleCost(cost float64) bool {
	return !math.IsInf(cost, 0) && !math.IsNaN(cost) && cost > 0
}

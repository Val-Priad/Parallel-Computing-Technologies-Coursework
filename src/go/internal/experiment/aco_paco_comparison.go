package experiment

import (
	"fmt"
	"math"
	"parallel-aco/internal/solver"
	"parallel-aco/internal/vrp"
)

type comparisonRow struct {
	ID             string
	Customers      int
	Vehicles       int
	ACOCost        float64
	ACODurationMS  float64
	PACOCost       float64
	PACODurationMS float64
}

func RunACOvsPACOComparison() {
	experiments := GetLargeComparisonExperiments()

	fmt.Println("\n=== ACO vs PACO comparison on large instances ===")
	fmt.Printf("%-18s %10s %10s %14s %14s %14s %14s %12s %10s\n",
		"experiment",
		"customers",
		"vehicles",
		"aco_cost",
		"aco_ms",
		"paco_cost",
		"paco_ms",
		"delta_cost",
		"speedup",
	)

	for _, exp := range experiments {
		row := runComparisonExperiment(exp)
		deltaCost := row.PACOCost - row.ACOCost
		speedupText := formatSpeedup(row.ACODurationMS, row.PACODurationMS)

		fmt.Printf("%-18s %10d %10d %14.3f %14.3f %14.3f %14.3f %12.3f %10s\n",
			row.ID,
			row.Customers,
			row.Vehicles,
			row.ACOCost,
			row.ACODurationMS,
			row.PACOCost,
			row.PACODurationMS,
			deltaCost,
			speedupText,
		)
	}
}

func runComparisonExperiment(exp ExperimentConfig) comparisonRow {
	_, instance := vrp.GenerateInstance(vrp.GeneratorConfig{
		NumCustomers: exp.NumCustomers,
		Vehicles:     exp.Vehicles,
		Width:        exp.Width,
		Height:       exp.Height,
		Seed:         exp.Seed,
		CapacityMode: exp.CapacityMode,
	})

	acoCfg := solver.DefaultACOConfig()
	acoCfg.Seed = exp.Seed
	acoCfg.NumAnts = 90
	acoCfg.Iterations = 220
	acoCfg.Alpha = 0.9
	acoCfg.Beta = 4.2
	acoCfg.Evaporation = 0.15
	acoCfg.InitialPheromone = 0.50
	acoCfg.EliteWeight = 3.0

	acoSolution := solver.SolveACO(instance, nil, acoCfg)

	pacoCfg := solver.PACOConfig{
		BaseConfig: acoCfg,
		NumWorkers: 12,
	}
	pacoSolution := solver.SolvePACO(instance, pacoCfg)

	return comparisonRow{
		ID:             fmt.Sprintf("c%d_v%d_s%d", exp.NumCustomers, exp.Vehicles, exp.Seed),
		Customers:      exp.NumCustomers,
		Vehicles:       exp.Vehicles,
		ACOCost:        acoSolution.Cost,
		PACOCost:       pacoSolution.Cost,
		ACODurationMS:  acoSolution.DurationMS,
		PACODurationMS: pacoSolution.DurationMS,
	}
}

func formatSpeedup(acoDurationMS, pacoDurationMS float64) string {
	if pacoDurationMS <= 0 || math.IsInf(pacoDurationMS, 0) || math.IsNaN(pacoDurationMS) {
		return "n/a"
	}

	speedup := acoDurationMS / pacoDurationMS
	if math.IsInf(speedup, 0) || math.IsNaN(speedup) {
		return "n/a"
	}

	return fmt.Sprintf("%.2fx", speedup)
}

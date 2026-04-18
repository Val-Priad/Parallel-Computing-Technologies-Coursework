package solver

import (
	"math"
	"reflect"
	"testing"

	"parallel-aco/internal/logging"
	"parallel-aco/internal/vrp"
)

const costEpsilon = 1e-9

type solverCase struct {
	name  string
	solve func(vrp.VRPInstance) vrp.Solution
}

type instanceCase struct {
	name  string
	build func() vrp.VRPInstance
}

type solverResult struct {
	name     string
	solution vrp.Solution
}

func TestSolvers_SmallInstances_CompareAgainstBruteForceOracle(t *testing.T) {
	for _, tc := range smallInstanceCases() {
		t.Run(tc.name, func(t *testing.T) {
			instance := tc.build()

			oracle := solveBruteForceOracle(t, instance)
			results := runAllSolvers(instance)

			assertAllSolutionsValid(t, instance, results)
			assertNoSolverBeatsOracle(t, oracle.Cost, results)
		})
	}
}

func TestSolvers_ComplexInstances_FeasibilityAndBruteForceBaseline(t *testing.T) {
	for _, tc := range complexInstanceCases() {
		t.Run(tc.name, func(t *testing.T) {
			instance := tc.build()

			oracle := solveBruteForceOracle(t, instance)
			results := runAllSolvers(instance)

			assertAllSolutionsValid(t, instance, results)
			assertNoSolverBeatsOracle(t, oracle.Cost, results)
		})
	}
}

func TestSolvers_Determinism_ACO(t *testing.T) {
	instance := generatedInstance(vrp.GeneratorConfig{
		NumCustomers: 5,
		Vehicles:     2,
		Width:        30,
		Height:       30,
		Seed:         301,
		MinDemand:    1,
		MaxDemand:    4,
		CapacityMode: vrp.CapacityAuto,
	})()

	cfg := ACOConfig{
		NumAnts:    12,
		Iterations: 12,
		Seed:       4242,
	}

	first := SolveACO(instance, nil, cfg)
	second := SolveACO(instance, nil, cfg)

	assertValidSolution(t, instance, first)
	assertValidSolution(t, instance, second)

	if diff := math.Abs(first.Cost - second.Cost); diff > costEpsilon {
		t.Fatalf(
			"ACO should be deterministic for the same seed: cost mismatch %.12f vs %.12f",
			first.Cost,
			second.Cost,
		)
	}

	if !reflect.DeepEqual(first.Routes, second.Routes) {
		t.Fatalf(
			"ACO should be deterministic for the same seed: routes differ: %+v vs %+v",
			first.Routes,
			second.Routes,
		)
	}
}

func allSolverCases() []solverCase {
	baseACOConfig := ACOConfig{
		NumAnts:    20,
		Iterations: 20,
		Seed:       42,
	}

	return []solverCase{
		{
			name: "aco",
			solve: func(instance vrp.VRPInstance) vrp.Solution {
				return SolveACO(instance, nil, baseACOConfig)
			},
		},
		{
			name: "paco",
			solve: func(instance vrp.VRPInstance) vrp.Solution {
				return SolvePACO(instance, PACOConfig{
					BaseConfig: baseACOConfig,
					NumWorkers: 2,
				})
			},
		},
	}
}

func smallInstanceCases() []instanceCase {
	return []instanceCase{
		{
			name:  "empty",
			build: emptyInstance,
		},
		{
			name: "small tight capacity",
			build: generatedInstance(vrp.GeneratorConfig{
				NumCustomers: 3,
				Vehicles:     1,
				Width:        20,
				Height:       20,
				Seed:         101,
				MinDemand:    1,
				MaxDemand:    3,
				CapacityMode: vrp.CapacityTight,
			}),
		},
		{
			name: "small two vehicles",
			build: generatedInstance(vrp.GeneratorConfig{
				NumCustomers: 4,
				Vehicles:     2,
				Width:        25,
				Height:       25,
				Seed:         102,
				MinDemand:    1,
				MaxDemand:    4,
				CapacityMode: vrp.CapacityAuto,
			}),
		},
		{
			name: "small fixed capacity",
			build: generatedInstance(vrp.GeneratorConfig{
				NumCustomers:  5,
				Vehicles:      2,
				Width:         30,
				Height:        30,
				Seed:          103,
				MinDemand:     1,
				MaxDemand:     3,
				CapacityMode:  vrp.CapacityFixed,
				FixedCapacity: 5,
			}),
		},
	}
}

func complexInstanceCases() []instanceCase {
	return []instanceCase{
		{
			name: "complex mixed auto capacity",
			build: generatedInstance(vrp.GeneratorConfig{
				NumCustomers:  6,
				Vehicles:      2,
				Width:         60,
				Height:        40,
				Seed:          201,
				MinDemand:     1,
				MaxDemand:     5,
				CapacityMode:  vrp.CapacityAuto,
				CapacitySlack: 1.10,
			}),
		},
		{
			name: "complex tight three vehicles",
			build: generatedInstance(vrp.GeneratorConfig{
				NumCustomers: 7,
				Vehicles:     3,
				Width:        70,
				Height:       50,
				Seed:         202,
				MinDemand:    1,
				MaxDemand:    4,
				CapacityMode: vrp.CapacityTight,
			}),
		},
	}
}

func generatedInstance(cfg vrp.GeneratorConfig) func() vrp.VRPInstance {
	return func() vrp.VRPInstance {
		_, instance := vrp.GenerateInstance(cfg)
		return instance
	}
}

func emptyInstance() vrp.VRPInstance {
	return vrp.VRPInstance{}
}

func solveBruteForceOracle(t *testing.T, instance vrp.VRPInstance) vrp.Solution {
	t.Helper()

	oracle := SolveBruteForce(instance, logging.NewLogger(false))
	assertValidSolution(t, instance, oracle)

	computedCost := solutionCost(instance, oracle)
	if diff := math.Abs(oracle.Cost - computedCost); diff > costEpsilon {
		t.Fatalf(
			"oracle cost mismatch: got %.12f, computed %.12f",
			oracle.Cost,
			computedCost,
		)
	}

	return oracle
}

func runAllSolvers(instance vrp.VRPInstance) []solverResult {
	solvers := allSolverCases()
	results := make([]solverResult, 0, len(solvers))

	for _, s := range solvers {
		results = append(results, solverResult{
			name:     s.name,
			solution: s.solve(instance),
		})
	}

	return results
}

func assertAllSolutionsValid(t *testing.T, instance vrp.VRPInstance, results []solverResult) {
	t.Helper()

	for _, result := range results {
		t.Run(result.name, func(t *testing.T) {
			assertValidSolution(t, instance, result.solution)
		})
	}
}

func assertNoSolverBeatsOracle(t *testing.T, oracleCost float64, results []solverResult) {
	t.Helper()

	if len(results) == 0 {
		t.Fatal("no solver results")
	}

	for _, result := range results {
		if result.solution.Cost+costEpsilon < oracleCost {
			t.Fatalf(
				"solver %q produced cost below oracle optimum: got %.12f, oracle %.12f",
				result.name,
				result.solution.Cost,
				oracleCost,
			)
		}
	}
}

func assertValidSolution(t *testing.T, instance vrp.VRPInstance, solution vrp.Solution) {
	t.Helper()

	if math.IsNaN(solution.Cost) || math.IsInf(solution.Cost, 0) {
		t.Fatalf("solution cost must be finite, got %.12f", solution.Cost)
	}

	if solution.Cost < -costEpsilon {
		t.Fatalf("solution cost must be non-negative, got %.12f", solution.Cost)
	}

	if len(solution.Routes) > instance.Vehicles {
		t.Fatalf("too many routes: got %d, vehicles %d", len(solution.Routes), instance.Vehicles)
	}

	computedCost := solutionCost(instance, solution)
	if diff := math.Abs(solution.Cost - computedCost); diff > costEpsilon {
		t.Fatalf(
			"solution cost mismatch: got %.12f, computed %.12f",
			solution.Cost,
			computedCost,
		)
	}

	if !isFeasibleCoverageAndCapacity(instance, solution) {
		t.Fatalf("solution is infeasible: %+v", solution)
	}
}

func isFeasibleCoverageAndCapacity(instance vrp.VRPInstance, solution vrp.Solution) bool {
	if len(instance.Customers) == 0 {
		return len(solution.Routes) == 0
	}

	servedByCustomerID := make(map[int]bool, len(instance.Customers))
	demandByCustomerID := make(map[int]int, len(instance.Customers))

	for _, customer := range instance.Customers {
		if customer.ID <= 0 {
			return false
		}
		if _, exists := servedByCustomerID[customer.ID]; exists {
			return false
		}

		servedByCustomerID[customer.ID] = false
		demandByCustomerID[customer.ID] = customer.Demand
	}

	for _, route := range solution.Routes {
		if len(route.Nodes) == 0 {
			continue
		}

		load := 0
		for _, customerID := range route.Nodes {
			alreadyServed, exists := servedByCustomerID[customerID]
			if !exists || alreadyServed {
				return false
			}

			servedByCustomerID[customerID] = true
			load += demandByCustomerID[customerID]
		}

		if load > instance.VehicleCapacity {
			return false
		}
	}

	for _, served := range servedByCustomerID {
		if !served {
			return false
		}
	}

	return true
}

func solutionCost(instance vrp.VRPInstance, solution vrp.Solution) float64 {
	total := 0.0

	for _, route := range solution.Routes {
		total += computeRouteCost(route.Nodes, instance.Dist)
	}

	return total
}

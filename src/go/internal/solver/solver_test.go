package solver

import (
	"math"
	"testing"

	"parallel-aco/internal/logging"
	"parallel-aco/internal/vrp"
)

func TestSolversOnSimpleInstance(t *testing.T) {
	instance := simpleInstance()
	expectedCost := 3.0

	tests := []struct {
		name   string
		solver func(vrp.VRPInstance) vrp.Solution
	}{
		{
			name: "brute force",
			solver: func(instance vrp.VRPInstance) vrp.Solution {
				return SolveBruteForce(instance, logging.NewLogger(false))
			},
		},
		{
			name: "greedy",
			solver: func(instance vrp.VRPInstance) vrp.Solution {
				return SolveGreedy(instance, nil)
			},
		},
		{
			name: "aco",
			solver: func(instance vrp.VRPInstance) vrp.Solution {
				return SolveACO(instance, nil, ACOConfig{
					NumAnts:    20,
					Iterations: 20,
					Seed:       42,
				})
			},
		},
		{
			name: "paco",
			solver: func(instance vrp.VRPInstance) vrp.Solution {
				return SolvePACO(instance, PACOConfig{
					BaseConfig: ACOConfig{
						NumAnts:    20,
						Iterations: 20,
						Seed:       42,
					},
					NumWorkers: 2,
				})
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			solution := tc.solver(instance)
			assertFeasibleSolution(t, instance, solution)

			if diff := math.Abs(solution.Cost - expectedCost); diff > 1e-9 {
				t.Fatalf("unexpected cost: got %.12f, want %.12f", solution.Cost, expectedCost)
			}

			if diff := math.Abs(solutionCost(instance, solution) - solution.Cost); diff > 1e-9 {
				t.Fatalf("solution cost does not match route cost: got %.12f, computed %.12f", solution.Cost, solutionCost(instance, solution))
			}
		})
	}
}

func TestSolversOnEmptyInstance(t *testing.T) {
	instance := vrp.VRPInstance{}

	tests := []struct {
		name   string
		solver func(vrp.VRPInstance) vrp.Solution
	}{
		{
			name: "brute force",
			solver: func(instance vrp.VRPInstance) vrp.Solution {
				return SolveBruteForce(instance, logging.NewLogger(false))
			},
		},
		{
			name: "greedy",
			solver: func(instance vrp.VRPInstance) vrp.Solution {
				return SolveGreedy(instance, nil)
			},
		},
		{
			name: "aco",
			solver: func(instance vrp.VRPInstance) vrp.Solution {
				return SolveACO(instance, nil, DefaultACOConfig())
			},
		},
		{
			name: "paco",
			solver: func(instance vrp.VRPInstance) vrp.Solution {
				return SolvePACO(instance, PACOConfig{BaseConfig: DefaultACOConfig()})
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			solution := tc.solver(instance)
			assertFeasibleSolution(t, instance, solution)

			if solution.Cost != 0 {
				t.Fatalf("unexpected cost: got %.12f, want 0", solution.Cost)
			}
		})
	}
}

func TestSolversOnHardInstance(t *testing.T) {
	instance := hardInstance()

	results := []struct {
		name     string
		solution vrp.Solution
	}{
		{
			name:     "brute force",
			solution: SolveBruteForce(instance, logging.NewLogger(false)),
		},
		{
			name:     "greedy",
			solution: SolveGreedy(instance, nil),
		},
		{
			name: "aco",
			solution: SolveACO(instance, nil, ACOConfig{
				NumAnts:    20,
				Iterations: 20,
				Seed:       42,
			}),
		},
		{
			name: "paco",
			solution: SolvePACO(instance, PACOConfig{
				BaseConfig: ACOConfig{
					NumAnts:    20,
					Iterations: 20,
					Seed:       42,
				},
				NumWorkers: 2,
			}),
		},
	}

	for _, result := range results {
		t.Run(result.name, func(t *testing.T) {
			assertFeasibleSolution(t, instance, result.solution)

			if diff := math.Abs(solutionCost(instance, result.solution) - result.solution.Cost); diff > 1e-9 {
				t.Fatalf("solution cost does not match route cost: got %.12f, computed %.12f", result.solution.Cost, solutionCost(instance, result.solution))
			}
		})
	}

	bruteForceCost := results[0].solution.Cost
	greedyCost := results[1].solution.Cost
	acoCost := results[2].solution.Cost
	pacoCost := results[3].solution.Cost

	if bruteForceCost > greedyCost {
		t.Fatalf("expected brute force to be no worse than greedy: brute force %.12f, greedy %.12f", bruteForceCost, greedyCost)
	}

	if acoCost < bruteForceCost || acoCost > greedyCost {
		t.Fatalf("expected ACO cost to be between brute force and greedy: brute force %.12f, aco %.12f, greedy %.12f", bruteForceCost, acoCost, greedyCost)
	}

	if pacoCost < bruteForceCost || pacoCost > greedyCost {
		t.Fatalf("expected PACO cost to be between brute force and greedy: brute force %.12f, paco %.12f, greedy %.12f", bruteForceCost, pacoCost, greedyCost)
	}
}

func simpleInstance() vrp.VRPInstance {
	return vrp.VRPInstance{
		Vehicles:        1,
		VehicleCapacity: 10,
		Customers: []vrp.Point{
			{ID: 1, X: 1, Y: 0, Demand: 1},
			{ID: 2, X: 0, Y: 1, Demand: 1},
		},
		Dist: [][]float64{
			{0, 1, 1},
			{1, 0, 1},
			{1, 1, 0},
		},
	}
}

func hardInstance() vrp.VRPInstance {
	_, instance := vrp.GenerateInstance(vrp.GeneratorConfig{
		NumCustomers: 8,
		Vehicles:     3,
		Width:        100,
		Height:       100,
		Seed:         42,
	})
	return instance
}

func assertFeasibleSolution(t *testing.T, instance vrp.VRPInstance, solution vrp.Solution) {
	t.Helper()

	if !isSolutionFeasible(instance, solution) {
		t.Fatalf("solution is not feasible: %+v", solution)
	}
}

func isSolutionFeasible(instance vrp.VRPInstance, solution vrp.Solution) bool {
	if len(instance.Customers) == 0 {
		return len(solution.Routes) == 0
	}

	seen := make(map[int]bool, len(instance.Customers))
	demandByID := make(map[int]int, len(instance.Customers))
	totalDemand := 0

	for _, customer := range instance.Customers {
		if customer.ID < 0 {
			return false
		}
		if _, exists := seen[customer.ID]; exists {
			return false
		}

		seen[customer.ID] = false
		demandByID[customer.ID] = customer.Demand
		totalDemand += customer.Demand
	}

	servedDemand := 0
	for _, route := range solution.Routes {
		load := 0
		for _, node := range route.Nodes {
			if _, exists := seen[node]; !exists || seen[node] {
				return false
			}

			seen[node] = true
			load += demandByID[node]
			servedDemand += demandByID[node]
		}

		if load > instance.VehicleCapacity {
			return false
		}
	}

	if servedDemand != totalDemand {
		return false
	}

	for _, wasSeen := range seen {
		if !wasSeen {
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

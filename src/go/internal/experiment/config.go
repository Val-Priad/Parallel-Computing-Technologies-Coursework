package experiment

import (
	"parallel-aco/internal/vrp"
)

type ExperimentConfig struct {
	NumCustomers int
	Vehicles     int
	Width        float64
	Height       float64

	CapacityMode vrp.CapacityMode
	Seed         int64
}

func generateExperiments(
	count int,
	startCustomers, customerStep int,
	vehiclesFn func(i int) int,
	sizeFn func(i int) float64,
	capacityMode vrp.CapacityMode,
	seedStart int64,
) []ExperimentConfig {
	experiments := make([]ExperimentConfig, count)
	for i := 0; i < count; i++ {
		size := sizeFn(i)
		experiments[i] = ExperimentConfig{
			NumCustomers: startCustomers + i*customerStep,
			Vehicles:     vehiclesFn(i),
			Width:        size,
			Height:       size,
			CapacityMode: capacityMode,
			Seed:         seedStart + int64(i),
		}
	}

	return experiments
}

func GetExperiments() []ExperimentConfig {
	return generateExperiments(
		3,
		4,
		2,
		func(i int) int { return 2 + (i+1)/2 },
		func(i int) float64 {
			if i == 0 {
				return 50
			}
			return float64(40 + i*20)
		},
		vrp.CapacityTight,
		67,
	)
}

func GetLargeComparisonExperiments(qty int) []ExperimentConfig {
	return generateExperiments(
		qty,
		50,
		10,
		func(i int) int { return 8 + i },
		func(i int) float64 { return float64(180 + i*20) },
		vrp.CapacityAuto,
		25,
	)
}

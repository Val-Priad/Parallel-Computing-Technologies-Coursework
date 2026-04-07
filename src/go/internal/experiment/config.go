package experiment

import "parallel-aco/internal/vrp"

type ExperimentConfig struct {
	Name         string
	NumCustomers int
	Vehicles     int
	Width        float64
	Height       float64

	CapacityMode vrp.CapacityMode
	Seed         int64
}

func GetExperiments() []ExperimentConfig {
	return []ExperimentConfig{
		{
			Name:         "exp_1_easy",
			NumCustomers: 6,
			Vehicles:     2,
			Width:        50,
			Height:       50,
			CapacityMode: vrp.CapacityTight,
			Seed:         42,
		},
		{
			Name:         "exp_2_medium",
			NumCustomers: 8,
			Vehicles:     3,
			Width:        60,
			Height:       60,
			CapacityMode: vrp.CapacityTight,
			Seed:         43,
		},
		{
			Name:         "exp_3_hard",
			NumCustomers: 10,
			Vehicles:     3,
			Width:        80,
			Height:       80,
			CapacityMode: vrp.CapacityTight,
			Seed:         44,
		},
		{
			Name:         "exp_4_very_hard",
			NumCustomers: 12,
			Vehicles:     4,
			Width:        100,
			Height:       100,
			CapacityMode: vrp.CapacityTight,
			Seed:         45,
		},
		{
			Name:         "exp_5_extreme",
			NumCustomers: 14,
			Vehicles:     4,
			Width:        120,
			Height:       120,
			CapacityMode: vrp.CapacityTight,
			Seed:         46,
		},
	}
}

func GetLargeComparisonExperiments() []ExperimentConfig {
	return []ExperimentConfig{
		{
			Name:         "large_1_60",
			NumCustomers: 60,
			Vehicles:     8,
			Width:        180,
			Height:       180,
			CapacityMode: vrp.CapacityAuto,
			Seed:         142,
		},
		{
			Name:         "large_2_80",
			NumCustomers: 80,
			Vehicles:     10,
			Width:        220,
			Height:       220,
			CapacityMode: vrp.CapacityAuto,
			Seed:         143,
		},
		{
			Name:         "large_3_100",
			NumCustomers: 100,
			Vehicles:     12,
			Width:        260,
			Height:       260,
			CapacityMode: vrp.CapacityAuto,
			Seed:         144,
		},
		{
			Name:         "large_4_120",
			NumCustomers: 120,
			Vehicles:     14,
			Width:        300,
			Height:       300,
			CapacityMode: vrp.CapacityAuto,
			Seed:         145,
		},
	}
}

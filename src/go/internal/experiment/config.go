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

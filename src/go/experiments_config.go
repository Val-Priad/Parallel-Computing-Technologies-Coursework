package main

type ExperimentConfig struct {
	Name         string
	NumCustomers int
	Vehicles     int
	Width        float64
	Height       float64

	CapacityMode CapacityMode
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
			CapacityMode: CapacityAuto,
			Seed:         42,
		},

		{
			Name:         "exp_2_light",
			NumCustomers: 8,
			Vehicles:     3,
			Width:        60,
			Height:       60,
			CapacityMode: CapacityAuto,
			Seed:         43,
		},

		{
			Name:         "exp_3_medium",
			NumCustomers: 10,
			Vehicles:     3,
			Width:        80,
			Height:       80,
			CapacityMode: CapacityAuto,
			Seed:         44,
		},

		{
			Name:         "exp_4_hard",
			NumCustomers: 12,
			Vehicles:     4,
			Width:        100,
			Height:       100,
			CapacityMode: CapacityAuto,
			Seed:         45,
		},

		{
			Name:         "exp_5_extreme",
			NumCustomers: 14,
			Vehicles:     4,
			Width:        120,
			Height:       120,
			CapacityMode: CapacityAuto,
			Seed:         46,
		},
	}
}

package main

import (
	"math"
	"math/rand"
	"time"
)

type CapacityMode string

const (
	CapacityAuto  CapacityMode = "auto"
	CapacityTight CapacityMode = "tight"
	CapacityLoose CapacityMode = "loose"
	CapacityFixed CapacityMode = "fixed"
)

type GeneratorConfig struct {
	NumCustomers int
	Vehicles     int
	Width        float64
	Height       float64
	Seed         int64

	MinDemand int
	MaxDemand int

	CapacityMode  CapacityMode
	CapacitySlack float64
	FixedCapacity int
}

func GenerateInstance(cfg GeneratorConfig) ([]Point, VRPInstance) {
	if cfg.Seed == 0 {
		cfg.Seed = time.Now().UnixNano()
	}

	if cfg.MinDemand <= 0 {
		cfg.MinDemand = 1
	}
	if cfg.MaxDemand < cfg.MinDemand {
		cfg.MaxDemand = cfg.MinDemand
	}
	if cfg.Width <= 0 {
		cfg.Width = 100
	}
	if cfg.Height <= 0 {
		cfg.Height = 100
	}
	if cfg.CapacitySlack <= 0 {
		cfg.CapacitySlack = 1.15
	}
	if cfg.CapacityMode == "" {
		cfg.CapacityMode = CapacityAuto
	}

	rng := rand.New(rand.NewSource(cfg.Seed))

	points := make([]Point, 0, cfg.NumCustomers+1)

	depot := Point{
		ID:     0,
		X:      cfg.Width / 2,
		Y:      cfg.Height / 2,
		Demand: 0,
	}
	points = append(points, depot)

	totalDemand := 0
	maxDemand := 0

	for i := 1; i <= cfg.NumCustomers; i++ {
		demand := rng.Intn(cfg.MaxDemand-cfg.MinDemand+1) + cfg.MinDemand

		totalDemand += demand
		if demand > maxDemand {
			maxDemand = demand
		}

		points = append(points, Point{
			ID:     i,
			X:      rng.Float64() * cfg.Width,
			Y:      rng.Float64() * cfg.Height,
			Demand: demand,
		})
	}

	vehicleCapacity := computeCapacity(cfg, totalDemand, maxDemand)

	dist := BuildDistanceMatrix(points)

	instance := VRPInstance{
		Depot:           points[0],
		Customers:       points[1:],
		Vehicles:        cfg.Vehicles,
		VehicleCapacity: vehicleCapacity,
		Dist:            dist,
	}

	return points, instance
}

func computeCapacity(cfg GeneratorConfig, totalDemand, maxDemand int) int {
	if totalDemand <= 0 {
		return 1
	}

	if cfg.Vehicles <= 0 {
		return max(maxDemand, 1)
	}

	switch cfg.CapacityMode {
	case CapacityTight:
		cap := totalDemand / cfg.Vehicles
		if cap < maxDemand {
			cap = maxDemand
		}
		return max(cap, 1)

	case CapacityLoose:
		return max(totalDemand, maxDemand)
	case CapacityFixed:
		if cfg.FixedCapacity <= 0 {
			return max(maxDemand, 1)
		}
		return max(cfg.FixedCapacity, maxDemand)
	case CapacityAuto:
		fallthrough
	default:
		cap := int(math.Ceil(float64(totalDemand) * cfg.CapacitySlack / float64(cfg.Vehicles)))
		if cap < maxDemand {
			cap = maxDemand
		}
		return max(cap, 1)
	}
}

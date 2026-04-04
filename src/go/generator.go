package main

import (
	"math/rand"
	"time"
)

type GeneratorConfig struct {
	NumCustomers int
	Vehicles     int
	Width        float64
	Height       float64
	Seed         int64
}

func GenerateInstance(cfg GeneratorConfig) ([]Point, VRPInstance) {
	if cfg.Seed == 0 {
		cfg.Seed = time.Now().UnixNano()
	}
	rng := rand.New(rand.NewSource(cfg.Seed))

	points := make([]Point, 0, cfg.NumCustomers+1)
	totalDemand := 0
	maxDemand := 0

	depot := Point{
		ID:     0,
		X:      cfg.Width / 2,
		Y:      cfg.Height / 2,
		Demand: 0,
	}
	points = append(points, depot)

	for i := 1; i <= cfg.NumCustomers; i++ {
		demand := rng.Intn(5) + 1
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

	vehicleCapacity := totalDemand
	if cfg.Vehicles > 1 && totalDemand > 0 {
		vehicleCapacity = (totalDemand * 125) / (cfg.Vehicles * 100)
		if vehicleCapacity < maxDemand {
			vehicleCapacity = maxDemand
		}
	}
	if vehicleCapacity < 1 {
		vehicleCapacity = 1
	}

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

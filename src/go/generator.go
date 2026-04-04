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

	depot := Point{
		ID: 0,
		X:  cfg.Width / 2,
		Y:  cfg.Height / 2,
	}
	points = append(points, depot)

	for i := 1; i <= cfg.NumCustomers; i++ {
		points = append(points, Point{
			ID: i,
			X:  rng.Float64() * cfg.Width,
			Y:  rng.Float64() * cfg.Height,
		})
	}

	dist := BuildDistanceMatrix(points)

	instance := VRPInstance{
		Depot:     points[0],
		Customers: points[1:],
		Vehicles:  cfg.Vehicles,
		Dist:      dist,
	}

	return points, instance
}

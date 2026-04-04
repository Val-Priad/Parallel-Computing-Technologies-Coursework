package main

type Point struct {
	ID     int
	X      float64
	Y      float64
	Demand int
}

type VRPInstance struct {
	Depot           Point
	Customers       []Point
	Vehicles        int
	VehicleCapacity int
	Dist            [][]float64
}

type SearchMetrics struct {
	DurationMS float64 `json:"duration_ms"`
}

type Route struct {
	VehicleID int
	Nodes     []int // IDs клиентов
}

type Solution struct {
	Routes  []Route
	Cost    float64
	Metrics SearchMetrics `json:"metrics"`
}

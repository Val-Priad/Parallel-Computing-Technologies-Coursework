package main

type Point struct {
	ID int
	X  float64
	Y  float64
}

type VRPInstance struct {
	Depot     Point
	Customers []Point
	Vehicles  int
	Dist      [][]float64
}

type Route struct {
	VehicleID int
	Nodes     []int // IDs клиентов
}

type Solution struct {
	Routes []Route
	Cost   float64
}

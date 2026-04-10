package vrp

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
	CapacityMode    CapacityMode
	Dist            [][]float64
}

type Route struct {
	VehicleID int
	Nodes     []int
}

type Solution struct {
	Routes     []Route
	Cost       float64
	DurationMS float64 `json:"duration_ms"`
}

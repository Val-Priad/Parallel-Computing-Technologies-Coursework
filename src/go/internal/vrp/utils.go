package vrp

import "math"

func Distance(a, b Point) float64 {
	return math.Sqrt((a.X-b.X)*(a.X-b.X) + (a.Y-b.Y)*(a.Y-b.Y))
}

func BuildDistanceMatrix(points []Point) [][]float64 {
	n := len(points)
	dist := make([][]float64, n)
	for i := 0; i < n; i++ {
		dist[i] = make([]float64, n)
		for j := 0; j < n; j++ {
			dist[i][j] = Distance(points[i], points[j])
		}
	}
	return dist
}

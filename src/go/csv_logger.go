package main

import (
	"encoding/csv"
	"os"
	"strconv"
)

type CSVLogger struct {
	file   *os.File
	writer *csv.Writer
}

func NewCSVLogger(filename string) (*CSVLogger, error) {
	file, err := os.Create(filename)
	if err != nil {
		return nil, err
	}

	writer := csv.NewWriter(file)

	writer.Write([]string{
		"experiment",
		"algorithm",
		"customers",
		"vehicles",
		"capacity",
		"time_ms",
		"cost",
		"feasible",
	})

	return &CSVLogger{
		file:   file,
		writer: writer,
	}, nil
}

func (c *CSVLogger) Log(
	expName string,
	algorithm string,
	instance VRPInstance,
	solution Solution,
) {
	feasible := isSolutionComplete(instance, solution)

	c.writer.Write([]string{
		expName,
		algorithm,
		strconv.Itoa(len(instance.Customers)),
		strconv.Itoa(instance.Vehicles),
		strconv.Itoa(instance.VehicleCapacity),
		strconv.FormatFloat(solution.Metrics.DurationMS, 'f', 3, 64),
		strconv.FormatFloat(solution.Cost, 'f', 3, 64),
		strconv.FormatBool(feasible),
	})
}

func isSolutionComplete(instance VRPInstance, solution Solution) bool {
	if len(instance.Customers) == 0 {
		return true
	}

	if len(solution.Routes) == 0 {
		return false
	}

	maxID := 0
	totalDemand := 0
	for _, customer := range instance.Customers {
		if customer.ID > maxID {
			maxID = customer.ID
		}
		totalDemand += customer.Demand
	}

	seen := make([]bool, maxID+1)
	demandByID := make([]int, maxID+1)
	isCustomer := make([]bool, maxID+1)

	for _, customer := range instance.Customers {
		if customer.ID < 0 || customer.ID >= len(seen) {
			return false
		}
		demandByID[customer.ID] = customer.Demand
		isCustomer[customer.ID] = true
	}

	servedDemand := 0
	for _, route := range solution.Routes {
		routeLoad := 0
		for _, node := range route.Nodes {
			if node < 0 || node >= len(seen) {
				return false
			}
			if !isCustomer[node] || seen[node] {
				return false
			}

			seen[node] = true
			routeLoad += demandByID[node]
			servedDemand += demandByID[node]
		}
		if routeLoad > instance.VehicleCapacity {
			return false
		}
	}

	if servedDemand != totalDemand {
		return false
	}

	for _, customer := range instance.Customers {
		if !seen[customer.ID] {
			return false
		}
	}

	return true
}

func (c *CSVLogger) Close() {
	c.writer.Flush()
	c.file.Close()
}

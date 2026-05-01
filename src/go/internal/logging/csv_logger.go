package logging

import (
	"encoding/csv"
	"os"
	"parallel-aco/internal/vrp"
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
		"capacity_mode",
		"time_ms",
		"cost",
	})

	return &CSVLogger{
		file:   file,
		writer: writer,
	}, nil
}

func (c *CSVLogger) Log(
	expName string,
	algorithm string,
	instance vrp.VRPInstance,
	solution vrp.Solution,
) {
	c.writer.Write([]string{
		expName,
		algorithm,
		strconv.Itoa(len(instance.Customers)),
		strconv.Itoa(instance.Vehicles),
		strconv.Itoa(instance.VehicleCapacity),
		string(instance.CapacityMode),
		strconv.FormatFloat(solution.DurationMS, 'f', 3, 64),
		strconv.FormatFloat(solution.Cost, 'f', 3, 64),
	})
}

func (c *CSVLogger) Close() {
	c.writer.Flush()
	c.file.Close()
}

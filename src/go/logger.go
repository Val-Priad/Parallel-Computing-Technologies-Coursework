package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Step struct {
	StepID int     `json:"step_id"`
	Routes [][]int `json:"routes"`
	Cost   float64 `json:"cost"`
}

type LogOutput struct {
	Points map[string][]float64 `json:"points"`
	Steps  []Step               `json:"steps"`
}

type Logger struct {
	Enabled bool
	Steps   []Step
}

func NewLogger(enabled bool) *Logger {
	return &Logger{
		Enabled: enabled,
	}
}

func (l *Logger) Log(step Step) {
	if !l.Enabled {
		return
	}
	l.Steps = append(l.Steps, step)
}

func (l *Logger) SaveToFile(filename string, points []Point) error {
	if !l.Enabled {
		return nil
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	pointsMap := make(map[string][]float64, len(points))
	for _, p := range points {
		pointsMap[fmt.Sprintf("%d", p.ID)] = []float64{p.X, p.Y}
	}

	output := LogOutput{
		Points: pointsMap,
		Steps:  l.Steps,
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(output)
}

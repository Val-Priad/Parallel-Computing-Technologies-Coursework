package logging

import (
	"encoding/json"
	"fmt"
	"os"
	"parallel-aco/internal/vrp"
	"path/filepath"
)

type Step struct {
	StepID int     `json:"step_id"`
	Routes [][]int `json:"routes"`
	Cost   float64 `json:"cost"`
}

type LogOutput struct {
	Points     map[string][]float64 `json:"points"`
	Steps      []Step               `json:"steps"`
	DurationMS float64              `json:"duration_ms"`
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

func (l *Logger) SaveToFile(filename string, points []vrp.Point, durationMS float64) error {
	if !l.Enabled {
		return nil
	}

	if dir := filepath.Dir(filename); dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
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
		Points:     pointsMap,
		Steps:      l.Steps,
		DurationMS: durationMS,
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(output)
}

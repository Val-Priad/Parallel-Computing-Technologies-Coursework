package main

import (
	"encoding/json"
	"os"
)

type Step struct {
	StepID int     `json:"step_id"`
	Routes [][]int `json:"routes"`
	Cost   float64 `json:"cost"`
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

func (l *Logger) SaveToFile(filename string) error {
	if !l.Enabled {
		return nil
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(l.Steps)
}

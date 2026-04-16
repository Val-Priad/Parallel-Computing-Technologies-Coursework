package experiment

import (
	"fmt"
	"parallel-aco/internal/logging"
	"parallel-aco/internal/vrp"
	"path/filepath"
)

const resultsDir = "results"

func writeExperimentResult(
	csvLogger *logging.CSVLogger,
	runName string,
	algo string,
	logSubdir string,
	instance vrp.VRPInstance,
	solution vrp.Solution,
	logger *logging.Logger,
	points []vrp.Point,
) {
	fmt.Printf("%s: cost=%.3f, time=%.3fms\n", algo, solution.Cost, solution.DurationMS)

	if logger != nil {
		logDir := resultsDir
		if logSubdir != "" {
			logDir = filepath.Join(resultsDir, logSubdir)
		}

		logPath := filepath.Join(logDir, fmt.Sprintf("%s_%s.json", runName, algo))
		if err := logger.SaveToFile(logPath, points, solution.DurationMS); err != nil {
			fmt.Printf("Failed to save %s log: %v\n", algo, err)
		}
	}

	csvLogger.Log(runName, algo, instance, solution)
}

package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"wintercloak/internal/engines"
	"wintercloak/internal/utils/engineutils"
)

func (a *App) ListAvailableEngines() []engines.EngineMeta {
	return engines.Available.List()
}

type GeneratePatchedFileArgs struct {
	Engine   string             `json:"engine"`
	Content  string             `json:"content"`
	Settings engines.EngineArgs `json:"settings"`
}

func (a *App) GeneratePatchedFile(frontArgs GeneratePatchedFileArgs) (string, error) {
	var patched []byte

	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "wclk-patch-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer func() {
		_ = tmpFile.Close()
		_ = os.Remove(tmpFile.Name())
	}()

	destFilepath := tmpFile.Name()

	// Write the content to temp file
	if err = os.WriteFile(destFilepath, []byte(frontArgs.Content), 0644); err != nil {
		return "", fmt.Errorf("failed to write content to temp file: %w", err)
	}

	//Initialize the appropriate engine based on engine name
	eng, err := engineutils.NewEngine(strings.ToLower(frontArgs.Engine), frontArgs.Settings)
	if err != nil {
		return "", fmt.Errorf("failed to create engine: %w", err)
	}

	// Apply the patch
	// Skip go build tag handling (not configurable from front)
	if patched, err = eng.Run(engines.RunOptions{
		TargetPath:       destFilepath,
		IgnoreExtensions: true,
	}); err != nil {
		return "", fmt.Errorf("failed to apply patch: %w", err)
	}

	return string(patched), nil
}

func (a *App) RunPreprocessor(runArgs engines.RunArgs) (string, error) {
	count := 0
	var out []byte

	err := filepath.Walk(runArgs.RunOptions.TargetDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip if it's a directory
		if info.IsDir() {
			return nil
		}

		runArgs.RunOptions.TargetPath = path

		eng, err := engineutils.NewEngine(strings.ToLower(runArgs.Engine), runArgs.Args)
		if err != nil {
			return err
		}

		out, err = eng.Run(runArgs.RunOptions)
		if err != nil {
			//	skipping
			log.Warnf("failed to process file %s: %s", path, err.Error())
			return nil
		}

		err = eng.WritePatchedFile(out)
		if err != nil {
			return err
		}

		count++
		return nil
	})

	if err != nil {
		return "", fmt.Errorf("failed to run preprocessor: %w", err)
	}

	return fmt.Sprintf("Processed %d files", count), nil
}

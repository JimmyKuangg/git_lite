package core

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func EnsureGitliteRepo() (string, error) {
	current, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("cannot get working directory: %w", err)
	}

	for {
		gitDir := filepath.Join(current, ".gitlite")
		info, err := os.Stat(gitDir)

		if err == nil && info.IsDir() {
			return current, nil
		}

		if err != nil && !os.IsNotExist(err) {
			return "", fmt.Errorf("error checking directory: %w", err)
		}

		parent := filepath.Dir(current)

		if parent == current {
			return "", errors.New("error: not a gitlite repository (or any of the parent directories)")
		}

		current = parent
	}
}

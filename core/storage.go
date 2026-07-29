package core

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const GitDir = ".gitlite"

func WriteObject(content []byte) (string, error) {
	hash := Hash(content)

	if len(hash) < 40 {
		return "", errors.New("invalid hash")
	}

	dir := filepath.Join(GitDir, "objects", hash[:2])
	file := filepath.Join(dir, hash[2:])

	err := os.MkdirAll(dir, 0o755)
	if err != nil {
		return "", fmt.Errorf("error creating directory %s: %w", dir, err)
	}

	err = os.WriteFile(file, content, 0o644)
	if err != nil {
		return "", fmt.Errorf("error writing file %s: %w", file, err)
	}

	return hash, nil
}

func ReadObject(hash string) ([]byte, error) {
	if len(hash) < 40 {
		return nil, errors.New("invalid hash")
	}

	file := filepath.Join(GitDir, "objects", hash[:2], hash[2:])
	bytes, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("error reading object %s: %w", hash, err)
	}

	return bytes, nil
}

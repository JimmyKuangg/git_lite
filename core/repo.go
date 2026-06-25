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
		info, err := os.Stat(gitDir);

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

func ScanWorkingDirectory() (map[string]string, error) {
	root, err := EnsureGitliteRepo()
	if err != nil {
		return nil, err
	}

	workingIndex, err := BuildWorkingSnapshot(root)
	if err != nil {
		return nil, err
	}

	stagedIndex, err := ReadIndex()
	if err != nil {
		return nil, err
	}

	for path, workingHash := range workingIndex {
		stagedHash, exists := stagedIndex.Entries[path]
		if !exists {
			fmt.Printf("untracked: %s\n", path)
			continue
		}

		if stagedHash != workingHash {
			fmt.Printf("modified: %s\n", path)
		}
	}

	for path := range stagedIndex.Entries {
    _, exists := workingIndex[path]
    if !exists {
        fmt.Printf("deleted: %s\n", path)
    }
}

	return nil, nil
}

func BuildWorkingSnapshot(root string) (map[string]string, error) {
	index := make(map[string]string)

	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, fmt.Errorf("error reading directory, %w", err)
	}

	for _, entry := range entries {
		fullPath := filepath.Join(root, entry.Name())
		relativePath, err := filepath.Rel(root, fullPath)
		if err != nil {
				return nil, err
		}

		if !entry.IsDir() {
			content, err := os.ReadFile(fullPath)
			if err != nil {
				return nil, err
			}
			hash := Hash(content)

			index[relativePath] = hash
		} else if entry.Name() != ".gitlite" && entry.Name() != ".git" {
			subIndex, err := BuildWorkingSnapshot(fullPath)
			if err != nil {
				return nil, fmt.Errorf("error index from root, %w", err)
			}

			for key, val := range subIndex {
				index[relativePath + "/" + key] = val
			}
		}
	}

	return index, nil
}
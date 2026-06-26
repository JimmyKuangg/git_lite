package core

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type Status struct {
	Untracked []string
	Modified 	[]string
	Deleted 	[]string
}

type Staged struct {
	Untracked []string
	Modified 	[]string
	Deleted 	[]string
}

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

func ScanWorkingDirectory() (Status, error) {
	var statuses Status

	root, err := EnsureGitliteRepo()
	if err != nil {
		return statuses, err
	}

	workingIndex, err := BuildWorkingSnapshot(root)
	if err != nil {
		return statuses, err
	}

	stagedIndex, err := ReadIndex()
	if err != nil {
		return statuses, err
	}

	for path, workingHash := range workingIndex {
		stagedHash, exists := stagedIndex.Entries[path]
		if !exists {
			statuses.Untracked = append(statuses.Untracked, path)
			continue
		}

		if stagedHash != workingHash {
			statuses.Modified = append(statuses.Modified, path)
		}
	}

	for path := range stagedIndex.Entries {
    _, exists := workingIndex[path]
    if !exists {
      statuses.Deleted = append(statuses.Deleted, path)
    }
	}

	return statuses, nil
}

func ScanStagedDirectory() (Staged, error) {
	var staged Staged

	headHash, err := ReadHEAD()
	if err != nil {
		return Staged{}, err
	}

	headBytes, err := ReadObject(headHash)
	if err != nil {
		return Staged{}, err
	}

	prevCommit, err := ParseCommit(headBytes)
	if err != nil {
		return Staged{}, err
	}

	prevTreeHash := prevCommit.Root
	prevCommitTree, err := FlattenTree(prevTreeHash)
	if err != nil {
		return Staged{}, err
	}

	stagedIndex, err := ReadIndex()
	if err != nil {
		return Staged{}, err
	}

	for path, hash := range stagedIndex.Entries {
		prevHash, exists := prevCommitTree[path]

		if !exists {
			staged.Untracked = append(staged.Untracked, path)
			continue
		}

		if prevHash != hash {
			staged.Modified = append(staged.Modified, path)
		}
	}

	for path := range prevCommitTree {
		if _, exists := stagedIndex.Entries[path]; !exists {
			staged.Deleted = append(staged.Deleted, path)
		}
	}

	return staged, nil
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
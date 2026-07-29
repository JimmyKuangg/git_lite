package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Diff struct {
	Untracked []string
	Modified  []string
	Deleted   []string
}

func ScanWorkingDirectory() (Diff, error) {
	root, err := EnsureGitliteRepo()
	if err != nil {
		return Diff{}, err
	}

	workingIndex, err := BuildWorkingSnapshot(root)
	if err != nil {
		return Diff{}, err
	}

	stagedIndex, err := ReadIndex()
	if err != nil {
		return Diff{}, err
	}

	return CompareSnapshots(workingIndex, stagedIndex.Entries), nil
}

func ScanStagedDirectory() (Diff, error) {
	headHash, err := ReadHEAD()
	if err != nil {
		return Diff{}, err
	}

	prevCommitTree := make(map[string]string)

	if headHash != "" && !strings.HasPrefix(headHash, "ref") {
		headBytes, err := ReadObject(headHash)
		if err != nil {
			return Diff{}, err
		}

		prevCommit, err := ParseCommit(headBytes)
		if err != nil {
			return Diff{}, err
		}

		prevTreeHash := prevCommit.Root
		prevCommitTree, err = FlattenTree(prevTreeHash)
		if err != nil {
			return Diff{}, err
		}
	}

	stagedIndex, err := ReadIndex()
	if err != nil {
		return Diff{}, err
	}

	return CompareSnapshots(stagedIndex.Entries, prevCommitTree), nil
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
				index[relativePath+"/"+key] = val
			}
		}
	}

	return index, nil
}

func CompareSnapshots(current, previous map[string]string) Diff {
	var diffs Diff

	for path, currentHash := range current {
		previousHash, exists := previous[path]

		if !exists {
			diffs.Untracked = append(diffs.Untracked, path)
			continue
		}

		if currentHash != previousHash {
			diffs.Modified = append(diffs.Modified, path)
		}
	}

	for path := range previous {
		if _, exists := current[path]; !exists {
			diffs.Deleted = append(diffs.Deleted, path)
		}
	}

	return diffs
}

func (d Diff) Empty() bool {
	return len(d.Deleted) == 0 && len(d.Modified) == 0 && len(d.Untracked) == 0
}

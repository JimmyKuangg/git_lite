package core

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type IndexEntry struct {
	Path string
	Hash string
}

type Index struct {
	Entries map[string]string
}

const IndexPath = ".gitlite/index"

func (i *Index) Add(path string, hash string) {
	i.Entries[path] = hash
}

func ReadIndex() (*Index, error) {
	content, err := os.ReadFile(IndexPath)
	
  if os.IsNotExist(err) {
    return &Index{
      Entries: make(map[string]string),
    }, nil
	} else if err != nil {
		return nil, fmt.Errorf("unable to read index: %w", err)
	}

	index := &Index{
    Entries: make(map[string]string),
	}
	lines := strings.Split(string(content), "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}

		split := strings.Split(line, " ")
		if len(split) != 2 {
			return nil, fmt.Errorf("invalid index entry: %s", line)
		}
		
		index.Entries[split[0]] = split[1]
	}

	return index, nil
}

func (i *Index) Save() error {
	var entriesSlice []IndexEntry

	for path, hash := range i.Entries {
		entriesSlice = append(entriesSlice, IndexEntry{
			Path: path, 
			Hash: hash,
		})
	}

	sort.Slice(entriesSlice, func(i, j int) bool {
		return entriesSlice[i].Path < entriesSlice[j].Path
	})

	var buf bytes.Buffer

	for _, entry := range entriesSlice {
		buf.WriteString(entry.Path)
		buf.WriteByte(' ')
		buf.WriteString(entry.Hash)
		buf.WriteByte('\n')
	}

	err := os.WriteFile(IndexPath, buf.Bytes(), 0644)
	if err != nil {
		return fmt.Errorf("error writing index file: %w", err)
	}

	return nil
}

func ScanWorkingDirectory() (map[string]string, error) {
	root, err := EnsureGitliteRepo()
	if err != nil {
		return nil, err
	}

	_, err = BuildIndexFromRoot(root)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func BuildIndexFromRoot(root string) (map[string]string, error) {
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
			subIndex, err := BuildIndexFromRoot(fullPath)
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
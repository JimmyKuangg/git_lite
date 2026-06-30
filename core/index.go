package core

import (
	"bytes"
	"fmt"
	"os"
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

func (i *Index) Remove(path string) {
	delete(i.Entries, path)
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

func WriteIndex(snapshot map[string]string) error {
	i := Index{
		Entries: make(map[string]string),
	}

	for path, hash := range snapshot {
		i.Entries[path] = hash
	}

	return i.Save()
}
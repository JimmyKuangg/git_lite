package core

import (
	"fmt"
	"os"
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
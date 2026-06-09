package core

import (
	"bytes"
	"errors"
	"sort"
)

type TreeEntry struct {
  Name string
  Type string
  Hash string
}

type Tree struct {
  Entries []TreeEntry
}

const (
  BlobType = "blob"
  TreeType = "tree"
)

func (t *Tree) Add(name string, typeof string, hash string) error {
	if (name == "") {
		return errors.New("invalid name")
	}

	if (typeof != BlobType && typeof != TreeType) {
		return errors.New("invalid type")
	}

	if (hash == "") {
		return errors.New("invalid hash")
	}

	t.Entries = append(t.Entries, TreeEntry{
		Name: name,
		Type: typeof,
		Hash: hash,
	})

	return nil
}

func (t *Tree) Encode() []byte {
	sort.Slice(t.Entries, func(i, j int) bool {
		return t.Entries[i].Name < t.Entries[j].Name
	})

	var buf bytes.Buffer

	for _, entry := range t.Entries {
		buf.WriteString(entry.Type)
		buf.WriteByte(' ')
		buf.WriteString(entry.Name)
		buf.WriteByte(' ')
		buf.WriteString(entry.Hash)
		buf.WriteByte('\n')
	}

	return buf.Bytes()
}
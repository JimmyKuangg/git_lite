package core

import "errors"

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
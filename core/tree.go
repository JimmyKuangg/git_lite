package core

import (
	"bytes"
	"errors"
	"sort"
	"strings"
)

type TreeNode struct {
  Name string
  Files []TreeEntry
  Children map[string]*TreeNode
}

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

// func BuildTree(path string) (string, error) {
// 	tree := Tree{}
// 	entries, err := os.ReadDir(path)
// 	if (err != nil) {
// 		return "", fmt.Errorf("error reading directory, %w", err)
// 	}

// 	for _, entry := range entries {
// 		fullPath := filepath.Join(path, entry.Name())

// 		if !entry.IsDir() {
// 			content, err := os.ReadFile(fullPath)
// 			if err != nil {
// 				return "", fmt.Errorf("error reading file, %w", err)
// 			}

// 			blobHash, err := WriteObject(content)
// 			if err != nil {
// 					return "", err
// 			}

// 			tree.Add(entry.Name(), BlobType, blobHash )
// 		} else {
// 			subTreeHash, err := BuildTree(fullPath)
// 			if err != nil {
// 				return "", fmt.Errorf("error building tree, %w", err)
// 			}
// 			tree.Add(entry.Name(), TreeType, subTreeHash)
// 		}
// 	}

// 	treeBytes := tree.Encode()
// 	hash, err := WriteObject(treeBytes)
// 	if err != nil {
// 			return "", err
// 	}
// 	return hash, nil
// }

func BuildTree() (*TreeNode, error) {
	tree := &TreeNode{
    Children: make(map[string]*TreeNode),
	}

	index, err := ReadIndex()
	if err != nil {
		return nil, err
	}

	for path, hash := range index.Entries {
		pathSlice := strings.Split(path, "/")

		dirs := pathSlice[:len(pathSlice)-1]
		file := pathSlice[len(pathSlice)-1]

		current := tree

		for _, dir := range dirs {
			if current.Children[dir] == nil {
				current.Children[dir] = &TreeNode{
					Children: make(map[string]*TreeNode),
				}
			}

			current = current.Children[dir]
		}

		current.Files = append(current.Files, TreeEntry{
			Name: file,
			Type: BlobType,
			Hash: hash,
		})
	}

	return tree, nil
}
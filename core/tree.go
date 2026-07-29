package core

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
)

type TreeNode struct {
	Files    []TreeEntry
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

func (t *Tree) Add(entry TreeEntry) {
	t.Entries = append(t.Entries, entry)
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

func WriteTree(node *TreeNode) (string, error) {
	tree := Tree{}

	for _, file := range node.Files {
		tree.Add(file)
	}

	for name, child := range node.Children {
		subTreeHash, err := WriteTree(child)
		if err != nil {
			return "", fmt.Errorf("error building tree, %w", err)
		}

		tree.Add(TreeEntry{
			Name: name,
			Hash: subTreeHash,
			Type: TreeType,
		})
	}

	treeBytes := tree.Encode()
	hash, err := WriteObject(treeBytes)
	if err != nil {
		return "", err
	}

	return hash, nil
}

func FlattenTree(treeHash string) (map[string]string, error) {
	flattened := make(map[string]string)

	treeBytes, err := ReadObject(treeHash)
	if err != nil {
		return nil, err
	}

	treeData := string(treeBytes)
	dataSlice := strings.Split(strings.TrimSpace(treeData), "\n")

	for _, entry := range dataSlice {
		entryParts := strings.Split(entry, " ")
		entryType := entryParts[0]
		entryName := entryParts[1]
		entryHash := entryParts[2]

		if entryType == BlobType {
			flattened[entryName] = entryHash
		} else {
			flattenedSub, err := FlattenTree(entryHash)
			if err != nil {
				return nil, err
			}

			for key, flattenedEntry := range flattenedSub {
				flattened[entryName+"/"+key] = flattenedEntry
			}
		}
	}

	return flattened, nil
}

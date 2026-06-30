package commands

import "git_lite/core"

func Add(path string) error {
	root, err := core.EnsureGitliteRepo()
	if err != nil {
		return err
	}

	working, err := core.BuildWorkingSnapshot(root)
	if err != nil {
		return err
	}

	index, err := core.ReadIndex()
	if err != nil {
		return err
	}

	if path == "." {
		for file, hash := range working {
			index.Add(file, hash)
		}

		for file := range index.Entries {
			if _, exists := working[file]; !exists {
				index.Remove(file)
			}
		}

		return index.Save()
	}

	hash, exists := working[path]

	if exists {
		index.Add(path, hash)
		return index.Save()
	}

	index.Remove(path)
	return index.Save()
}
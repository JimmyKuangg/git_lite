package commands

import (
	"git_lite/core"
	"os"
	"path/filepath"
)

func Add(path string) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		fullPath := filepath.Join(path, entry.Name())

		if entry.IsDir() {
				err := Add(fullPath)
				if err != nil {
						return err
				}

				continue
		}

		content, err := os.ReadFile(fullPath)
		if err != nil {
			return err
		}

		_, err = core.WriteObject(content)
		if err != nil {
			return err
		}
	}

	return nil
}
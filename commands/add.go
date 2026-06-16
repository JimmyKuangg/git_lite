package commands

import (
	"git_lite/core"
	"os"
	"path/filepath"
)

func Add(path string) error {
    info, err := os.Stat(path)
    if err != nil {
      return err
    }

    if !info.IsDir() {
      return stageFile(path)
    }

    entries, err := os.ReadDir(path)
    if err != nil {
      return err
    }

    for _, entry := range entries {
      fullPath := filepath.Join(path, entry.Name())
      err := Add(fullPath)
      
			if err != nil {
          return err
      }
    }

    return nil
}

func stageFile(path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	
	hash, err := core.WriteObject(content)
	if err != nil {
		return err
	}

	index, err := core.ReadIndex()
  if err != nil {
      return err
  }

  index.Add(path, hash)

  return index.Save()
}
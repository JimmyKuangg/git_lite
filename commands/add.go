package commands

import (
	"git_lite/core"
	"os"
	"path/filepath"
)

func Add(path string) error {
  _, err := core.EnsureGitliteRepo()
	if err != nil {
		return err
	}

  index, err := core.ReadIndex()
  if err != nil {
    return err
  }

  err = addHelper(path, index)
  if err != nil && !os.IsNotExist(err) {
    return err
  }

  if path == "." {
    for path := range index.Entries {
      if _, err := os.Stat(path); os.IsNotExist(err) {
        index.Remove(path)
      }
    } 
  } else {
    if _, err := os.Stat(path); os.IsNotExist(err) {
      index.Remove(path)
    }
  }

  return index.Save()
}

func stageFile(path string, index *core.Index) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	
	hash, err := core.WriteObject(content)
	if err != nil {
		return err
	}

  index.Add(path, hash)
  return nil
}

func addHelper(path string, index *core.Index) error {
  info, err := os.Stat(path)
  if err != nil {
    return err
  }

  if !info.IsDir() {
    return stageFile(path, index)
  }

  entries, err := os.ReadDir(path)
  if err != nil {
    return err
  }

  for _, entry := range entries {
    if entry.IsDir() {
      if entry.Name() == ".gitlite" || entry.Name() == ".git" {
        continue
      }
    }

    fullPath := filepath.Join(path, entry.Name())
    err := addHelper(fullPath, index)
    
		if err != nil {
      return err
    }
  }

  return nil
}
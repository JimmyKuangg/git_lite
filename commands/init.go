package commands

import (
	"fmt"
	"os"
)

func Init() error {
err := os.MkdirAll(".gitlite/objects", 0755)
	if err != nil {
		return fmt.Errorf("error creating objects directory: %w", err)
	}

	err = os.MkdirAll(".gitlite/refs", 0755)
	if err != nil {
		return fmt.Errorf("error creating refs directory: %w", err)
	}

	err = os.WriteFile(".gitlite/index", []byte{}, 0644)
	if err != nil {
		return fmt.Errorf("error creating index file: %w", err)
	}

	err = os.WriteFile(".gitlite/HEAD", []byte("ref: refs/heads/main"), 0644)
	if err != nil {
		return fmt.Errorf("error creating HEAD file: %w", err)
	}

  return nil
}
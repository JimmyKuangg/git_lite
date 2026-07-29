package core

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

func WriteHEAD(hash string) error {
	if hash == "" {
		return errors.New("cannot write empty HEAD")
	}

	err := os.WriteFile(".gitlite/HEAD", []byte(hash), 0o644)
	if err != nil {
		return fmt.Errorf("failed to write HEAD: %w", err)
	}

	return nil
}

func ReadHEAD() (string, error) {
	data, err := os.ReadFile(".gitlite/HEAD")
	if err != nil {
		return "", fmt.Errorf("failed to read HEAD: %w", err)
	}

	return strings.TrimSpace(string(data)), nil
}

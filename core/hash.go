package core

import (
	"crypto/sha1"
	"fmt"
)

func Hash(content []byte) string {
	sum := sha1.Sum(content)
	return fmt.Sprintf("%x", sum)
}

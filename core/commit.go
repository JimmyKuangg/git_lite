package core

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Commit struct {
	Timestamp time.Time
	Author 		string
	Message 	string
	Root 			string
	Parent 		string
}

func BuildCommit(message string, root string) (Commit, error) {
	author := os.Getenv("USER")
	if author == "" {
		author = "unknown"
	}

	head, err := ReadHEAD() 
	if err != nil {
		return Commit{}, err
	}

	if strings.HasPrefix(head, "ref: ") {
    head = ""
	}

	commit := Commit{
		Author: author,
		Message: message,
		Root: root,
		Parent: head,
		Timestamp: time.Now(),
	}

	return commit, nil
}

func (c *Commit) Encode() []byte {
	var buf bytes.Buffer

	buf.WriteString("tree ")
	buf.WriteString(c.Root)
	buf.WriteByte('\n')

	if c.Parent != "" {
		buf.WriteString("parent ")
		buf.WriteString(c.Parent)
		buf.WriteByte('\n')
	}

	buf.WriteString("author ")
	buf.WriteString(c.Author)
	buf.WriteByte('\n')

	buf.WriteString("timestamp ")
	buf.WriteString(c.Timestamp.UTC().Format(time.RFC3339))
	buf.WriteByte('\n')

	buf.WriteString("message ")
	buf.WriteString(c.Message)
	return buf.Bytes()
}

func ParseCommit(data []byte) (Commit, error) {
	stream := string(data)
	lines := strings.Split(stream, "\n")

	foundRoot := false
	foundAuthor := false
	foundTimestamp := false
	foundMessage := false

	var c Commit

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		split := strings.SplitN(line, " ", 2)
		if len(split) < 2 {
			continue 
		}

		key := split[0]
		val := split[1]

		switch key {
		case "author":
			c.Author = val
			foundAuthor = true
		case "message":
			c.Message = val
			foundMessage = true
		case "tree":
			c.Root = val
			foundRoot = true
		case "parent":
			c.Parent = val
		case "timestamp":
			t, err := time.Parse(time.RFC3339, val)
			if err != nil {
				return Commit{}, err
			}
			c.Timestamp = t
			foundTimestamp = true
		}
	}

	if !foundRoot || !foundAuthor || !foundTimestamp || !foundMessage {
    return Commit{}, errors.New("invalid commit object")
	}

	return c, nil
}

func ReadCommit(hash string) (Commit, error) {
  data, err := ReadObject(hash)
  if err != nil {
    return Commit{}, err
  }

  return ParseCommit(data)
}

func ApplySnapshot(snapshot map[string]string) error {
	for path, hash := range snapshot {
		data, err := ReadObject(hash)
		if err != nil {
			return err
		}

		err = os.MkdirAll(filepath.Dir(path), 0755)
		if err != nil {
			return fmt.Errorf("error creating directory in ApplySnapshot: %w", err)
		}

		err = os.WriteFile(path, data, 0644) 
		if err != nil {
			return fmt.Errorf("error writing to file in ApplySnapshot: %w", err)
		}
	}

	return nil
}

func CleanupRemovedFiles(working, snapshot map[string]string) error {
	for path := range working {
    if _, exists := snapshot[path]; !exists {
      err := os.Remove(path)
			if err != nil {
				return fmt.Errorf("error removing file at %s: %w", path, err)
			}
    }
	}
	
	return nil
}
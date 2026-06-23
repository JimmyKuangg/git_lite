package core

import (
	"bytes"
	"os"
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

	buf.WriteByte('\n')

	buf.WriteString(c.Message)
	return buf.Bytes()
}

func ReadCommit(bytes []byte) (Commit, error) {
	stream := string(bytes)
	lines := strings.Split(stream, "\n")

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
		case "message":
			c.Message = val
		case "root":
			c.Root = val
		case "parent":
			c.Parent = val
		case "timestamp":
			t, err := time.Parse(time.RFC3339, val)
			if err != nil {
				return Commit{}, err
			}
			c.Timestamp = t
		}
	}

	return c, nil
}
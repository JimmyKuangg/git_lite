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
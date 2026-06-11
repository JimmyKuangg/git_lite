package core

import (
	"bytes"
	"time"
)

type Commit struct {
	Timestamp time.Time
	Author 		string
	Message 	string
	Root 			string
	Parent 		string
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
package core

type IndexEntry struct {
	Path string
	Hash string
}

type Index struct {
	Entries map[string]string
}

func (i *Index) Add(path string, hash string) {
	i.Entries[path] = hash
}
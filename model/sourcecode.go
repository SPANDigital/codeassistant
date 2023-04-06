package model

import (
	"os"
	"path/filepath"
)

type SourceCode struct {
	Filename string
	Language string
	Content  string
}

func (c SourceCode) Save(directory string) error {
	println("Saving ", c.Filename, "in", directory)
	return os.WriteFile(filepath.Join(directory, c.Filename), []byte(c.Content), 0644)
}

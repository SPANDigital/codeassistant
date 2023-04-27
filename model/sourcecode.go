// SPDX-License-Identifier: MIT

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
	return os.WriteFile(filepath.Join(directory, c.Filename), []byte(c.Content), 0644)
}

func (c SourceCode) String() string {
	return c.Content
}

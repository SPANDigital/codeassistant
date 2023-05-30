// SPDX-License-Identifier: MIT

package model

import (
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func BuildLibraries() map[string]*Library {

	libraries := make(map[string]*Library)

	var promptsLibrary string
	if viper.GetString("promptsLibraryDir") == "" {

		// Find home directory.
		home, err := os.UserHomeDir()
		if err != nil {
			return libraries
		}

		promptsLibrary = filepath.Join(home, "prompts-library")
	} else {
		promptsLibrary = viper.GetString("promptsLibraryDir")
	}

	libraryFromDir := func(path string) *Library {
		relPath, err := filepath.Rel(promptsLibrary, path)
		if err != nil {
			return nil
		}
		if relPath[0:1] == "." {
			return nil
		}
		library, found := libraries[relPath]
		if found {
			return library
		} else {
			library := &Library{
				Name:     relPath,
				FullPath: path,
				Index:    "",
				Commands: make(map[string]*Command),
			}
			libraries[relPath] = library
			return library
		}

	}

	_ = filepath.WalkDir(promptsLibrary, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		base := filepath.Base(path)
		frontName := strings.Split(base, ".")[0]
		ext := filepath.Ext(path)
		if d.IsDir() && d.Name()[0:1] != "." {
			_ = libraryFromDir(path)
		} else if !d.IsDir() && d.Name() == "_index.md" {
			library := libraryFromDir(filepath.Dir(path))
			if library != nil {
				data, err := os.ReadFile(path)
				if err == nil {
					library.Index = string(data)
				}
			}
		} else if !d.IsDir() && (ext == ".yml" || ext == ".yaml" || ext == ".js") {
			library := libraryFromDir(filepath.Dir(path))
			if library != nil {
				data, err := os.ReadFile(path)
				if err == nil {
					if base == "_index.yml" {
						err = yaml.Unmarshal(data, &library)
						return err
					} else if base == "_data.yml" {
						err = yaml.Unmarshal(data, &library.Data)
						return err
					} else if ext == ".yaml" || ext == ".yml" {
						command := library.getCommand(frontName)
						err := yaml.Unmarshal(data, &command)
						if err != nil {
							return err
						}
						if command.DisplayName == "" {
							command.DisplayName = strings.ReplaceAll(command.Name, "-", " ")
						}
					} else if ext == ".js" {
						command := library.getCommand(frontName)
						command.Script = string(data)
					}
				}
			}
		}
		return nil
	})
	return libraries
}

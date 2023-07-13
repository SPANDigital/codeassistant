package embeddings

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
	embeddingsLibrary := viper.GetString("embeddingsLibraryDir")
	if embeddingsLibrary == "" {
		return libraries
	}
	libraryFromDir := func(path string) *Library {
		relPath, err := filepath.Rel(embeddingsLibrary, path)
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
				Name:       relPath,
				Path:       path,
				Index:      "",
				Embeddings: make(map[string]*Embedding),
			}
			libraries[relPath] = library
			return library
		}
	}

	_ = filepath.WalkDir(embeddingsLibrary, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		base := filepath.Base(path)
		ext := filepath.Ext(path)
		if d.IsDir() && d.Name()[0:1] != "." {
			_ = libraryFromDir(path)
		} else if !d.IsDir() && d.Name() == "_index.md" {
			library := libraryFromDir(filepath.Dir(path))
			if library != nil {
				data, err := os.ReadFile(path)
				if err == nil {
					library.Index = string(data)
					library.addBuildPath(path)
				}
			}
		} else if !d.IsDir() && (ext == ".yml" || ext == ".yaml" || ext == ".js") {
			library := libraryFromDir(filepath.Dir(path))
			if library != nil {
				data, err := os.ReadFile(path)
				if err == nil {
					if base == "_index.yml" {
						err = yaml.Unmarshal(data, &library)
						if err != nil {
							library.addBuildPath(path)
						}
						return err
					} else if base == "_data.yml" {
						err = yaml.Unmarshal(data, &library.Data)
						if err != nil {
							library.addBuildPath(path)
						}
						return err
					} else if ext == ".yaml" || ext == ".yml" {
						embedding := library.getEmbedding(path)
						err := yaml.Unmarshal(data, &embedding)
						if err != nil {
							return err
						}
						if embedding.DisplayName == "" {
							embedding.DisplayName = strings.ReplaceAll(embedding.Name, "-", " ")
						}
					} else if ext == ".js" {
						embedding := library.getEmbedding(path)
						embedding.Script = string(data)
					}
				}
			}
		}
		return nil
	})
	return libraries
}

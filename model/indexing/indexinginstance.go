package indexing

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Masterminds/sprig/v3"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/robertkrimen/otto"
	"github.com/spandigitial/codeassistant/client"
	"github.com/spf13/viper"
	"github.com/yargevad/filepathx"
	"io"
	_ "modernc.org/sqlite"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"text/template"
	"time"
)

type IndexingInstance struct {
	Indexing         *Indexing
	Params           map[string]string
	EnableStdin      bool
	Data             map[string]interface{}
	WorkingDirectory string
	client           client.LLMClient
	Index            *Index
}

type TemplateFunc func() (string, error)

func stdin() (string, error) {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func fakeStdin() (string, error) {
	return "", nil
}

var extensionToLanguage = map[string]string{
	".rb":     "Ruby",
	".java":   "Java",
	".python": "Python",
	".js":     "Javascript",
	".ts":     "Typescript",
}

func (ei *IndexingInstance) combineParamsData() map[string]interface{} {
	combo := make(map[string]interface{})
	for k, v := range ei.Params {
		combo[k] = v
	}
	for k, v := range ei.Data {
		combo[k] = v
	}
	combo["WorkingDirectory"] = ei.WorkingDirectory
	combo["Collection"] = ei.Index
	return combo
}

func (ei *IndexingInstance) runRowTemplate(input string, row map[string]interface{}) (string, error) {

	tmpl, err := template.New("paramTemplate").Funcs(sprig.FuncMap()).Funcs(map[string]any{
		"stdin": (map[bool]TemplateFunc{true: stdin, false: fakeStdin})[ei.EnableStdin],
	}).Parse(input)
	if err != nil {
		return "", err
	}
	buff := bytes.Buffer{}
	err = tmpl.Execute(&buff, row)
	if err != nil {
		return "", err
	}
	return buff.String(), nil
}

func (ei *IndexingInstance) runParamTemplate(input string) (string, error) {

	tmpl, err := template.New("paramTemplate").Funcs(sprig.FuncMap()).Funcs(map[string]any{
		"stdin": (map[bool]TemplateFunc{true: stdin, false: fakeStdin})[ei.EnableStdin],
	}).Parse(input)
	if err != nil {
		return "", err
	}
	buff := bytes.Buffer{}
	err = tmpl.Execute(&buff, ei.combineParamsData())
	if err != nil {
		return "", err
	}
	return buff.String(), nil
}

type onDemandFile struct {
	path string
}

func (o onDemandFile) String() string {
	b, _ := os.ReadFile(o.path)
	return string(b)
}

func (o onDemandFile) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.String())
}

func (ei *IndexingInstance) runTemplate(content string) (string, error) {
	tmpl, err := template.New("runTemplate").Funcs(sprig.FuncMap()).Funcs(template.FuncMap{
		"FilepathJoin": func(elems ...string) string {
			return filepath.Join(elems...)
		},
		"FilepathGlob": func(pattern string) []string {
			matches, _ := filepathx.Glob(pattern)
			return matches
		},
		"SetIndex": func(c []map[string]interface{}) string {
			ei.Index = NewIndex(c)
			return ""
		},
		"ExtractFiles": func(ps []string) ([]map[string]interface{}, error) {
			extracted := make([]map[string]interface{}, len(ps))
			for idx, p := range ps {

				info, err := os.Stat(p)
				if err != nil {
					return nil, err
				}

				extracted[idx] = map[string]interface{}{
					"Name":     info.Name(),
					"Path":     strings.TrimPrefix(strings.TrimPrefix(p, ei.WorkingDirectory), "/"),
					"Content":  onDemandFile{path: p},
					"Language": extensionToLanguage[filepath.Ext(p)],
					"Size":     info.Size(),
					"ModTime":  info.ModTime(),
				}
				for k, v := range ei.Indexing.ExtractFilesEnrich {
					result, err := ei.runRowTemplate(v, extracted[idx])
					if err != nil {
						return nil, err
					}
					s := fmt.Sprintf("cd %s; %s", ei.WorkingDirectory, result)
					fmt.Println("Calling command", s)
					output, err := exec.Command("/bin/sh", "-c", s).Output()
					if err != nil {
						return nil, err
					}
					extracted[idx][k] = strings.TrimSpace(string(output))
				}
			}
			return extracted, nil
		},
	}).Parse(content)
	if err != nil {
		return "", err
	}
	buff := bytes.Buffer{}
	err = tmpl.Execute(&buff, ei.combineParamsData())
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(buff.String()), nil
}

func (ei *IndexingInstance) runVectorTemplates(item IndexItem) (map[string]string, error) {

	results := map[string]string{}

	for k, v := range ei.Indexing.Vectors {
		tmpl, err := template.New("inputTemplate").Funcs(sprig.FuncMap()).Funcs(map[string]any{
			"stdin": (map[bool]TemplateFunc{true: stdin, false: fakeStdin})[ei.EnableStdin],
		}).Parse(v)
		if err != nil {
			return nil, err
		}
		buff := bytes.Buffer{}
		err = tmpl.Execute(&buff, item.Content)
		if err != nil {
			return nil, err
		}
		results[k] = strings.TrimSpace(buff.String())
	}
	return results, nil

}

func (ei *IndexingInstance) buildParams(defaultParams map[string]string, args []string) {
	params := make(map[string]string)
	for k, v := range ei.Indexing.Params {
		value, found := defaultParams[k]
		if found {
			params[k] = value
		} else {
			value, err := ei.runParamTemplate(v)
			if err == nil {
				params[k] = value
			}
		}
	}
	for _, arg := range args {
		parts := strings.Split(arg, ":")
		if len(parts) == 2 {
			params[parts[0]] = parts[1]
		}
	}
	ei.Params = params
}

func (ei *IndexingInstance) runScript() error {
	paramsJson, err := json.Marshal(ei.Params)
	if err != nil {
		return err
	}
	dataJson, err := json.Marshal(ei.Data)
	if err != nil {
		return err
	}
	if ei.Indexing.Script != "" {
		vm := otto.New()

		_, err := vm.Run(fmt.Sprintf("var data=%s;\nvar params=%s;\n%s",
			string(dataJson),
			string(paramsJson),
			ei.Indexing.Script))
		if err != nil {
			return err
		}
	}
	return nil
}

func NewEmbeddingInstance(client client.LLMClient, enableStdin bool, defaultParams map[string]string, args ...string) (*IndexingInstance, error) {
	if len(args) < 2 {
		return nil, errors.New("at least two arguments are required to construct a embedding instance")
	}
	libraryName, embeddingName := args[0], args[1]
	library, found := BuildLibraries()[libraryName]
	if !found {
		return nil, fmt.Errorf("library: '%s' not found", libraryName)
	}
	embedding, found := library.Embeddings[embeddingName]
	if !found {
		return nil, fmt.Errorf("indexing: '%s' not found in library: '%s", embeddingName, libraryName)
	}
	var workingDir string
	var err error
	if embedding.WorkingDirectory == "" {
		workingDir, err = os.MkdirTemp("", "embedding-working")
		if err != nil {
			return nil, err
		}
	} else {
		workingDir = embedding.WorkingDirectory
		_, err := os.Stat(workingDir)
		if err != nil {
			os.MkdirAll(workingDir, 0700)
		}
	}

	embeddingInstance := &IndexingInstance{
		EnableStdin:      enableStdin,
		Data:             library.Data,
		Indexing:         embedding,
		WorkingDirectory: workingDir,
		Index:            NewIndex(nil),
		client:           client,
	}
	embeddingInstance.buildParams(defaultParams, args[2:])
	return embeddingInstance, nil
}

func (ei *IndexingInstance) Collect() error {
	if ei.Indexing.Clone != "" {
		cloneOptions := &git.CloneOptions{
			URL:      ei.Indexing.Clone,
			Progress: os.Stdout,
		}

		if ei.Indexing.AccessToken != "" {
			cloneOptions.Auth = &http.BasicAuth{
				Username: "abc123", //does not matter
				Password: ei.Indexing.AccessToken,
			}
		}
		fmt.Fprintf(os.Stdout, "Cloning to %s\n", ei.WorkingDirectory)

		_, err := git.PlainClone(ei.WorkingDirectory, false, cloneOptions)
		if err != nil && err.Error() != "repository already exists" {
			return err
		}

	}
	if ei.Index.Length() == 0 {
		_, err := ei.runTemplate(ei.Indexing.Collect)
		if err != nil {
			return err
		}
	}

	return nil
}

type job struct {
	Index int
}

type jobResult struct {
	Index      int
	KeyVectors map[string][]float64
}

const numberOfWorkers = 10

func (ei *IndexingInstance) Fetch() error {
	start := time.Now()

	var wg sync.WaitGroup
	wg.Add(numberOfWorkers)
	jobChannel := make(chan job)
	jobResultChannel := make(chan jobResult, ei.Index.Length())

	fmt.Println("Workers starting")

	// Start the workers
	for i := 0; i < numberOfWorkers; i++ {
		go ei.worker(i, &wg, jobChannel, jobResultChannel)
	}

	// Send jobs to worker
	for idx, _ := range ei.Index.Items {
		jobChannel <- job{
			Index: idx,
		}

	}

	close(jobChannel)
	wg.Wait()
	close(jobResultChannel)

	// Receive job results from workers
	for result := range jobResultChannel {
		fmt.Println("Results received for row", result.Index, len(result.KeyVectors))
		for k, v := range result.KeyVectors {
			ei.Index.Items[result.Index].Vectors[k] = v
		}
	}

	fmt.Printf("Took %s\n", time.Since(start))
	return nil
}

func (ei *IndexingInstance) worker(id int, wg *sync.WaitGroup, jobChannel <-chan job, resultChannel chan jobResult) {
	defer wg.Done()
	for job := range jobChannel {
		resultChannel <- ei.work(id, job)
	}
}

func (ei *IndexingInstance) work(workerId int, job job) jobResult {
	results, err := ei.runVectorTemplates(ei.Index.Items[job.Index])
	keyVectors := map[string][]float64{}
	if err == nil {
		for k, v := range results {
			vector, err := ei.client.Embeddings(viper.GetString("openAiEmbeddingsModel"), v)
			if err == nil {
				keyVectors[k] = vector
			} else {
				fmt.Fprintln(os.Stderr, "Error on row", job.Index, err)
			}
		}
	}
	return jobResult{
		Index:      job.Index,
		KeyVectors: keyVectors,
	}
}

func (ei *IndexingInstance) Load() error {
	if ei.Indexing.LoadIfPossible {
		bytes, err := os.ReadFile(ei.Indexing.IndexJson)
		if err != nil {
			if err.Error() != "no such file or directory" {
				return nil
			}
			return err
		}
		json.Unmarshal(bytes, &ei.Index)
		return err
	}
	return nil
}

func (ei *IndexingInstance) Save() error {
	if ei.Indexing.SaveIfPossible {
		bytes, err := json.Marshal(ei.Index)
		if err != nil {
			return err
		}
		err = os.WriteFile(ei.Indexing.IndexJson, bytes, 0644)
		return err
	}
	return nil
}

func (ei *IndexingInstance) Search(llmClient client.LLMClient, content string, query string, messageParts chan<- client.MessagePart) error {
	queryVector, err := ei.client.Embeddings(viper.GetString("openAiEmbeddingsModel"), query)
	if err != nil {
		return err
	}
	closest, err := ei.Index.Closest(content, 5, queryVector)
	if err != nil {
		return err
	}
	var context string
	for _, i := range closest {
		context = context + i.Content["Content"].(string) + "###"
	}

	question := "Answer the question based on the context below, and if the question can't be answered based on the context, say \"I don't know\"\n\nContext: " + context + "\n\n---\n\nQuestion: " + query + "\nAnswer:"

	llmClient.SimpleCompletion("gpt-4", ei.Indexing.RoleHint, question, messageParts)

	return nil
}

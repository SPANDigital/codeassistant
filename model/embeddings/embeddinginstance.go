package embeddings

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
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"
	"time"
)

type EmbeddingInstance struct {
	Embedding        *Embedding
	Params           map[string]string
	EnableStdin      bool
	Data             map[string]interface{}
	WorkingDirectory string
	client           client.LLMClient
	collection       []*map[string]interface{}
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

func (ei *EmbeddingInstance) combineParamsData() map[string]interface{} {
	combo := make(map[string]interface{})
	for k, v := range ei.Params {
		combo[k] = v
	}
	for k, v := range ei.Data {
		combo[k] = v
	}
	combo["WorkingDirectory"] = ei.WorkingDirectory
	combo["Collection"] = ei.collection
	return combo
}

func (ei *EmbeddingInstance) runParamTemplate(input string) (string, error) {

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

func (ei *EmbeddingInstance) runTemplate(content string) (string, error) {
	tmpl, err := template.New("runTemplate").Funcs(sprig.FuncMap()).Funcs(template.FuncMap{
		"FilepathJoin": func(elems ...string) string {
			return filepath.Join(elems...)
		},
		"FilepathGlob": func(pattern string) []string {
			matches, _ := filepathx.Glob(pattern)
			return matches
		},
		"SetCollection": func(c []*map[string]interface{}) string {
			ei.collection = c
			return ""
		},
		"ExtractFiles": func(ps []string) []*map[string]interface{} {
			extracted := make([]*map[string]interface{}, len(ps))
			for idx, p := range ps {
				extracted[idx] = &map[string]interface{}{
					"Path":    p,
					"Content": onDemandFile{path: p},
				}
			}
			return extracted
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
	return buff.String(), nil
}

func (ei *EmbeddingInstance) runInputTemplate(row *map[string]interface{}) (string, error) {

	tmpl, err := template.New("inputTemplate").Funcs(sprig.FuncMap()).Funcs(map[string]any{
		"stdin": (map[bool]TemplateFunc{true: stdin, false: fakeStdin})[ei.EnableStdin],
	}).Parse(ei.Embedding.Input)
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

func (ei *EmbeddingInstance) buildParams(defaultParams map[string]string, args []string) {
	params := make(map[string]string)
	for k, v := range ei.Embedding.Params {
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

func (ei *EmbeddingInstance) runScript() error {
	paramsJson, err := json.Marshal(ei.Params)
	if err != nil {
		return err
	}
	dataJson, err := json.Marshal(ei.Data)
	if err != nil {
		return err
	}
	if ei.Embedding.Script != "" {
		vm := otto.New()

		_, err := vm.Run(fmt.Sprintf("var data=%s;\nvar params=%s;\n%s",
			string(dataJson),
			string(paramsJson),
			ei.Embedding.Script))
		if err != nil {
			return err
		}
	}
	return nil
}

func NewEmbeddingInstance(client client.LLMClient, enableStdin bool, defaultParams map[string]string, args ...string) (*EmbeddingInstance, error) {
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
		return nil, fmt.Errorf("embeddings: '%s' not found in library: '%s", embeddingName, libraryName)
	}
	workingDir, err := os.MkdirTemp("", "embedding-working")
	if err != nil {
		return nil, err
	}
	embeddingInstance := &EmbeddingInstance{
		EnableStdin:      enableStdin,
		Data:             library.Data,
		Embedding:        embedding,
		WorkingDirectory: workingDir,
		client:           client,
	}
	embeddingInstance.buildParams(defaultParams, args[2:])
	return embeddingInstance, nil
}

func (ei *EmbeddingInstance) Collect() error {
	if ei.Embedding.Clone != "" {
		cloneOptions := &git.CloneOptions{
			URL:      ei.Embedding.Clone,
			Progress: os.Stdout,
		}

		if ei.Embedding.AccessToken != "" {
			cloneOptions.Auth = &http.BasicAuth{
				Username: "abc123", //does not matter
				Password: ei.Embedding.AccessToken,
			}
		}
		fmt.Fprintf(os.Stdout, "Cloning to %s\n", ei.WorkingDirectory)

		_, err := git.PlainClone(ei.WorkingDirectory, false, cloneOptions)
		if err != nil {
			return err
		}

	}
	_, err := ei.runTemplate(ei.Embedding.Build)
	if err != nil {
		return err
	}

	return nil
}

type job struct {
	Index int
}

type jobResult struct {
	Index  int
	Error  error
	Vector []float32
}

const numberOfWorkers = 10

func (ei *EmbeddingInstance) Fetch() error {
	start := time.Now()

	var jobs []job
	for i := 0; i < len(ei.collection); i++ {
		jobs = append(jobs, job{Index: i})
	}

	var wg sync.WaitGroup
	wg.Add(numberOfWorkers)
	jobChannel := make(chan job)
	jobResultChannel := make(chan jobResult, len(jobs))

	// Start the workers
	for i := 0; i < numberOfWorkers; i++ {
		go ei.worker(i, &wg, jobChannel, jobResultChannel)
	}

	// Send jobs to worker
	for _, job := range jobs {
		jobChannel <- job
	}

	close(jobChannel)
	wg.Wait()
	close(jobResultChannel)

	// Receive job results from workers
	for result := range jobResultChannel {
		if result.Error == nil {
			fmt.Fprintln(os.Stdout, result.Index, "->", result.Vector)
			(*ei.collection[result.Index])["vector"] = result.Vector
		} else {
			fmt.Fprintln(os.Stderr, "No vector for index", result.Index, result.Error)
		}
	}

	fmt.Printf("Took %s", time.Since(start))
	return nil
}

func (ei *EmbeddingInstance) worker(id int, wg *sync.WaitGroup, jobChannel <-chan job, resultChannel chan jobResult) {
	defer wg.Done()
	for job := range jobChannel {
		resultChannel <- ei.work(id, job)
	}
}

func (ei *EmbeddingInstance) work(workerId int, job job) jobResult {
	input, err := ei.runInputTemplate(ei.collection[job.Index])
	vector, err := ei.client.Embeddings(viper.GetString("openAiEmbeddingsModel"), input)
	return jobResult{
		Index:  job.Index,
		Error:  err,
		Vector: vector,
	}
}

# codeassistant

codeassistant is toolinf which automates interactions with the
[OpenAI Completions API](https://platform.openai.com/docs/api-reference/completions) and
the [Vertex AI Predict API](https://cloud.google.com/vertex-ai/docs/predictions/get-predictions).

Prompts are organized in a directory (or _prompts library_) as YAML configuration files with documentation implemented
in Markdown.
An example of such a library can be found [here](https://github.com/SPANDigital/prompts-library).

We are looking for contributors, please see how you can [contribute](CONTRIBUTING.md), and our [code of conduct](CODE_OF_CONDUCT.md).

It fulfills these purposes:

- A tool for prompt engineers to prototype prompts, and rapidly iterate on them
- The ability to parameterize prompts with light templating in the handling of input
- Allows prompts to be integrated with other software such as shell scripts
- Provides a Web UI

It has two main modes of operation:

- CLI: Suitable for shell scripts. Output of prompts can be redirected from STDOUT.
- WebUI: Useful for testing prompts.

## OpenAI Configuration

You will need to configure an OpenAI API Key before usage.

### Sample configuration

It is recommended you set up `codeassistant` with a config file at `$HOME/.codeassistant.yaml` for default values:

```yaml
backend: openai
openAiApiKey: "<api key>"
openAiUserId: "<your email address>"
promptsLibraryDir: <directory to load prompts, defaults to $HOME/prompts-library>
```

## VertexAI Configuration

You will need to install the gcloud sdk before using VertexAI.
You will need a project on Google Cloud which gives your user access to Vertex AI.

### Sample configuration

```yaml
backend: vertexai
vertexAiProjectId: "<project-id>"
vertexAiLocation: "us-central1"
vertexAiModel: "text-bison@001"
promptsLibraryDir: <directory to load prompts, defaults to $HOME/prompts-library>
```

## Additional configuration

More keys are available for debugging

```yaml
debug:
  - configuration
  - first-response-time
  - last-response-time
  - request-header
  - request-time
  - request-tokens
  - response-header
  - sent-prompt
  - webserver
```

## Installing and running via Docker

```bash
docker run --rm --name codeassistant \
  --volume $HOME/.codeassistant.yaml:/.codeassistant.yaml:ro \
  --volume $HOME/prompts-library:/prompts-library:ro \
  -p8989:8989  \
  ghcr.io/spandigital/codeassistant:latest serve
```

In this example `.codeassistant.yaml` is `$HOME/.codeassistant.yaml and`
prompts-library and the prompts-library folder is in `$HOME/.prompts-library`
On the docker container $HOME is defined as /

## Installing an running via MacOS X

### Initial installation

```bash
brew tap SPANDigital/homebrew-tap
brew install codeassistant
```

### Upgrades

```bash
brew up
brew reinstall codeassistant
```

## Usage

### Run web front-end

```bash
codeassistant serve
```

or to override the default model

```bash
codeassistant serve --defaultModel gpt-4
```

### List all the commands in your prompt libraries

```bash
codeassistant list
```

### List commands for a specific prompt library

```bash
codeassistant run <library> <command> <var1:value> <vae2:value>
```

or to override the default model

```bash
codeassistant run <library> <command> <var1:value> <vae2:value> --defaultModel gpt-4
```

### List available  models (beta)

```bash
codeassistant list-models
```


This `README.md` file is documentation:

`SPDX-License-Identifier: MIT`

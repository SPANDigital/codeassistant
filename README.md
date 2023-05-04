# codeassistant
Use ChatGPT API to automate programming tasks

## OpenAI API Key

codeassistant requires a OpenAI:wq API Key to use.

## Default values

It is recommended you set up codeassistant with a config file at `$HOME/.codeassistant.yaml` for default values

```yaml
openAiApiKey: <api key>
user: <your email adresss>
promptsLibraryDir: <directory to load prompts, defaults to $HOME/prompts-library>
```

## installation

Installation requires a personal access token with at least repo access in the HOMEBREW_GITHUB_API_TOKEN environment variable.

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

## List available comand in your prompt libraries

```bash
codeassistant list
```

## List commands in a prompt library

```bash
codeassistant run <library> <command> <var1:value> <vae2:value>
```

This `README.md` file is documentation:

`SPDX-License-Identifier: MIT`
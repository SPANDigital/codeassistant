# codeassistant
Use ChatGPT API to automate programming tasks

## OpenAI API Key

You will need to configure an OpenAI API Key before usage.

## Default values

It is recommended you set up `codeassistant` with a config file at `$HOME/.codeassistant.yaml` for default values:

```yaml
openAiApiKey: "<api key>"
userEmail: "<your email address>"
promptsLibraryDir: <directory to load prompts, defaults to $HOME/prompts-library>
```

More complex configurations are possible:

```yaml
openAiApiKey: "<api key>"
userEmail: "<your email address>"
promptsLibraryDir: <directory to load prompts, defaults to $HOME/prompts-library>
userAgent: "<use this for user agent header>"
debug:
  - configuration
  - first-response-time
  - last-response-time
<<<<<<< HEAD
  - request-header
  - request-time
  - request-tokens
  - response-header
  - sent-prompt
=======
  - webserver
>>>>>>> 9b2ac38 (Webserver)
```

## Installation

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

### List all the commands in your prompt libraries

```bash
codeassistant list
```

### List commands for a specific prompt library

```bash
codeassistant run <library> <command> <var1:value> <vae2:value>
```

This `README.md` file is documentation:

`SPDX-License-Identifier: MIT`
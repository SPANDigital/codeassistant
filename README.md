# codeassistant
Use ChatGPT API to automate programming tasks

## OpenAI API Key

codeassistant requires a OpenAI API Key to use.

## Default values

It is recommended you set up codeassistant with a config file at `$HOME/.codeassistant.yaml` for default values

```yaml
openAiApiKey: <api key>
user: <your email adresss>
```



## installation

Installation requires a personal access token with at least repo access in the HOMEBREW_GITHUB_API_TOKEN environment variable.

```bash
export HOMEBREW_GITHUB_API_TOKEN="<personal access token>"
brew tap SPANDigital/homebrew-tap
homebrew install codeassistant
```

## Use with Projects with Jetbrains IDEs

### Clone the code assistant repo

### Install in JetBrains IDE

- RubyMine

- WebStorm

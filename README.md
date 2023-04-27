# codeassistant
Use ChatGPT API to automate programming tasks

## OpenAI API Key

codeassistant requires a OpenAI:wq API Key to use.

## Default values

It is recommended you set up codeassistant with a config file at `$HOME/.codeassistant.yaml` for default values

```yaml
openAiApiKey: <api key>
user: <your email adresss>
```

## installation

Installation requires a personal access token with at least repo access in the HOMEBREW_GITHUB_API_TOKEN environment variable.

### Initial installation

```bash
export HOMEBREW_GITHUB_API_TOKEN="<personal access token>"
brew tap SPANDigital/homebrew-tap
brew install codeassistant
```

### Upgrades

```bash
export HOMEBREW_GITHUB_API_TOKEN="<personal access token>"
brew up
brew reinstall codeassistant
```

## Use with Projects with Jetbrains IDEs

Clone the code assistant repo

```bash
git clone git@github.com:SPANDigital/codeassistant.git
```

### Install in JetBrains IDE

- RubyMine
  1. Install the Flora plugin in Rubymine
  2. Craate a '.plugins' in a Rubymine project
  3. ```bash
     cp <codeassistentDirectory>/jetbrains/rubymine/*.plugin.js <projectDiretory>/.plugins
     ```
  4. Reload your RubyMine project   

- WebStorm
  1. Install the Flora plugin in Rubymine
  2. Craate a '.plugins' in a Rubymine project
  3. ```bash
     cp <codeassistentDirectory>/jetbrains/webstorm/*.plugin.js <projectDiretory>/.plugins
     ```
  4. Reload your Webstorm project

## Commands

### Retrieve article abiut something in markdown

```bash
codeassistant article <term>
```

### Generate NestJS entities and basic services from a Ruby on Rails Schema

```bash
codeassistant rails2nestjs schema2entities \
  --schemaFilename <pathToRubyOnRailsSchema> \
  --entitiesDirectory <directoryToSaveEntities> \
  --servicesDirectory <directoryToSaveServices> 
```

### Convert Ruby on Rails Controllers to NestJS Controllers

```bash
codeassistant rails2nestjs convert \
  --railstype controller \
  --nestjstype controller \
  --src <sourceToRubyFile> \
  --dest <destinationTypeScript>
```

### Convert Ruby on Rails RSpec to NestJS Tests using Jest

```bash
codeassistant rails2nestjs convert \
  --railstype spec \
  --nestjstype "test using jest" \ 
  --src <sourceToRubyFile> \
  --dest <destinationTypeScript>
```

This `README.md` file is documentation:

`SPDX-License-Identifier: MIT`
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

## Commands

### Retrieve definition of something markdown

```bash
codeassistant whatis <term>
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
# cfn-compose
A command-line tool for managing CloudFormation Stacks at scale.

By Balman Rawat while working at [CloudFactory](https://www.cloudfactory.com/)

## Features
* Create/Update/Delete multiple CloudFormation stacks parallely or sequentially
* Customize the CloudFormation stacks dependency using yaml config
* Delete multiple CloudFormation stacks respecting the creation sequence
* DryRun mode to plan the change
* Generate/Validate/visualize configuration with ease

![Demo](./docs/images/demo.gif)

## Overview
As the infrastrucutre evolves and gets complicated we need to manage/maintain multiple CloudFormation Stacks. When we want to `create/update/delete` these stacks we need to manually apply the actions one at a time. Deletion mostly in dev/test environment can be hectic as we should delete the stacks in the reverse of creation order. **cfn-compose** helps to manage multiple stacks that are closely related using declarative language.

![overview image](./docs/images/cfn-compose.svg)

## Usage
```shell
➜ cfn-compose --help          
Manage cloudformation stacks at scale. Design and deploy multiple cloudformation stacks either in sequence or in prallel using declarative configuration

Usage:
  cfn-compose [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  config      Generate, validate and visualize the compose configuration
  deploy      Deploys the stacks based on the sequence specified in the compose configuration
  destroy     Destroys all the stacks in the reverse order of creation
  help        Help about any command

Flags:
  -c, --config string     File path to compose file (default "cfn-compose.yml")
  -d, --dry-run           Run commands in dry run mode
  -h, --help              help for cfn-compose
  -l, --loglevel string   Specify Log Levels. Valid Levels are: DEBUG, INFO, WARN, ERROR (default "INFO")
  -v, --version           version for cfn-compose

Use "cfn-compose [command] --help" for more information about a command.
```

#### Examples
```shell
## Deploy
cfn-compose deploy
## Deploy in dry run mode
cfn-compose deploy -d

## Destroy
cfn-compose destroy
## Destroy in dry run mode
cfn-compose destroy -d

## Generate Validate and Visualize compose configuration
cfn-compose config generate
cfn-compose config validate
cfn-compose config visualize
```

## Parameters
| Command | Options | Description |
| ------------- | ------------- | ------------- |
| cfn-compose | -h, --help, help | Get description of cfn-compose |
| cfn-compose | -d, --dry-run | enable dry run mode |
| cfn-compose| -l, --loglevel |  Specify Log Levels. Valid Levels are: DEBUG, INFO, WARN, ERROR (default "INFO") |
| cfn-compose| -c, --config | File path to compose file (default "cfn-compose.yml") |
| cfn-compose deploy | with no flag| deploys all the stacks |
| cfn-compose deploy | -f, --flow | Cherry pick specific flow to deploy |
| cfn-compose destroy | with no flag| destroys all the stacks |
| cfn-compose destroy | -f, --flow | Cherry pick specific flow to destroy |
| cfn-compose config generate | no flags | Generates compose template |
| cfn-compose config validate | no flags | Validates the compose configuration |
| cfn-compose config visualize | no flags | Visualize the stacks dependencies and creation order |
| cfn-compose | -v, --version |  version for cfn-compose |

## Installation
Binary is available for Linux, Windows and Mac OS (amd64 and arm64). Download the binary for your respective platform from the [releases page](https://github.com/rbalman/cfn-compose/releases).

Linux:
```
curl -sSLO https://github.com/rbalman/cfn-compose/releases/download/v0.0.1-beta/cfn-compose-v0.0.1-beta-linux-amd64.tar.gz
```
```
tar zxf cfn-compose-v0.0.1-beta-linux-amd64.tar.gz
```
```
sudo install -m 0755 cfn-compose /usr/local/bin/cfn-compose
```

macOS (Intel):
```
curl -sSLO https://github.com/rbalman/cfn-compose/releases/download/v0.0.1-beta/cfn-compose-v0.0.1-beta-darwin-amd64.tar.gz
```
```
tar zxf cfn-compose-v0.0.1-beta-darwin-amd64.tar.gz
```
```
sudo install -m 0755 cfn-compose /usr/local/bin/cfn-compose
```

macOS (Apple Silicon):
```
curl -sSLO https://github.com/rbalman/cfn-compose/releases/download/v0.0.1-beta/cfn-compose-v0.0.1-beta-darwin-arm64.tar.gz
```
```
tar zxf cfn-compose-v0.0.1-beta-darwin-arm64.tar.gz
```
```
sudo install -m 0755 cfn-compose /usr/local/bin/cfn-compose
```

Windows:
```
curl -sSLO https://github.com/rbalman/cfn-compose/releases/download/v0.0.1-beta/cfn-compose-v0.0.1-beta-windows-amd64.zip
```
```
unzip cfn-compose-v0.0.1-beta-windows-amd64.zip
```

## Development

If you wish to contribute or compile from source code, you'll first need Go installed on your machine. Go version 1.18+ is required.

```
git clone https://github.com/rbalman/cfn-compose
cd cfn-compose 
go build
```

## Contributors
- [Balman Rawat](https://github.com/rbalman)

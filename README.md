# cfn-compose
A command-line tool for managing CloudFormation Stacks at scale.

By Balman Rawat, started as an experiment project at [CloudFactory](https://www.cloudfactory.com/)

## Features
* Create/Update/Delete multiple CloudFormation stacks parallely or sequentially
* Customize the CloudFormation stacks dependency using yaml config
* Delete multiple CloudFormation stacks respecting the creation sequence
* DryRun mode to plan the change
* Generate/Validate/visualize configuration with ease
* Supports Go Templating for dynamic value substitution

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

## Man
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

## Compose Configuration
**Syntax:**

```yaml
Description: Sample CloudFormation Compose file
Vars:
  Key: Value
Flows:
  Flow1:
    Order: 0
    Description: Flow1 Description
    Stacks:
    - Stack1
    - Stack2
  Flow2:
    Order: 1
    Description: Flow2 description
    Stacks:
    - Stack1
    - Stack2
```
A typical compose configuration contains:
- Optional `Description`
- Optional `Vars` section to define variables in `Key: Value` mapping
eg:

```yaml
Vars:
  ENV_TYPE: "nonproduction"
  ENV_NAME: "demo"
  AWS_PROFILE: "demo"
```
- Mandatory `Flows:` section
`Flow` is a collection of CloudFormation stacks that are deployed sequentially. `Flows` is collection of flows which can be ordered using `Order` property. `Flows` can run in parallel or sequentially based on the Order property. 
  - Optional `Order` can be any `unsigned` integer. Default `Order` is set to `0`. Flow with lowest orders are deployed first.
  - Optinal `Description`
  - Mandatory `Stacks` which is the collection of CFN stack. Below are the supported attributes of the stack object
    - mandatory `template_file` or `template_url` (only s3 url)
    - mandatory `stack_name`
    - optinal `capabilities`
    - optinal `parameters`
    - optinal `tags`
    - optinal `tags`

**Sample:**
```yaml
Description: Sample CloudFormation Compose file
Vars:
  ENV_NAME: cfn-compose
  ENV_TYPE: nonproduction
Flows:
  SecurityGroup:
    Order: 0
    Description: Creates SecurityGroup
    Stacks:
    - template_file: <cfn-template-path>
      stack_name: stack-name1
      parameters:
        EnvironmentName: '{{ .ENV_NAME }}'
        EnvironmentType: '{{ .ENV_TYPE }}'
      tags:
        EnvironmentName: '{{ .ENV_NAME }}'
        EnvironmentType: '{{ .ENV_TYPE }}'

  EC2Instance:
    Order: 1
    Description: Deploying EC2 Instance
    Stacks:
    - template_file: <cfn-template-path>
      stack_name: stack-name2
      parameters:
        EnvironmentName: '{{ .ENV_NAME }}'
        EnvironmentType: '{{ .ENV_TYPE }}'
      tags:
        EnvironmentName: '{{ .ENV_NAME }}'
        EnvironmentType: '{{ .ENV_TYPE }}'
```

[Details Example](examples/ec2-sqs/Readme.md)

## Limitations
* Supports limited CFN attributes
* No Retry Mechanism
* No Configurable worker pool. One Go routine is spun for every flow.
* Single Compose Configuration can only have up to 50 flows and each flow can contain only upto 50 stacks

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

## Contribution
There is a lot of room for enhancements and you are more than welcome to contribute. If any concerns or recommendations [create issues](https://github.com/rbalman/cfn-compose/issues). If want to contribute [create PR](https://github.com/rbalman/cfn-compose/pulls)

## Contributors
- [Balman Rawat](https://github.com/rbalman)

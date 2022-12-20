# cfn-composea
A command-line tool for managing CloudFormation Stacks at scale.

By Balman Rawat while working at [CloudFactory](https://www.cloudfactory.com/)

## Features
* Create/Update/Delete multiple CloudFormation stacks parallely or sequentially
* Customize the CloudFormation stacks dependency using yaml config
* Delete multiple CloudFormation stacks respecting the creation sequence
* DryRun mode to plan the change
* Generate/Validate/visualize configuration with ease

## Background
As the infrastrucutre evolves and gets complicated we need to manage/maintain multiple CloudFormation Stacks. When we want to `create/update/delete` these stacks we need to manually apply the actions one at a time. Deletion mostly in dev/test environment can be hectic as we should delete the stacks in the reverse of creation order. **cfn-compose** helps to manage multiple stacks that are closely related using declarative language.

![overview image](./docs/images/cfn-compose.svg)

## Background
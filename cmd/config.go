package cmd

import (
	"github.com/spf13/cobra"
	"github.com/balmanrawat/cfn-compose/config"
	"github.com/balmanrawat/cfn-compose/cfn"
	"github.com/balmanrawat/cfn-compose/compose"
	"gopkg.in/yaml.v2"
	"errors"
	"fmt"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "helper functions to work with compose file",
	Aliases: []string{"c"},
	Long:  `can be used to validate, generate, read configuration`,
}

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "validates the compose file configuration",
	Aliases: []string{"vd"},
	Long:  `validates the compose file configuration. It could be helpful when developing and testing out new configuration`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cc, err := config.GetComposeConfig(configFile)
		if err != nil {
			return errors.New(fmt.Sprintf("Failed while fetching compose file: %s\n", err.Error()))
		}

		err = cc.Validate()
		if err != nil {
			return errors.New(fmt.Sprintf("Failed while validating compose file: %s\n", err.Error()))
		}
		
		fmt.Printf("All good!!")
		return nil
	},
}

var visualizeCmd = &cobra.Command{
	Use:   "visualize",
	Short: "visualize flows and stacks",
	Aliases: []string{"vz"},
	Long:  `parses the configuration and shows flows and stacks in defined order`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cc, err := config.GetComposeConfig(configFile)
		if err != nil {
			return errors.New(fmt.Sprintf("Failed while fetching compose file: %s\n", err.Error()))
		}

		err = cc.Validate()
		if err != nil {
			return errors.New(fmt.Sprintf("Failed while validating compose file: %s\n", err.Error()))
		}
		
		flowsMap := compose.SortFlows(cc.Flows)
		compose.PrintFlowsMap(flowsMap)

		return nil
	},
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "generates configuration template",
	Aliases: []string{"gen"},
	Long:  `generates the sample bootstrap template to speed up the process`,
	RunE: func(cmd *cobra.Command, args []string) error {
		sampleConfig := config.ComposeConfig {
			Description: "Sample CloudFormation Compose file",
			Vars: map[string]string{
				"ENV_NAME": "cfn-compose", 
				"ENV_TYPE": "nonproduction", 
				"DelaySeconds": "60",
				"VPC_ID": "",
				"SUBNET_ID": "",
			},
			Flows: map[string]config.Flow{
				"EC2Instance": config.Flow{
					Description: "Creates EC2 Instance Security Group",
					Order: 0,
					Stacks: []cfn.Stack{
						cfn.Stack{
							StackName: "sample-{{ .ENV_NAME }}-security-group",
							TemplateFile: "sg.yml",
							Parameters: map[string]string{"EnvironmentName": "{{ .ENV_NAME }}", "EnvironmentType": "{{ .ENV_TYPE }}","VpcId": "{{ .VPC_ID }}"},
							Tags: map[string]string{"EnvironmentName": "{{ .ENV_NAME }}", "EnvironmentType": "{{ .ENV_TYPE }}"},
						},
						cfn.Stack{
							StackName: "sample-{{ .ENV_NAME }}-ec2-instance",
							TemplateFile: "ec2.yml",
							Parameters: map[string]string{"EnvironmentName": "{{ .ENV_NAME }}", "EnvironmentType": "{{ .ENV_TYPE }}", "SubnetId": "{{ .SUBNET_ID }}"},
							Tags: map[string]string{"EnvironmentName": "{{ .ENV_NAME }}", "EnvironmentType": "{{ .ENV_TYPE }}"},
						},
					},
				},
				"MessageQueue": config.Flow{
					Description: "Deploying Queuing Resources",
					Order: 1,
					Stacks: []cfn.Stack{
						cfn.Stack{
							StackName: "sample-{{ .ENV_NAME }}-sqs",
							TemplateFile: "sqs.yml",
							Parameters: map[string]string{"EnvironmentName": "{{ .ENV_NAME }}", "EnvironmentType": "{{ .ENV_TYPE }}", "DelaySeconds": "{{ .DelaySeconds }}"},
							Tags: map[string]string{"EnvironmentName": "{{ .ENV_NAME }}", "EnvironmentType": "{{ .ENV_TYPE }}"},
						},
					},
				},
			},
		}

		d, err := yaml.Marshal(&sampleConfig)
		if err != nil {
			return errors.New(fmt.Sprintf("Failed to generate sample compose file: %s\n", err.Error()))
		}

		fmt.Printf("%s", d)

		return nil
	},
}

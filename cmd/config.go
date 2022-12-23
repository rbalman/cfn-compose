package cmd

import (
	"errors"
	"fmt"
	"github.com/rbalman/cfn-compose/cfn"
	"github.com/rbalman/cfn-compose/compose"
	"github.com/rbalman/cfn-compose/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var configCmd = &cobra.Command{
	Use:     "config",
	Short:   "Generate, validate and visualize the compose configuration",
	Aliases: []string{"c"},
	Long:    `Generate, validate and visualize the compose configuration`,
}

var validateCmd = &cobra.Command{
	Use:     "validate",
	Short:   "Validates the compose configuration",
	Aliases: []string{"vd"},
	Long:    `Static validation of the compose configuration. helps to debug configuration issues`,
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
	Use:     "visualize",
	Short:   "Visualize the stacks dependencies and creation order",
	Aliases: []string{"vz"},
	Long:    `Visualize the stacks dependencies and creation order specified in the compose file`,
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
		compose.VisualizeFlowsMap(flowsMap)

		return nil
	},
}

var generateCmd = &cobra.Command{
	Use:     "generate",
	Short:   "Generates compose template",
	Aliases: []string{"gen"},
	Long:    `Generates the sample bootstrap compose template`,
	RunE: func(cmd *cobra.Command, args []string) error {
		sampleConfig := config.ComposeConfig{
			Description: "Sample CloudFormation Compose file",
			Vars: map[string]string{
				"ENV_NAME":  "cfn-compose",
				"ENV_TYPE":  "nonproduction",
				"VPC_ID":    "",
				"SUBNET_ID": "",
			},
			Flows: map[string]config.Flow{
				"SecurityGroup": {
					Description: "Creates Sample Security Group",
					Order:       0,
					Stacks: []cfn.Stack{
						{
							StackName:    "sample-{{ .ENV_NAME }}-security-group",
							TemplateFile: "sg.yml",
							Parameters:   map[string]string{"EnvironmentName": "{{ .ENV_NAME }}", "EnvironmentType": "{{ .ENV_TYPE }}", "VpcId": "{{ .VPC_ID }}"},
							Tags:         map[string]string{"EnvironmentName": "{{ .ENV_NAME }}", "EnvironmentType": "{{ .ENV_TYPE }}"},
						},
					},
				},
				"EC2Instance": {
					Description: "Creates EC2 Instance",
					Order:       1,
					Stacks: []cfn.Stack{
						{
							StackName:    "sample-{{ .ENV_NAME }}-ec2-instance",
							TemplateFile: "ec2.yml",
							Parameters:   map[string]string{"EnvironmentName": "{{ .ENV_NAME }}", "EnvironmentType": "{{ .ENV_TYPE }}", "SubnetId": "{{ .SUBNET_ID }}"},
							Tags:         map[string]string{"EnvironmentName": "{{ .ENV_NAME }}", "EnvironmentType": "{{ .ENV_TYPE }}"},
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

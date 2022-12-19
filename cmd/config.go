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

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all the jobs and stacks",
	Aliases: []string{"ls"},
	Long:  `parses the configuration and shows jobs and stacks in defined order`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cc, err := config.GetComposeConfig(configFile)
		if err != nil {
			return errors.New(fmt.Sprintf("Failed while fetching compose file: %s\n", err.Error()))
		}

		err = cc.Validate()
		if err != nil {
			return errors.New(fmt.Sprintf("Failed while validating compose file: %s\n", err.Error()))
		}
		
		jobsMap := compose.SortJobs(cc.Jobs)
		compose.PrintJobsMap(jobsMap)

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
			Vars: map[string]string{"ENV_TYPE": "nonproduction"},
			Jobs: map[string]config.Job{
				"DataStore": config.Job{
					Name: "DataStore",
					Description: "Creates Database and Security Group",
					Order: 0,
					Stacks: []cfn.Stack{
						cfn.Stack{
							StackName: "sample-database-network",
							TemplateFile: "path-to-database-network-cfn-template",
							Parameters: map[string]string{"DatabaseName": "sample-database-network"},
							Tags: map[string]string{"Name": "sample-database-network"},
						},
						cfn.Stack{
							StackName: "sample-database-stack",
							TemplateFile: "path-to-database-cfn-template",
							Parameters: map[string]string{"DatabaseName": "sample-database"},
							Tags: map[string]string{"Name": "sample-database"},
						},
					},
				},
				"LambdaJob": config.Job{
					Name: "LambdaJob",
					Description: "Deploy lambda that uses above datastore",
					Order: 1,
					Stacks: []cfn.Stack{
						cfn.Stack{
							StackName: "lambda-sqs-stack",
							TemplateFile: "path-to-lambda-sqs-template",
							Parameters: map[string]string{"DelaySeconds": "5"},
							Tags: map[string]string{"Name": "lamdbda-publisher-sqs"},
						},
						cfn.Stack{
							StackName: "database-consumer-lambda-stack",
							TemplateFile: "path-to-lambda-cfn-template",
							Parameters: map[string]string{"LambdaName": "database-consumer-lambda"},
							Tags: map[string]string{"Name": "database-consumer-lambda"},
						},
					},
				},
			},
		}

		d, err := yaml.Marshal(&sampleConfig)
		if err != nil {
			return errors.New(fmt.Sprintf("Failed to generate sample compose file: %s\n", err.Error()))
		}

		fmt.Printf("### SAMPLE TEMPLATE ###\n%s", d)

		return nil
	},
}

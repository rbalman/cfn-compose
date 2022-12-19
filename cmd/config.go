package cmd

import (
	"github.com/spf13/cobra"
	"github.com/balmanrawat/cfn-compose/config"
	"github.com/balmanrawat/cfn-compose/compose"
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

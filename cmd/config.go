package cmd

import (
	"github.com/spf13/cobra"
	"github.com/balmanrawat/cfn-compose/config"
	"path/filepath"
	"os"
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
		dir := filepath.Dir(configFile)
		file := filepath.Base(configFile)
		os.Chdir(dir)

		fmt.Println("##########################")
		fmt.Println("# Supplied Configuration #")
		fmt.Println("##########################")
		fmt.Printf("Config: %s\n\n", file)

		config, err := config.Parse(file)
		if err != nil {
			fmt.Printf("Failed while fetching compose file: %s\n", err.Error())
			return err
		}

		err = config.Validate()
		if err != nil {
			fmt.Printf("Failed while validating compose file: %s\n", err.Error())
			return err
		}
		
		fmt.Printf("All good!!")
		return nil
	},
}

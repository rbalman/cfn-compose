package cmd

import (
	"github.com/spf13/cobra"
	"github.com/balmanrawat/cfn-compose/config"
	"github.com/balmanrawat/cfn-compose/compose"
	"github.com/balmanrawat/cfn-compose/libs"
	"path/filepath"
	"os"
	"fmt"
)

var dryRun bool

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "deploy changes in the compose file",
	Aliases: []string{"d"},
	Long:  `Creates/Updates CFN resources if required`,
	Run: func(cmd *cobra.Command, args []string) {
		dir := filepath.Dir(configFile)
		file := filepath.Base(configFile)
		os.Chdir(dir)

		ll := libs.GetLogLevel(logLevel)

		fmt.Println("##########################")
		fmt.Println("# Supplied Configuration #")
		fmt.Println("##########################")
		fmt.Printf("Config: %s\n", file)
		fmt.Printf("DryRun: %t\n\n", dryRun)

		config, err := config.Parse(file)
		if err != nil {
			fmt.Printf("Failed while fetching compose file: %s\n", err.Error())
			os.Exit(1)
		}

		err = config.Validate()
		if err != nil {
			fmt.Printf("Failed while validating compose file: %s\n", err.Error())
			os.Exit(1)
		}

		compose.Apply(config, ll, dryRun)
	},
}

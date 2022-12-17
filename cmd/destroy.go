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

var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "destroys all the stacks which are part of the compose file",
	Aliases: []string{"ds"},
	Long:  `destroy respects the order specified in the compose file and applies the changes accordingly for the individual CFN stacks. It just does the reverse thing deploy sub-command does. Supports dryRun mode, use --dry-run or -d flag.`,
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

		compose.Apply(config, ll, jobName, false, dryRun)
	},
}

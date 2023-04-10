package cmd

import (
	"github.com/rbalman/cfn-compose/compose"
	"github.com/spf13/cobra"
)

var destroyCmd = &cobra.Command{
	Use:     "destroy",
	Short:   "Destroys all the stacks in the reverse order of creation",
	Aliases: []string{"ds"},
	Long:    `Destroys all the stacks in the reverse order of creation as specified in the compose configuration. Supports dryRun mode, use --dry-run or -d flag.`,
	Run: func(cmd *cobra.Command, args []string) {
		c := compose.Composer{
			LogLevel:         logLevel,
			CherryPickedFlow: flowName,
			DeployMode:       false,
			DryRun:           dryRun,
			ConfigFile:       configFile,
			WorkersCount:     workersCount,
		}

		c.PrintConfig()
		c.Apply()
	},
}
